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

type uploadRequest struct {
	Option   string
	FileName string
	Size     string
	ChunkNum string
	MD5      string
	ChunkPos string
}

/* 引用session中的session.GlobalSessions */
var globalSessions *session.Manager

func Init() {
	globalSessions = session.GlobalSessions
}

func TransportPage(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"HTTP": r.Method,
		"FUNC": "TransportPage",
	}).Info("HTTP REQUEST")

	/* session authentication */
	sess := session.GlobalSessions.SessionCheck(w, r)
	if sess != nil {
		// log.Info("sess check error")
		// w.WriteHeader(404)
		// return
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, err := template.ParseFiles(config.TransportHTML)
		if err != nil {
			log.WithFields(log.Fields{
				"HTML File": config.TransportHTML,
			}).Error("HTML file is not exist")
		}
		t.Execute(w, token)
		return
	}
}

func SharePage(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"HTTP": r.Method,
		"FUNC": "SharePage",
	}).Info("HTTP REQUEST")

	/* session authentication */
	sess := session.GlobalSessions.SessionCheck(w, r)
	if sess != nil {
		// log.Info("sess check error")
		// w.WriteHeader(404)
		// return
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, err := template.ParseFiles(config.ShareHTML)
		if err != nil {
			log.WithFields(log.Fields{
				"HTML File": config.ShareHTML,
			}).Error("HTML file is not exist")
		}
		t.Execute(w, token)
		return
	}
}

func Download(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"HTTP": r.Method,
		"FUNC": "Download",
	}).Info("HTTP REQUEST")
	/* session authentication */
	sess := globalSessions.SessionCheck(w, r)
	if sess == nil {
		log.Error("session check failed")
		fmt.Fprintf(w, config.FmtStr, "3000", "error", "session check failed")
		return
	}
	current_dir := sess.Get("current_dir")
	log.Info("current dir: ", current_dir)
	up_dir := sess.Get("up_dir")
	log.Info("up_dir: ", up_dir)

	log.Info("method:", r.Method)
	fn := r.FormValue("downloadfile")
	log.Info("test")
	log.Info("fn: ", fn)
	//fn := "ONVIF-AccessControl-Service-Spec (2).pdf"
	//设置响应头
	header := w.Header()
	header.Add("Content-Type", "application/octet-stream")
	header.Add("Content-Disposition", "attachment;filename="+fn)

	/* 实现浏览器进度显示 */
	var fileSize string
	file, err := os.Stat(current_dir.(string) + fn)
	if err != nil {
		log.WithFields(log.Fields{
			"File": current_dir.(string) + fn,
		}).Error("read file stat failed")
	}
	fileSize = fmt.Sprintf("%d", file.Size())
	log.WithFields(log.Fields{
		"fileSize": fileSize,
	}).Info("Info")
	header.Add("Content-Length", fileSize) // 设置响应头的Content-Length属性，用于浏览器显示进度

	//使用ioutil包读取文件
	b, err := ioutil.ReadFile(current_dir.(string) + fn)
	if err != nil {
		log.Errorf("file is not exist: %s", fn)
		fmt.Fprintf(w, config.FmtStr, "3000", "error", "file is not exist")
		return
	}
	//写入到响应流中
	w.Write(b)
}

func UploadRequest(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"HTTP": r.Method,
		"FUNC": "UploadRequest",
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

	log.Info("UploadRequest json: ", string(con))

	var uploadReq uploadRequest
	err := json.Unmarshal([]byte(con), &uploadReq) /* 解析json字符串数据到结构体中 */
	if err != nil {
		log.Error("json unmarshal error")
		fmt.Fprintf(w, config.FmtStr, "2000", "error", "json unmarshal error")
		return
	}

	switch uploadReq.Option {
	case "reUploadFile":
		{
			sess.Set("currentFile", uploadReq)
			err := os.Remove("./tmp/" + uploadReq.FileName + "/" + uploadReq.FileName + "_" + uploadReq.ChunkPos)
			if err != nil {
				log.Error(err)
			}
			fmt.Fprintf(w, config.FmtStr, "1000", "success", "reupload request success")
		}
	case "uploadFile":
		{
			err = os.Mkdir("./tmp/"+uploadReq.FileName, 0777)
			if err != nil {
				log.Info(err)
				fmt.Fprintf(w, config.FmtStr, "3000", "error", "mkdir error")
			} else {
				log.Infof("create dir %s ok\r\n", uploadReq.FileName)
			}
			// 文件上传信息保存，保存到session中，用于分片续传时使用
			sess.Set(uploadReq.FileName, uploadReq)
			sess.Set("currentFile", uploadReq)

			fmt.Fprintf(w, config.FmtStr, "1000", "success", "upload request success")
		}
	case "uploadCancel":
		{
			err = os.RemoveAll("./tmp/" + uploadReq.FileName)
			if err != nil {
				log.Fatal(err)
			}

			sess.Delete(uploadReq.FileName)
			fmt.Fprintf(w, config.FmtStr, "1000", "success", "upload file success")
			return
		}
	default:
		{
			fmt.Fprintf(w, config.FmtStr, "2000", "error", "upload request option error")
		}
	}

}

/* 分片上传功能实现 */
func Upload(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"HTTP": r.Method,
		"FUNC": "Upload",
	}).Info("HTTP REQUEST")

	header := w.Header()
	header.Add("Content-Type", "application/json")

	/* session authentication */
	sess := globalSessions.SessionCheck(w, r) //session 检查的接口封装
	if sess == nil {
		log.Error("session check failed")
		fmt.Fprintf(w, config.FmtStr, "3000", "error", "session check failed")
		return
	}
	current_dir := sess.Get("current_dir")
	log.Info("current dir: ", current_dir)
	up_dir := sess.Get("up_dir")
	log.Info("up_dir: ", up_dir)

	// var current_dir string = "./river/test"
	/* 表单上传文件 */
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, _, err := r.FormFile("uploadfile")
		if err != nil {
			log.Error(err)
			return
		}
		defer file.Close()

		var uploadFileInfo uploadRequest
		uploadFileInfo = sess.Get("currentFile").(uploadRequest)

		// fmt.Fprintf(w, "%v", sess.Get("fileName"))
		var fileName string = uploadFileInfo.FileName + "_" + uploadFileInfo.ChunkPos

		/* 判断上传的文件是否已经存在，可能发生同名的情况 */
		log.Debug("fileName: ", fileName)
		_, err = os.Stat("./tmp/" + uploadFileInfo.FileName + "/" + fileName)
		var f *os.File
		if err == nil {
			fmt.Fprintf(w, config.FmtStr, "3000", "error", "file is exsited")
			return
		} else {
			f, err = os.OpenFile("./tmp/"+uploadFileInfo.FileName+"/"+fileName, os.O_WRONLY|os.O_CREATE, 0666) // 此处假设当前目录下已存在test目录
		}

		defer f.Close()

		if err != nil {
			log.Error(err)
			fmt.Fprintf(w, config.FmtStr, "3000", "error", "create file error")
			return
		}
		io.Copy(f, file) //这里进行大文件copy会导致内存占用过大。一段时候后会进行GC回收
		log.Info("copy ok")

		// 判断是否单个文件是否上传完毕，上传完毕则合并文件到目标文件夹
		if uploadFileInfo.ChunkPos == uploadFileInfo.ChunkNum {
			f.Close()
			_, err = os.Stat(current_dir.(string) + uploadFileInfo.FileName)
			var fii *os.File
			if err == nil {
				// 如果已经存在这个文件了，则在当前文件名后增加日期
				time := time.Now().Format("2006-01-02_15-04-05")
				log.Debug("Time: ", time)

				var fileName []byte = []byte(current_dir.(string) + uploadFileInfo.FileName)
				var fileNameFinal []byte = fileName
				log.Debug("file_name: ", string(fileName))
				var preFileName []byte
				var sufFileName []byte
				for i := len(fileName) - 1; i > 0; i-- {
					if fileName[i] == '.' {
						preFileName = fileName[:i] // 文件名前缀 test
						sufFileName = fileName[i:] // 文件名后缀 .txt
						fileNameFinal = []byte(string(preFileName) + "_" + time + string(sufFileName))
						log.Info("fileNameFile: ", fileNameFinal)
						break
					}
				}

				fii, err = os.OpenFile(string(fileNameFinal), os.O_WRONLY|os.O_CREATE, 0777) // 此处假设当前目录下已存在test目录
			} else {
				fii, err = os.OpenFile(current_dir.(string)+uploadFileInfo.FileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
			}
			defer fii.Close()

			if err != nil {
				log.Error(err)
				fmt.Fprintf(w, config.FmtStr, "3000", "error", "Open object file error")
				return
			}
			index, _ := strconv.Atoi(uploadFileInfo.ChunkNum)
			for i := 1; i <= index; i++ {
				f11, err := os.OpenFile("./tmp/"+uploadFileInfo.FileName+"/"+uploadFileInfo.FileName+"_"+strconv.Itoa(int(i)), os.O_RDONLY, os.ModePerm)
				if err != nil {
					log.Error(err)
					err = os.RemoveAll("./tmp/" + uploadFileInfo.FileName) //出错删除文件
					if err != nil {
						log.Error(err)
					}
					fmt.Fprintf(w, config.FmtStr, "3000", "error", "Open slice file error")
					return
				}
				b, err := ioutil.ReadAll(f11)
				if err != nil {
					log.Error(err)
					err = os.RemoveAll("./tmp/" + uploadFileInfo.FileName) //出错删除文件
					if err != nil {
						log.Error(err)
					}
					fmt.Fprintf(w, config.FmtStr, "3000", "error", "ioutil readall error")
					return
				}
				fii.Write(b)
				f11.Close()
			}

			// err = os.RemoveAll("./tmp/" + uploadFileInfo.FileName)
			// if err != nil {
			// 	log.Error(err)
			// }

			sess.Delete(uploadFileInfo.FileName)

			fii.Close()
			md5str := FileMD5(fii.Name())
			log.Infof("server md5str:%s\n", md5str)
			log.Infof("web md5:%s\n", uploadFileInfo.MD5)
			if uploadFileInfo.MD5 != md5str {
				// err := os.Remove(fii.Name())
				// if err != nil {
				// 	log.Error(err)
				// }
				fmt.Fprintf(w, config.FmtStr, "3000", "error", "upload file md5 error")
			} else {
				err = os.RemoveAll("./tmp/" + uploadFileInfo.FileName)
				if err != nil {
					log.Error(err)
				}
				fmt.Fprintf(w, config.FmtStr, "1000", "success", "upload one file all slice success")
			}
			return
		}

		chunkPos, err := strconv.Atoi(uploadFileInfo.ChunkPos)
		if err != nil {
			log.Error(err)
		}
		uploadFileInfo.ChunkPos = strconv.Itoa(chunkPos + 1)
		sess.Set(uploadFileInfo.FileName, uploadFileInfo)
		sess.Set("currentFile", uploadFileInfo)

		fmt.Fprintf(w, config.FmtStr, "1000", "success", "upload file success")
	}
}

func FileMD5(file string) string {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		log.Info(err)
		return ""
	}

	buffer, _ := ioutil.ReadAll(f)
	data := buffer
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	log.Infof("MD5:%s\n", md5str)
	return md5str
}

// /* 单文件上传 */
// func Upload(w http.ResponseWriter, r *http.Request) {
// 	log.WithFields(log.Fields{
// 		"HTTP": r.Method,
// 		"FUNC": "Upload",
// 	}).Info("HTTP REQUEST")

// 	/* session authentication */
// 	sess := globalSessions.SessionCheck(w, r) //session 检查的接口封装
// 	if sess == nil {
// 		log.Error("session check failed")
// 		w.WriteHeader(404)
// 		return
// 	}
// 	current_dir := sess.Get("current_dir")
// 	log.Info("current dir: ", current_dir)
// 	up_dir := sess.Get("up_dir")
// 	log.Info("up_dir: ", up_dir)

// 	// sid := globalSessions.GetCookieName(w, r) /* 取得当前http会话的cookie */
// 	// //sid, _ := url.QueryUnescape(cookie.Value)       //解析cookie中的session id
// 	// sess := session.GetFromMemory().GetSession(sid) //在服务器本地查找cookie传过来的session id
// 	// if sess == nil {
// 	// 	log.Info("server session id not exist")
// 	// 	w.WriteHeader(404)
// 	// 	return
// 	// }
// 	/* 与上面代码作用相同。识别是否存在session，存在则读出session内容，否则创建新的 */
// 	// sess := globalSessions.SessionStart(w, r)
// 	// val := sess.Get("username")
// 	// if val != nil {
// 	// 	log.Info(val)
// 	// }

// 	/* binary文件上传 */
// 	// log.Info("test")
// 	// f, err := os.OpenFile(disk_dir+"testfile", os.O_WRONLY|os.O_CREATE, 0666) // 此处假设当前目录下已存在test目录
// 	// if err != nil {
// 	// 	log.Fatal("cannot open temp file", err)
// 	// }
// 	// defer f.Close()
// 	// io.Copy(f, r.Body)
// 	// w.WriteHeader("1")

// 	// var current_dir string = "./river/test"
// 	/* 表单上传文件 */
// 	if r.Method == "GET" {
// 		crutime := time.Now().Unix()
// 		h := md5.New()
// 		io.WriteString(h, strconv.FormatInt(crutime, 10))
// 		token := fmt.Sprintf("%x", h.Sum(nil))

// 		t, _ := template.ParseFiles("upload.gtpl")
// 		t.Execute(w, token)
// 	} else {
// 		r.ParseMultipartForm(32 << 20)
// 		file, handler, err := r.FormFile("uploadfile")
// 		if err != nil {
// 			log.Error(err)
// 			return
// 		}
// 		defer file.Close()
// 		fmt.Fprintf(w, "%v", handler.Header)

// 		/* 判断上传的文件是否已经存在，可能发生同名的情况 */
// 		log.Debug("handler.Filename: ", handler.Filename)
// 		_, err = os.Stat(current_dir.(string) + handler.Filename)
// 		var f *os.File
// 		if err == nil {
// 			// 如果已经存在这个文件了，则在当前文件名后增加日期
// 			time := time.Now().Format("2006-01-02_15-04-05")
// 			log.Debug("Time: ", time)

// 			var fileName []byte = []byte(current_dir.(string) + handler.Filename)
// 			var fileNameFinal []byte = fileName
// 			log.Debug("file_name: ", string(fileName))
// 			var preFileName []byte
// 			var sufFileName []byte
// 			for i := len(fileName) - 1; i > 0; i-- {
// 				if fileName[i] == '.' {
// 					preFileName = fileName[:i] // 文件名前缀 test
// 					sufFileName = fileName[i:] // 文件名后缀 .txt
// 					fileNameFinal = []byte(string(preFileName) + "_" + time + string(sufFileName))
// 					log.Info("fileNameFile: ", fileNameFinal)
// 					break
// 				}
// 			}

// 			f, err = os.OpenFile(string(fileNameFinal), os.O_WRONLY|os.O_CREATE, 0666) // 此处假设当前目录下已存在test目录
// 		} else {
// 			f, err = os.OpenFile(current_dir.(string)+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666) // 此处假设当前目录下已存在test目录
// 		}

// 		if err != nil {
// 			log.Error(err)
// 			return
// 		}
// 		defer f.Close()
// 		io.Copy(f, file) //这里进行大文件copy会导致内存占用过大。一段时候后会进行GC回收
// 		log.Info("copy ok")
// 		fmt.Fprintf(w, "1")
// 	}
// }
