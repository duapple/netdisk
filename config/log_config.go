package config

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// var LogSetReportCaller bool
// var LogSetFormatter int8
// var LogSetOutput int8
// var LogSetLevel int8

func log_init() {
	if config_t.LogSetReportCaller {
		log.SetReportCaller(true)
	}

	// 设置Text文本格式输出log
	if config_t.LogSetFormatter == 0 {
		log.SetFormatter(&log.TextFormatter{
			// DisableColors: false,
			FullTimestamp: true,
			ForceColors:   true,
		})
	} else if config_t.LogSetFormatter == 1 {
		log.SetFormatter(&log.JSONFormatter{})
	}

	switch config_t.LogSetOutput {
	case 0:
		log.SetOutput(os.Stdout)
	case 1:
		log.Info("create logrus.log")
		// You could set this to any `io.Writer` such as a file
		file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log.SetOutput(file)
		} else {
			log.Info("Failed to log to file, using default stderr")
		}
	}

	switch config_t.LogSetLevel {
	case 0:
		log.SetLevel(log.DebugLevel)
	case 1:
		log.SetLevel(log.InfoLevel)
	}

	log.Info("logrus module init ok")
}
