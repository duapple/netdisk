package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/duapple/netdisk/config"
	log "github.com/sirupsen/logrus"
)

//session存储方式接口
type Provider interface {
	//初始化一个session，sid根据需要生成后传入
	SessionInit(sid string) (Session, error)
	//根据sid,获取session
	SessionRead(sid string) (Session, error)
	//销毁session
	SessionDestroy(sid string) error
	//回收
	SessionGC(maxLifeTime int64)
}

//Session操作接口
type Session interface {
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(ket interface{}) error
	SessionID() string
}

type Manager struct {
	cookieName  string
	lock        sync.Mutex //互斥锁
	provider    Provider   //存储session方式
	maxLifeTime int64      //有效期
}

var provides = make(map[string]Provider) //生成session存储方式的映射关系

//实例化一个session管理器
func NewSessionManager(provideName, cookieName string, maxLifeTime int64) (*Manager, error) {
	provide, ok := provides[provideName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q ", provideName)
	}
	return &Manager{cookieName: cookieName, provider: provide, maxLifeTime: maxLifeTime}, nil
}

//注册 由实现Provider接口的结构体调用
func Register(name string, provide Provider) {
	if provide == nil {
		panic("session: Register provide is nil")
	}
	if _, ok := provides[name]; ok { /* ok应该为false，说明provider还没有被注册 */
		panic("session: Register called twice for provide " + name)
	}
	/* 没有注册，则注册provider */
	provides[name] = provide
}

//生成sessionId
func (manager *Manager) sessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	//加密
	return base64.URLEncoding.EncodeToString(b)
}

func (manager Manager) GetCookieName(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil {
		return ""
	}
	sid, _ := url.QueryUnescape(cookie.Value) //解析cookie中的session id
	log.Info("sid: ", sid)
	return sid
}

func (manager *Manager) SessionCheck(w http.ResponseWriter, r *http.Request) (session Session) {
	manager.lock.Lock() //加锁
	defer manager.lock.Unlock()
	/* session authentication */
	sid := manager.GetCookieName(w, r) /* 取得当前http会话的cookie */
	//sid, _ := url.QueryUnescape(cookie.Value)       //解析cookie中的session id
	sess := GetFromMemory().GetSession(sid) //在服务器本地查找cookie传过来的session id
	if sess == nil {
		log.Error("server session id not exist")
		return nil
	}
	return sess
}

//判断当前请求的cookie中是否存在有效的session，存在返回，否则创建
func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session Session) {
	manager.lock.Lock() //加锁
	defer manager.lock.Unlock()
	cookie, err := r.Cookie(manager.cookieName) /* 取得当前http会话的cookie */
	if err != nil || cookie.Value == "" {       /* 前面连接的cookie不存在 */
		//创建一个
		sid := manager.sessionId()                     /* 生成sessionID */
		session, _ = manager.provider.SessionInit(sid) /* 通过sessionID 生成session */
		cookie := http.Cookie{
			Name:     manager.cookieName,
			Value:    url.QueryEscape(sid), //转义特殊符号@#s￥%+*-等
			Path:     "/",
			Domain:   config.HostIPV4,
			HttpOnly: true,
			MaxAge:   int(manager.maxLifeTime),
			//Expires:  time.Now().Add(time.Duration(manager.maxLifeTime)), //不设置此属性，将在浏览器关闭时关闭cookie
			//MaxAge和Expires都可以设置cookie持久化时的过期时长，Expires是老式的过期方法，
			// 如果可以，应该使用MaxAge设置过期时间，但有些老版本的浏览器不支持MaxAge。
			// 如果要支持所有浏览器，要么使用Expires，要么同时使用MaxAge和Expires。
		}
		http.SetCookie(w, &cookie)
	} else {
		log.Info("cookie value: ", cookie.Value)
		sid, _ := url.QueryUnescape(cookie.Value) //反转义特殊符号
		session, _ = manager.provider.SessionRead(sid)
	}
	return session
}

//销毁session 同时删除cookie
func (manager *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		return
	} else {
		manager.lock.Lock()
		defer manager.lock.Unlock()
		sid, _ := url.QueryUnescape(cookie.Value)
		manager.provider.SessionDestroy(sid) /* 删除session */
		expiration := time.Now()
		cookie := http.Cookie{
			Name:     manager.cookieName,
			Path:     "/",
			HttpOnly: true,
			Expires:  expiration,
			MaxAge:   -1}
		http.SetCookie(w, &cookie)
	}
}

func (manager *Manager) GC() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.provider.SessionGC(manager.maxLifeTime)
	time.AfterFunc(time.Duration(manager.maxLifeTime), func() { manager.GC() })
}
