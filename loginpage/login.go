package loginpage

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/duapple/netdisk/config"
	"github.com/duapple/netdisk/database"
	"github.com/duapple/netdisk/session"

	log "github.com/sirupsen/logrus"
)

type UserInfo struct {
	Username string
	Password string
}

/* 引用session中的session.GlobalSessions */
var globalSessions *session.Manager

/* init函数执行时，顺序不能确定，需要注意 */
func Init() {
	log.Info("login init")
	globalSessions = session.GlobalSessions
}

func Login(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"HTTP": r.Method,
		"FUNC": "Login",
	}).Info("HTTP REQUEST")
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, err := template.ParseFiles(config.LoginHTML)
		if err != nil {
			log.Error("login.html is not exist")
		}
		t.Execute(w, token)

	} else {
		//请求的是登陆数据，那么执行登陆的逻辑判断
		r.ParseForm()
		log.Info("username length:", len(r.Form["username"][0]))
		log.Info("username:", template.HTMLEscapeString(r.Form.Get("username"))) //输出到服务器端
		log.Info("password:", template.HTMLEscapeString(r.Form.Get("password")))
		// template.HTMLEscape(w, []byte(r.Form.Get("username"))) //输出到客户端

		http.Redirect(w, r, "/upload", 302) //表达功能实现登陆后重定向页面
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"HTTP": r.Method,
		"FUNC": "Logout",
	}).Info("HTTP REQUEST")
	//销毁
	globalSessions.SessionDestroy(w, r)
	log.Info("session destroy")

	crutime := time.Now().Unix()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(crutime, 10))
	token := fmt.Sprintf("%x", h.Sum(nil))

	t, err := template.ParseFiles(config.LoginHTML)
	if err != nil {
		log.Error("login.html is not exist")
		w.WriteHeader(404)
		return
	}
	t.Execute(w, token)
}

var user_t UserInfo = UserInfo{Username: "river", Password: "123123"}

func user_check(user UserInfo) (eq bool) {
	eq = false
	if user.Username == user_t.Username {
		if user.Password == user_t.Password {
			eq = true
		} else {
			log.Error("password error")
		}
	} else {
		log.Error("user not exist")
	}
	return eq
}

func Login_auth(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"HTTP": r.Method,
		"FUNC": "Login_auth",
	}).Info("HTTP REQUEST")

	header := w.Header()
	header.Add("Content-Type", "application/json")

	defer r.Body.Close()
	con, _ := ioutil.ReadAll(r.Body) //获取post的数据
	log.Infof("%T %s\n", con, string(con))

	var user UserInfo
	err := json.Unmarshal([]byte(con), &user)
	if err != nil {
		log.Error("json unmarshal error\n")
		fmt.Fprintf(w, config.FmtStr, "2000", "error", "json unmarshal error")
		return
	}
	log.Info("username: ", user.Username)
	log.Info("password: ", user.Password)

	// eq := user_check(user) //用户校验，验证用户名和密码
	user_t := database.UserInfo{user.Username, user.Password}
	eq := user_t.Check_user_from_db()
	if eq {
		/* session */
		/* globalSessions = session.GlobalSessions //引用session.GlobalSessions */
		sess := globalSessions.SessionStart(w, r)
		val := sess.Get("username")
		if val != nil {
			log.Info(val)
		} else {
			sess.Set("username", user.Username)
			log.Info("set session")
		}
		/* session end */

		/* root dir set */
		current_dir := config.DiskRootDir + user.Username + "/"
		up_dir := current_dir
		sess.Set("current_dir", current_dir)
		sess.Set("up_dir", up_dir)
		fmt.Fprintf(w, config.FmtStr, "1000", "success", "login auth sucess")
	} else {
		log.Info(string(con))
		log.Error("username or password is error")
		fmt.Fprintf(w, config.FmtStr, "2000", "error", "login auth error")
	}
}

func SayhelloName(w http.ResponseWriter, r *http.Request) {
	log.Infof("[%s] func: SayhelloName", r.Method)
	/* session */
	cookie, err := r.Cookie("name")
	if err == nil {
		log.Info(cookie.Value)
		log.Info(cookie.Domain)
		log.Info(cookie.Expires)
	}

	r.ParseForm() //解析url传递的参数，对于POST则解析响应包的主体（request body）
	//注意:如果没有调用ParseForm方法，下面无法获取表单的数据
	log.Info(r.Form) //这些信息是输出到服务器端的打印信息
	log.Info("path", r.URL.Path)
	log.Info("scheme", r.URL.Scheme)
	log.Info(r.Form["url_long"])
	for k, v := range r.Form {
		log.Info("key:", k)
		log.Info("val:", strings.Join(v, ""))
	}

	header := w.Header()
	header.Add("Content-Type", "application/json")
	var json string = "{\"test\": \"test\"}"
	fmt.Fprintf(w, json) //写入到w的是输出到客户端的
}
