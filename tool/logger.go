/**
 * @Author: LOFTER
 * @Description:
 * @File:  logger
 * @Date: 2021/2/9 6:35 下午
 */
package tool

import (
	"errors"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"time"
	"upload/util"
)

var Log = logrus.New()

func InitLogger() {
	if util.IsFileExist("applog") == false {
		_ = os.MkdirAll("applog", 0755)
	}
	file, err := os.OpenFile("applog/app-"+time.Now().Format("2006-01-02")+".log", os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)
	if err != nil {
		log.Fatalln("log init failed")
	}

	var info os.FileInfo
	info, err = file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fileWriter := logFileWriter{file, info.Size()}
	Log.SetOutput(&fileWriter)
	//log.SetLevel(log.ErrorLevel)
	Log.SetReportCaller(true)
}

type logFileWriter struct {
	file *os.File
	size int64
}

func (p *logFileWriter) Write(data []byte) (n int, err error) {
	if p == nil {
		return 0, errors.New("logFileWriter is nil")
	}
	if p.file == nil {
		return 0, errors.New("file not opened")
	}
	n, e := p.file.Write(data)
	p.size += int64(n)
	//每天一个文件
	if p.file.Name() != "applog/app-"+time.Now().Format("2006-01-02")+".log" {
		p.file.Close()
		p.file, _ = os.OpenFile("applog/app-"+time.Now().Format("2006-01-02")+".log", os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)
		p.size = 0
	}
	return n, e
}
