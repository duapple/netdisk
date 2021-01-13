package ddns

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	log "github.com/sirupsen/logrus"
)

var accessKeyId string = ""
var accessSecret string = ""
var recordId string = ""

func DemoDDNS() {
	go SetDDNSService()
	for {
		time.Sleep(time.Duration(100) * time.Second)
	}
}

func SetDDNSService() {
	var WanIP string
	var RecordIP string = GetAliRecordIP() // 服务器启动时，从阿里云获取一次

	for {
		WanIP = GetWanIPStr()
		if WanIP != RecordIP {
			log.Info("Wan IP changed. Will change the record IP.")
			err := SetDDNS(WanIP)
			if err == nil {
				RecordIP = WanIP
			}

		} else {
			//log.Info("Wan IP hold.")
		}
		time.Sleep(time.Duration(3) * time.Second)
	}
}

func SetDDNS(wanIP string) (err error) {
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", accessKeyId, accessSecret)

	request := alidns.CreateUpdateDomainRecordRequest()
	request.Scheme = "https"

	request.RecordId = recordId
	request.RR = "@"
	request.Type = "A"
	request.Value = wanIP //GetWanIPStr() //"118.123.37.212"
	request.Lang = "en"
	request.UserClientIp = wanIP // "118.123.37.211"
	request.TTL = "600"
	request.Priority = "1"
	request.Line = "default"

	response, err := client.UpdateDomainRecord(request)
	if err != nil {
		fmt.Print(err.Error(), response)
		return err
	}
	fmt.Printf("response is %#v\n", response)
	return nil
}

func GetAliRecordIP() (recordIP string) {
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", accessKeyId, accessSecret)

	request := alidns.CreateDescribeDomainRecordInfoRequest()
	request.Scheme = "https"

	request.RecordId = recordId
	request.Lang = "en"
	request.UserClientIp = "118.123.37.211"

	response, err := client.DescribeDomainRecordInfo(request)
	if err != nil {
		fmt.Print(err.Error())
		return ""
	}
	log.Info("Record IP: ", response.Value)
	return response.Value
}

func GetWanIPStr() (wanip string) {
	cmd := exec.Command("curl", "ident.me")
	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal("error: ", err)
	}
	//fmt.Printf("in all caps: %q\n", out.String())

	wanip = out.String()
	if wanip != "" {
		//log.Info("Get WAN IP ok: ", wanip)
	} else {
		log.Warn("Get WAN IP failed")
	}
	return wanip
}

// func GetWanIPStr() (wanip string) {
// 	cmd := exec.Command("curl", "cip.cc")
// 	cmd.Stdin = strings.NewReader("some input")
// 	var out bytes.Buffer
// 	cmd.Stdout = &out
// 	err := cmd.Run()
// 	if err != nil {
// 		log.Fatal("error: ", err)
// 	}
// 	fmt.Printf("in all caps: %q\n", out.String())

// 	wanip = getWanIPStr(out.String())
// 	if wanip != "" {
// 		log.Info("Get WAN IP ok: ", wanip)
// 	} else {
// 		log.Warn("Get WAN IP failed")
// 	}
// 	return wanip
// }

// func getWanIPStr(ipInfo string) string {
// 	var ipInfoByte []byte = []byte(ipInfo)
// 	var len int = len(ipInfoByte)
// 	var preStr string = "IP\t: "
// 	var lenPre int = 5 //len(preStrByte)
// 	log.WithFields(log.Fields{
// 		"LEN": len,
// 	}).Info("")

// 	if len < lenPre {
// 		log.Error("out put ip info string is too short")
// 	}

// 	log.Info("preString: ", string(ipInfoByte[:lenPre]))

// 	if string(ipInfoByte[:lenPre]) != preStr {
// 		log.Error("out put ip info string is error")
// 		log.Error("error string: ", string(ipInfoByte))
// 	}

// 	if index := getEnterPos(string(ipInfoByte)); index > 0 {
// 		return string(ipInfoByte[lenPre:index])
// 	}

// 	return ""
// }

// func getEnterPos(str string) int {
// 	for i := 0; i < len(str); i++ {
// 		if string(str[i]) == "\n" {
// 			return i
// 		}
// 	}

// 	return -1
// }
