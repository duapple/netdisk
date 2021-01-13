package homepage

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/duapple/netdisk/config"
	"github.com/duapple/netdisk/session"
	log "github.com/sirupsen/logrus"
)

type Dir struct {
	DirName string
	Size    string
	ModTime string
	//ModTime time.Time
}

type File struct {
	FileName string
	Size     string
	ModTime  string
	//ModTime time.Time
}

type DirInfo struct {
	CurrentDir string
	Dirs       []Dir
	Files      []File
}

type Dir_Opt_e int32

/* 目录操作指令枚举类型 */
const (
	DirOptRead   Dir_Opt_e = 0 /* 读取目录内容 */
	DirOptCreate Dir_Opt_e = 1 /* 创建目录 */
	DirOptRemove Dir_Opt_e = 2 /* 删除目录 */
	DirOptRename Dir_Opt_e = 3 /* 重命名目录 */
)

type _Dir_Opt struct {
	Opt     Dir_Opt_e /* 目录操作指令 */
	DirName []string  /* 待目录名或者文件名 */
	/* 当操作指令为READ或者CREATE时，DirName数组只有第一个数据有效。
	 * 当为REMOVE时，整个数组数据都有效 */
}

func Index(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"HTTP": r.Method,
		"FUNC": "Index",
	}).Info("HTTP REQUEST")

	/* session authentication */
	sess := session.GlobalSessions.SessionCheck(w, r)
	if sess == nil {
		// log.Info("sess check error")
		// w.WriteHeader(404)
		// return
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, err := template.ParseFiles(config.LoginHTML)
		if err != nil {
			log.Error("login.html is not exist")
		}
		t.Execute(w, token)
		return
	}

	crutime := time.Now().Unix()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(crutime, 10))
	token := fmt.Sprintf("%x", h.Sum(nil))

	t, err := template.ParseFiles(config.IndexHTML)
	if err != nil {
		log.Error("index.html is not exist")
	}
	t.Execute(w, token)
}

func Dir_Opt(w http.ResponseWriter, r *http.Request) {

	log.WithFields(log.Fields{
		"HTTP": r.Method,
		"FUNC": "Dir_Opt",
	}).Info("HTTP REQUEST")

	header := w.Header()
	header.Add("Content-Type", "application/json")

	/* session authentication */
	sess := session.GlobalSessions.SessionCheck(w, r)
	if sess == nil {
		log.Error("sess check error")
		fmt.Fprintf(w, config.FmtStr, "3000", "error", "session check error")
		return
	}

	defer r.Body.Close()
	con, _ := ioutil.ReadAll(r.Body) //获取post的body数据

	var dir_Opt _Dir_Opt
	err := json.Unmarshal([]byte(con), &dir_Opt) /* 解析json字符串数据到结构体中 */
	if err != nil {
		log.Error("json unmarshal error")
		fmt.Fprintf(w, config.FmtStr, "2000", "error", "json unmarshal error")
		return
	}
	var currentDir string
	var upDir string
	/* 取得session中的当前目录 */
	current_Dir_sess_get := sess.Get("current_dir")
	up_Dir_sess_get := sess.Get("up_dir")
	username_sess_get := sess.Get("username")

	if up_Dir_sess_get != nil {
		upDir = up_Dir_sess_get.(string)
	} else {
		log.Error("up_dir not in sess")
	}

	var dir_pre string
	/* session中存在current_dir时，取得 */
	if current_Dir_sess_get != nil {
		dir_pre = current_Dir_sess_get.(string)
	} else {
		log.Warn("current_dir not in session")
		dir_pre = config.DiskRootDir + username_sess_get.(string)
	}

	/* 操作不同路径访问 */
	if dir_Opt.DirName[0] == ".." { /* 访问上一级目录，当current_dir在用户根目录时，UpDir为用户根目录 */
		currentDir = upDir
	} else if dir_Opt.DirName[0] == "." { /* 访问当前目录，current_dir不变 */
		currentDir = dir_pre
	} else if dir_Opt.DirName[0] == "" { /* 访问用户根目录，用以从任何地方跳转到用户根目录 */
		currentDir = config.DiskRootDir + sess.Get("username").(string) + "/"
	} else {
		currentDir = dir_pre + dir_Opt.DirName[0] + "/" /* 访问请求指定的目录，目录不存在时，将返回错误信息 */
	}

	switch dir_Opt.Opt {
	case DirOptRead: /* 读目录操作 */
		log.WithFields(log.Fields{
			"currentDir": currentDir,
		}).Info("")
		upDir = update_updir(currentDir, config.DiskRootDir+username_sess_get.(string)+"/") /* 根据当前路径得到上一级路径 */
		log.WithFields(log.Fields{
			"upDir": upDir,
		}).Info("")

		var dirInfo DirInfo
		dirJSONString, err := dirInfo.read_dir(currentDir)
		if err != nil {
			log.Error(err)
			fmt.Fprintf(w, config.FmtStr, "3000", "error", "read dir error")
			return
		}
		/* 将当前路径和上一级路径写入到session */
		sess.Set("current_dir", currentDir)
		sess.Set("up_dir", upDir)
		// var dirJSONString1 string
		// dirJSONString1 = "{\"CurrentDir\":\"" + currentDir + "\"," + dirJSONString + "}"
		log.Info(dirJSONString)
		fmt.Fprintf(w, dirJSONString)

	case DirOptCreate:
		err := os.Mkdir(currentDir, 0777)
		if err != nil {
			log.Info(err)
			fmt.Fprintf(w, config.FmtStr, "3000", "error", "make dir error")
		} else {
			log.Infof("create dir %s ok\r\n", dir_Opt.DirName[0])
			fi, _ := os.Stat(currentDir) //获取目录或者文件信息
			tmp := []byte(fmt.Sprintf("%s", fi.ModTime()))
			date := tmp[:10]
			time := tmp[11:19]
			modTime := fmt.Sprintf("%s %s", string(date), string(time))
			fmt.Fprintf(w, config.FmtStr, "1000", "success", modTime)
		}

	case DirOptRemove:
		upDir = update_updir(currentDir, config.DiskRootDir+username_sess_get.(string)+"/")

		for i := range dir_Opt.DirName {
			log.Info("i = ", i)
			log.Infof("dir_Opt.DirName[%d]: %s", i, dir_Opt.DirName[i])
			if dir_Opt.DirName[i] != "" && dir_Opt.DirName[i] != ".." && dir_Opt.DirName[i] != "." {
				err := os.Remove(upDir + dir_Opt.DirName[i])
				//log.Error("err: ", err)
				if err != nil {
					log.Error(err)
					err = os.RemoveAll(upDir + dir_Opt.DirName[i])
					if err != nil {
						log.Error(err)
						fmt.Fprintf(w, config.FmtStr, "3000", "error", "remove dir error")
					}
				} else {
					log.Infof("remove dir %s ok\r\n", dir_Opt.DirName[i])
				}
			} else {
				log.Error("DirName error")
			}
		}
		fmt.Fprintf(w, config.FmtStr, "1000", "success", "remove dir success")

	case DirOptRename:
		upDir = update_updir(currentDir, config.DiskRootDir+username_sess_get.(string)+"/")
		var tmp []byte = []byte(currentDir)
		len := len(tmp)
		currentDir_t := string(tmp[:len-1])
		newDirName := upDir + dir_Opt.DirName[1]
		log.Infof("currentDir: %s, newDirName: %s\r\n", currentDir_t, newDirName)
		err := os.Rename(currentDir_t, newDirName)
		if err != nil {
			log.Errorf("rename %s error\r\n", currentDir)
			log.Info(err)
			fmt.Fprintf(w, config.FmtStr, "3000", "error", "rename dir error")
			return
		}
		log.Infof("rename %s to %s\r\n", currentDir_t, newDirName)
		fmt.Fprintf(w, config.FmtStr, "1000", "sucess", "rename dir success")
	}
}

func update_updir(currentDir, RootDir string) (UpDir string) {
	var currentDir_t []byte = []byte(currentDir)
	var flag int16 = 0
	/* 得到当前路径的上一级路径，当前路径已经是根目录的话，则保持上一级路径与当前路径相同 */
	for i := len(currentDir) - 1; i >= len(RootDir); i-- {
		if currentDir_t[i] == '/' {
			if flag == 1 {
				break
			}
			flag = 1
		}
		currentDir_t = currentDir_t[:i]
	}
	return string(currentDir_t)
}

func (dirInfo DirInfo) read_dir(dirName string) (jsonString string, err error) {
	//获取文件或目录相关信息
	fileInfoList, err := ioutil.ReadDir(dirName)
	if err != nil {
		log.Info(err)
		return "", err
	}
	log.WithFields(log.Fields{
		"FileAndDir": len(fileInfoList),
	}).Info("")

	dirInfo.CurrentDir = dirName

	for i := range fileInfoList {
		fi, _ := os.Stat(dirName + fileInfoList[i].Name()) //获取目录或者文件信息
		tmp := []byte(fmt.Sprintf("%s", fi.ModTime()))
		date := tmp[:10]
		time := tmp[11:19]
		modTime := fmt.Sprintf("%s %s", string(date), string(time))
		size_int64 := fi.Size() / 1024
		var size string
		if size_int64 >= 1024 {
			size_int64 /= 1024
			if size_int64 >= 1024 {
				size_int64 /= 1024
				size = fmt.Sprintf("%dG", size_int64)
			} else {
				size = fmt.Sprintf("%dM", size_int64)
			}
		} else {
			size = fmt.Sprintf("%dKB", size_int64)
		}
		if fi != nil {
			if fi.IsDir() {
				// log.Info(fileInfoList[i].Name())                                                                       //打印当前文件或目录下的文件或目录名
				dirInfo.Dirs = append(dirInfo.Dirs, Dir{DirName: fileInfoList[i].Name(), Size: "-", ModTime: modTime}) //ModDate: string(ModDate_t), ModTime: string(ModTime_t)}) //添加目录信息到dir结构体
			} else {
				dirInfo.Files = append(dirInfo.Files, File{FileName: fileInfoList[i].Name(), Size: size, ModTime: modTime}) //ModDate: string(ModDate_t), ModTime: string(ModTime_t)}) //添加文件信息到dir结构体
			}
		}
	}

	b, err := json.Marshal(dirInfo) //将dir结构体解析为json字符串
	if err != nil {
		log.Error(err)
		return "", err
	}
	jsonString = string(b)
	return jsonString, err
}
