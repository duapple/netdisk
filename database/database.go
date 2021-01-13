package database

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/duapple/netdisk/config"
	log "github.com/sirupsen/logrus"
)

type UserInfo struct {
	UserName string
	PassWord string
}

type users_db map[string]int

type usersDatabase struct {
	Users []UserInfo
}

/*  */
var Users_map users_db  //数组存放用户账号数据
var Users usersDatabase //将用户名和数组进行映射

func Init() {
	/* 从json文件中取得所有用户数据 */
	Get_user_db(config.UserDBJSON)
	Get_user_map_db(config.UserMapDBJSON)
}

/* 从文件中读取用户账户数据 */
func (users *usersDatabase) get_users_data_from_json(db_json_name string) (err error) {
	log.Info("[func] get_users_data_from_json")

	f, err := os.Open(db_json_name)
	// f, err := os.OpenFile(db_json_name, os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()
	if err != nil {
		log.Error("open users data json file failed")
		return err
	}

	file_info, err := f.Stat()
	if err != nil {
		log.Error("get file info err")
		return err
	}
	filesize := file_info.Size()
	buff := make([]byte, filesize)
	_, err = f.Read(buff)
	if err != nil {
		log.Errorf("read %s file error\n", db_json_name)
		return err
	}

	log.Infof("buff:%s#\n", string(buff))
	//users.username = make(map[string]int, 0)

	err = json.Unmarshal([]byte(buff), &users) /* 解析json字符串数据到结构体中 */

	if err != nil {
		log.Fatal("json unmarshal error")
		return err
	}
	//log.Info(users)
	return nil
}

/* 从文件中获取用户映射数据 */
func (users *users_db) get_users_map_from_json(db_map_json_name string) (err error) {
	log.Info("[func] get_users_map_from_json")

	f, err := os.Open(db_map_json_name)
	// f, err := os.OpenFile(db_json_name, os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()
	if err != nil {
		log.Error("open users data json file failed")
		return err
	}

	file_info, err := f.Stat()
	if err != nil {
		log.Error("get file info err")
		return err
	}
	filesize := file_info.Size()
	buff := make([]byte, filesize)
	_, err = f.Read(buff)
	if err != nil {
		log.Errorf("read %s file error\n", db_map_json_name)
		return err
	}

	log.Infof("buff:%s#\n", string(buff))

	err = json.Unmarshal([]byte(buff), &users) // 解析到映射中
	return nil
}

/* 将新增用户数据写入到json中 */
func (user UserInfo) add_user_data_to_json(db_json_name string) (err error) {
	log.Info("[func] add_user_data_to_json")
	var users_data usersDatabase
	/* 先从文件中读取用户数据 */
	err = users_data.get_users_data_from_json(db_json_name)
	if err != nil {
		log.Warn("create acount failed")
		return err
	}

	/* 从文件中读取用户映射信息 */
	Users_map = make(users_db)
	Users_map.get_users_map_from_json(config.UserMapDBJSON)
	_, ok := Users_map[user.UserName]
	if ok {
		log.Warn("user has been existed")
		return nil
	}

	/* 新增用户添加到映射中 */
	len := len(users_data.Users)
	log.Debug("len: ", len)
	Users_map[user.UserName] = len
	// log.Debug("Users_map: ", Users_map)
	m, err := json.Marshal(Users_map)
	log.Debug("Users_map: ", string(m))
	f1, err := os.OpenFile(config.UserMapDBJSON, os.O_WRONLY|os.O_CREATE, 0666)
	defer f1.Close()
	if err != nil {
		log.Error("open users data json file failed")
		return err
	}
	f1.Write(m) //写入新增用户后的映射到json文件中

	/* 新增用户添加到用户账户信息结构体中 */
	users_data.Users = append(users_data.Users, UserInfo{UserName: user.UserName, PassWord: user.PassWord})
	log.Info(users_data)
	f2, err := os.OpenFile(db_json_name, os.O_WRONLY|os.O_CREATE, 0666)
	defer f2.Close()
	if err != nil {
		log.Error("open users data json file failed")
		return err
	}
	b, err := json.Marshal(users_data)
	log.Info(string(b))
	f2.Write(b) //将新增用户后的结构体写入到json文件中
	return nil
}

func Get_user_db(users_db_json string) (err error) {
	log.Info("[func] Get_user_db")

	err = Users.get_users_data_from_json(users_db_json)
	/* 添加用户 */
	// var user1 UserInfo = UserInfo{"hejiang", "321321"}
	// user1.add_user_data_to_json(users_db_json)
	if err != nil {
		return err
	}
	return nil
}

func Get_user_map_db(users_map_db_json string) (err error) {
	log.Info("[func] Get_user_map_db")

	err = Users_map.get_users_map_from_json(users_map_db_json)
	if err != nil {
		return err
	}
	return nil
}

func UpdateDB(users_db_json, users_map_db_json string) (err error) {
	err1 := Get_user_db(users_db_json)
	err2 := Get_user_map_db(users_map_db_json)
	if err1 != nil || err2 != nil {
		return err1
	}
	return nil
}

func Regist(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {

		log.WithFields(log.Fields{
			"HTTP": r.Method,
			"FUNC": "Regist",
		}).Info("HTTP REQUEST")

		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, err := template.ParseFiles(config.RegistHTML)
		if err != nil {
			log.Error("regist.html is not exist")
		}
		t.Execute(w, token)
	} else {
		AddUser(w, r)
	}
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"HTTP": r.Method,
		"FUNC": "AddUser",
	}).Info("HTTP REQUEST")

	header := w.Header()
	header.Add("Content-Type", "application/json")

	defer r.Body.Close()
	con, _ := ioutil.ReadAll(r.Body) //获取post的数据
	//fmt.Printf("%T %v\n", con, con)

	var user UserInfo
	err := json.Unmarshal([]byte(con), &user)
	if err != nil {
		log.Error("json unmarshal error\n")
		fmt.Fprintf(w, config.FmtStr, "2000", "error", "json unmarshal error")
		return
	}
	log.Info("add username: ", user.UserName)
	log.Info("add password: ", user.PassWord)

	var user_w UserInfo

	user_w.PassWord = fmt.Sprintf("%x", md5.Sum([]byte(user.PassWord))) //保存密码的MD5
	user_w.UserName = user.UserName
	log.WithFields(log.Fields{
		"PassWord MD5": user_w.PassWord,
	}).Info("")

	// err = user_w.add_user_data_to_json(config.UserDBJSON) // 使用密码的MD5作为密码
	err = user.add_user_data_to_json(config.UserDBJSON) // 使用原密码作为密码

	if err != nil {
		log.Error("add user data to json error")
		fmt.Fprintf(w, config.FmtStr, "3000", "error", "add user data to json error")
		return
	}
	/* 更新Users数据库 */
	Get_user_db(config.UserDBJSON)

	fmt.Fprintf(w, config.FmtStr, "1000", "success", "add user data success")

}

func (UserInfo_t UserInfo) Check_user_from_db() (ret bool) {
	log.Infof("UserInfo_t.UserName: %s\n", UserInfo_t.UserName)
	index, ok := Users_map[UserInfo_t.UserName]
	log.Infof("index: %d\n", index)
	if !ok {
		log.Warn("user error")
		return false
	}
	userInfo_t := Users.Users[index]
	if userInfo_t.PassWord != UserInfo_t.PassWord {
		log.Warn("password error")
		return false
	}

	fi, _ := os.Stat(config.DiskRootDir + UserInfo_t.UserName) //获取目录或者文件信息
	// if err != nil {
	// 	log.Info(err)
	// 	return false
	// }
	if fi == nil {
		err := os.Mkdir(config.DiskRootDir+UserInfo_t.UserName, 0777)
		if err != nil {
			log.Error(err)
			return false
		} else {
			log.Infof("create dir %s ok\n", config.DiskRootDir+UserInfo_t.UserName)
		}
	}

	return true
}
