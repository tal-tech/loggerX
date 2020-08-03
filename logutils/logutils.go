package logutils

import (
	"strings"

	"github.com/tal-tech/xtools/log4go"
)

var DefaultReplacer *strings.Replacer

func init() {
	DefaultReplacer = strings.NewReplacer("\t", "", "\r", "", "\n", "")
	initLevelMap()
}

func Filter(msg string, r ...string) string {
	replacer := DefaultReplacer
	if len(r) > 0 {
		replacer = strings.NewReplacer("\t", r[0], "\r", r[0], "\n", r[0])
	}
	return replacer.Replace(msg)
}

var Level string = "ERROR"
var SortLevel log4go.Level = 7

var LevelMap map[string]log4go.Level
var Inited bool

//初始化日志级别map
func initLevelMap() {
	LevelMap = make(map[string]log4go.Level, 0)
	LevelMap["FINEST"] = log4go.FINEST
	LevelMap["FINE"] = log4go.FINE
	LevelMap["DEBUG"] = log4go.DEBUG
	LevelMap["INFO"] = log4go.INFO
	LevelMap["WARNING"] = log4go.WARNING
	LevelMap["ERROR"] = log4go.ERROR
	LevelMap["TRACE"] = log4go.TRACE
	LevelMap["CRITICAL"] = log4go.CRITICAL
	LevelMap["FATAL"] = log4go.CRITICAL
}

func ValidLevel(lvl string) bool {
	if !Inited {
		return true
	}
	if l, ok := LevelMap[Level]; !ok {
		return false
	} else if l > LevelMap[lvl] {
		return false
	}
	return true
}

func LogLevel(lvl string) log4go.Level {
	if !Inited {
		return log4go.FINEST
	}

	l, ok := LevelMap[lvl]
	if !ok {
		return log4go.FINEST
	}
	return l
}
