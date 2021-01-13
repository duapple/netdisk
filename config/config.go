package config

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"
)

var config_t config
var config_json string = "./config/config.json" // 下面的配置参数来自于这个json文件
var FmtStr string = "{\"code\":\"%s\",\"cause\":\"%s\",\"description\":\"%s\"}"

/* 查询主机IP地址，在这个选项中进行设置 */

var (
	HostIPV4 string
	HostPort string

	dir        string
	LoginHTML  string
	IndexHTML  string
	RegistHTML string

	database_dir  string
	UserDBJSON    string
	UserMapDBJSON string

	DiskRootDir string

	TransportHTML string
	ShareHTML     string
)

type config struct {
	HostIPV4           string
	HostPort           string
	StaticDir          string
	LoginHTML          string
	IndexHTML          string
	RegistHTML         string
	DatabaseDir        string
	UserDBJSON         string
	UserMapDBJSON      string
	DiskRootDir        string
	LogSetReportCaller bool
	LogSetFormatter    int8
	LogSetOutput       int8
	LogSetLevel        int8
	TransportHTML      string
	ShareHTML          string
}

func Init() {
	f, err := os.Open(config_json)
	defer f.Close()
	if err != nil {
		log.Errorf("open %s file failed\n", config_json)
		return
	}
	file_info, err := f.Stat()
	if err != nil {
		log.Error("get file information error")
		return
	}
	filesize := file_info.Size()
	buff := make([]byte, filesize)
	_, err = f.Read(buff)
	if err != nil {
		log.Errorf("read %s file error\n", config_json)
		return
	}

	log.Infof("buff:%s#\n", string(buff))

	err = json.Unmarshal([]byte(buff), &config_t)
	if err != nil {
		log.WithFields(log.Fields{
			"FILE": config_json,
		}).Fatal("json unmarshal error")
		return
	}

	// 全局变量初始化
	HostIPV4 = config_t.HostIPV4
	HostPort = config_t.HostPort

	dir = config_t.StaticDir
	LoginHTML = dir + config_t.LoginHTML
	IndexHTML = dir + config_t.IndexHTML
	RegistHTML = dir + config_t.RegistHTML
	TransportHTML = dir + config_t.TransportHTML
	ShareHTML = dir + config_t.ShareHTML

	database_dir = config_t.DatabaseDir
	UserDBJSON = database_dir + config_t.UserDBJSON
	UserMapDBJSON = database_dir + config_t.UserMapDBJSON

	DiskRootDir = config_t.DiskRootDir

	log.WithFields(log.Fields{
		"HostIPV4":      HostIPV4,
		"HostPort":      HostPort,
		"dir":           dir,
		"LoginHTML":     LoginHTML,
		"IndexHTML":     IndexHTML,
		"RegistHTML":    RegistHTML,
		"database_dir":  database_dir,
		"UserDBJSON":    UserDBJSON,
		"UserMapDBJSON": UserMapDBJSON,
		"DiskRootDir":   DiskRootDir,
	}).Info("config info: ")

	log.Info("config module init ok")

	log_init() // 初始化logrus模块
}

// /* 查询主机IP地址，在这个选项中进行设置 */
// var HostIPV4 string = "localhost" // 主机IP地址
// // var HostIPV4 string = "192.168.43.20" // 主机IP地址

// var HostPort string = ":9090" // 主机服务端口号

// var dir string = "./static/"                //js css 文件路径
// var LoginHTML string = dir + "login.html"   //登陆页面
// var IndexHTML string = dir + "index.html"   //用户主页
// var RegistHTML string = dir + "regist.html" //注册页面

// var database_dir string = "./database/"
// var UserDBJSON string = database_dir + "user_db.json"        //用户数据文件
// var UserMapDBJSON string = database_dir + "user_db_map.json" //用户数据文件映射

// var DiskRootDir string = "./DiskRoot/" //磁盘根目录
