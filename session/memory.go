package session

import (
	"container/list"
	"sync"
	"time"
)

//session来自内存 实现
type FromMemory struct {
	lock     sync.Mutex               //用来锁
	sessions map[string]*list.Element //用来存储在内存, 将list中的元素和 sid进行映射，方便查找链表元素
	list     *list.List               //用来做 gc
}

/* golang 全局变量初始化是在何时完成的 */
var pder = &FromMemory{list: list.New()} //生成用于存储SessionStore的链表

func GetFromMemory() *FromMemory {
	return pder
}

/* 为什么要用方法而不是函数？ */
/* 接口是方法集合，方法可以单独存在 */
func (pder FromMemory) GetSession(sid string) (sess Session) {
	elem, ok := pder.sessions[sid]
	if ok {
		sess = elem.Value.(*SessionStore) /* list中的element的Value是interface{}, 存储的是SessionStore，通过类型断言得到SessionStore，然后赋值给SessionStore实现的接口 */
	} else {
		sess = nil
	}
	return sess
}

// func init() {
// 	fmt.Println("session memory init")
// 	pder.sessions = make(map[string]*list.Element, 0)
// 	//注册  memory 调用的时候一定要一致
// 	Register("memory", pder)
// }

//session实现
type SessionStore struct {
	sid              string                      //session id 唯一标示
	LastAccessedTime time.Time                   //最后访问时间
	value            map[interface{}]interface{} //session 里面存储的值
}

//设置
func (st *SessionStore) Set(key, value interface{}) error {
	st.value[key] = value
	pder.SessionUpdate(st.sid)
	return nil
}

//获取session
func (st *SessionStore) Get(key interface{}) interface{} {
	pder.SessionUpdate(st.sid)
	if v, ok := st.value[key]; ok {
		return v
	} else {
		return nil
	}
	return nil
}

//删除
func (st *SessionStore) Delete(key interface{}) error {
	/* 删除value和key的映射 */
	delete(st.value, key)
	pder.SessionUpdate(st.sid)
	return nil
}

func (st *SessionStore) SessionID() string {
	return st.sid
}

/* 初始化Session */
func (frommemory *FromMemory) SessionInit(sid string) (Session, error) {
	frommemory.lock.Lock()
	defer frommemory.lock.Unlock()
	v := make(map[interface{}]interface{}, 0)
	/* 将SessionStore结构体内容地址放入到 newsess 指针中。 */
	newsess := &SessionStore{sid: sid, LastAccessedTime: time.Now(), value: v}
	/* 在list最后插入一个新元素（指针），并返回这个元素（指针）。 */
	element := frommemory.list.PushBack(newsess)

	/* 将Fromemory中的sessions成员设置为 SessionStore 指针。 */
	frommemory.sessions[sid] = element //通过sessionId 与 SessionStore结构体做映射
	/* 当某一个类型实现了某个接口时，这个接口变量可以等于这个类型值。
	*  如 *SessionStore 实现了session.Session接口，这里要求的返回值类型为：Session = *SessionStore 。
	 */
	return newsess, nil
}

/* 通过sessionID从list中读取Session，不存在则新创建一个 */
func (frommemory *FromMemory) SessionRead(sid string) (Session, error) {
	if element, ok := frommemory.sessions[sid]; ok { // 利用映射数据结构，通过sid查找相应的session。
		/* element是FromMemory结构体的成员， 来自映射类型：map[string]*list.Element，链表元素。
		*  list包中的Element结构体
		*  Value是Element结构体成员，类型是interface{}，用于存放任意数据类型
		 */
		/* i.(T)类型断言，判断Value接口存放的值是否是*SessionStore类型，是则返回该类型的值，否则Panic */
		return element.Value.(*SessionStore), nil
	} else {
		sess, err := frommemory.SessionInit(sid)
		return sess, err
	}
	return nil, nil
}

func (frommemory *FromMemory) SessionDestroy(sid string) error {
	if element, ok := frommemory.sessions[sid]; ok {
		delete(frommemory.sessions, sid) /* 删除session */
		frommemory.list.Remove(element)  /* 将被删除的session对应的element从list中移除 */
		return nil
	}
	return nil
}

func (frommemory *FromMemory) SessionGC(maxLifeTime int64) {
	frommemory.lock.Lock()
	defer frommemory.lock.Unlock()
	for {
		element := frommemory.list.Back()
		if element == nil { /* list为空，结束 */
			break
		}
		/* 回收内存空间 */
		if (element.Value.(*SessionStore).LastAccessedTime.Unix() + maxLifeTime) <
			time.Now().Unix() { /* 超过maxLiftTime，则从list中移除session对应element */
			frommemory.list.Remove(element)
			delete(frommemory.sessions, element.Value.(*SessionStore).sid) /* 删除该session的映射 */
		} else {
			break
		}
	}
}

func (frommemory *FromMemory) SessionUpdate(sid string) error {
	frommemory.lock.Lock()
	defer frommemory.lock.Unlock()
	if element, ok := frommemory.sessions[sid]; ok {
		element.Value.(*SessionStore).LastAccessedTime = time.Now()
		frommemory.list.MoveToFront(element) /* 更新session时间，然后将session移动到list最前面。应该是在修改session时，执行这个操作 */
		return nil
	}
	return nil
}
