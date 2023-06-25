package middleware

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"sf-paas/config"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var log = Logger
var Logger *logrus.Logger

const (
	LOG_PATH       = "log"
	LOG_LEVEL_INFO = "info"
	TIME_FORMAT    = "2006-01-02 15:04:05"
)

func init() {
	Logger = logrus.New()

	logFilePath := LOG_PATH
	logFileName := config.Conf.Log.Logfile
	fmt.Printf("log file name:%v\n", logFileName)
	if !IsDir(logFilePath) {
		if !CreateDir(logFilePath) {
			fmt.Printf("create log dir fialed")
			return
		}
	}
	// 日志文件
	fileName := path.Join(logFilePath, logFileName)
	fmt.Printf("file name:%v\n", fileName)
	//输出日志中添加文件名和方法信息
	Logger.SetReportCaller(true)
	// 写入文件
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Printf("write log file failed:%v\n", err)
		return
	}
	// 设置输出
	Logger.SetOutput(bufio.NewWriter(file))

	if config.Conf.Log.Loglevel == LOG_LEVEL_INFO {
		// 设置日志级别

		Logger.SetLevel(logrus.InfoLevel)
	} else {
		Logger.SetLevel(logrus.DebugLevel)
	}

	writeFile(fileName)
	fmt.Printf("init log success\n")
}

func writeFile(fileName string) {
	// 设置 rotatelogs
	logWriter, _ := rotatelogs.New(
		// 分割后的文件名称
		fileName+".%Y%m%d%H",

		// 生成软链，指向最新日志文件Loglevel
		rotatelogs.WithLinkName(fileName),

		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(config.Conf.Log.Logmaxage*24*time.Hour),

		// 设置日志切割时间间隔(1小时)
		rotatelogs.WithRotationTime(time.Hour),
	)

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}

	lfHook := lfshook.NewHook(writeMap, &nested.Formatter{
		HideKeys:        true,
		NoFieldsColors:  false,
		CallerFirst:     false,
		TrimMessages:    true,
		TimestampFormat: TIME_FORMAT,
	})

	// 新增 Hook
	Logger.AddHook(lfHook)
}

func Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		uri := c.Request.RequestURI
		//开始时间
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		method := c.Request.Method
		statusCode := c.Writer.Status()
		ip := c.ClientIP()
		if config.Conf.Log.Loglevel == LOG_LEVEL_INFO {
			Logger.WithFields(logrus.Fields{
				//客户端ip
				"clientIp": ip,
				//状态码
				"statusCode": statusCode,
				//接口请求方法
				"reqMethod": method,
				//请求接口
				"reqUri": uri,
				//请求耗时
				"latencyTime": latencyTime,
			}).Info()
		} else {
			now := time.Now().Format(TIME_FORMAT)
			Logger.Infof("%s | %3d | %13v | %15s | %s  %s",
				now,
				statusCode,
				latencyTime,
				ip,
				method,
				uri,
			)
		}
	}
}
