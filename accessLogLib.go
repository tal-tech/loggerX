package logger

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/tal-tech/loggerX/builders"
	"go.uber.org/zap"
)

//global map
//used to register the other log lib
var accessLogMap = make(map[string]func(args interface{}), 0)

//support logrus and zap
func init() {
	accessLogMap["logrus"] = loadLogrus
	accessLogMap["zap"] = loadZap
}

/*

方法：AccessLogLib 接入loglib
功能：接入多样logLib 以支持其他组件中 loggerX 日志打印格式
@param libName 接入logLib库名 默认接入log4go
@param args[0] logrus全局logger实例 使用默认的logrus Builder
@param args[0] zap全局logger实例 使用默认的zap Builder
*/

func AccessLogLib(libName string, args interface{}) {
	if ok, fn := accessLogMap[libName]; ok {
		fn(args)
	}
}

//load logrus
func loadLogrus(args interface{}) {
	instance, ok := args.(*logrus.Logger)
	if !ok {
		fmt.Fprintf(os.Stderr, "SetLogLib Error: Could not get logrus instance!")
		os.Exit(1)
	}
	//logrus builder
	libBuild := builders.NewLogrusBuilder(instance)
	SetBuilder(libBuild)
}

//load zap
func loadZap(args interface{}) {
	instance, ok := args.(*zap.Logger)
	if !ok {
		fmt.Fprintf(os.Stderr, "SetLogLib Error: Could not get zap instance!")
		os.Exit(1)
	}
	//zap builder
	libBuild := builders.NewZapBuilder(instance)
	SetBuilder(libBuild)
}
