package session

import (
	"container/list"

	log "github.com/sirupsen/logrus"
)

var GlobalSessions *Manager

func Init() {
	log.Info("session main init")

	pder.sessions = make(map[string]*list.Element, 0)
	//注册  memory 调用的时候一定要一致
	/* pder接口体地址传给其实现的接口 */
	Register("memory", pder)

	var err error
	GlobalSessions, err = NewSessionManager("memory", "goSessionid", 360000)
	if err != nil {
		log.Info(err)
		return
	}
	go GlobalSessions.GC()
	log.Info("fd")
}
