package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/duapple/netdisk/config"
	"github.com/duapple/netdisk/database"
	"github.com/duapple/netdisk/homepage"
	"github.com/duapple/netdisk/loginpage"
	"github.com/duapple/netdisk/session"
)

func init() {
	// 使用init时，需要注意模块之间的相互依赖问题
	// loginpage 和 homepage 依赖于session，所以需要先初始化session，否则会导致出错
	config.Init()
	session.Init()
	loginpage.Init()
	homepage.Init()
	database.Init()
}

func main() {
	// go ddns.SetDDNSService()

	http.HandleFunc("/hello", loginpage.SayhelloName)                                          //设置访问的路由
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static")))) // 启动静态文件服务

	// http.HandleFunc("/", loginpage.Login)
	http.HandleFunc("/login", loginpage.Login)           //设置访问的路由
	http.HandleFunc("/login_auth", loginpage.Login_auth) //设置获取登陆信息路由
	http.HandleFunc("/logout", loginpage.Logout)         //销毁                                                    //设置访问的路由
	// http.HandleFunc("/count", count)
	http.HandleFunc("/regist", database.Regist) //注册用户

	http.HandleFunc("/upload", homepage.Upload)
	http.HandleFunc("/home", homepage.Dir_Opt)
	http.HandleFunc("/download", homepage.Download)
	http.HandleFunc("/index", homepage.Index)
	http.HandleFunc("/transport", homepage.TransportPage)
	http.HandleFunc("/share", homepage.SharePage)
	http.HandleFunc("/upload_request", homepage.UploadRequest)

	log.Info("route setting ok")

	err := http.ListenAndServe(config.HostPort, nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
