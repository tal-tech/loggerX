package stackerr

import (
	"fmt"
	"strings"

	"github.com/spf13/cast"
)

type StackErr struct {
	Filename   string
	Line       int
	Message    string //标准输出报错信息
	StackTrace string
	Code       int    //错误码
	Info       string //错误详情
	Position   string
	Level      int //0最高优先级 1-4 普通优先级 5 可不关注的异常
}

func (this *StackErr) ErrorInfo() string {
	return this.Info
}

func (this *StackErr) Error() string {
	return fmt.Sprintf("%d|%s", this.Code, this.Message)
}

func (this *StackErr) Stack() string {
	return fmt.Sprintf("(%s:%d)%s\tStack: %s", this.Filename, this.Line, this.Info, this.StackTrace)
}

func (this *StackErr) Detail() string {
	return fmt.Sprintf("(%s:%d)%s", this.Filename, this.Line, this.Info)
}

func (this *StackErr) Format(tag ...string) (data string) {
	var strs []string
	strs = append(strs, cast.ToString(this.Code))
	strs = append(strs, this.Message)
	strs = append(strs, this.Filename)
	strs = append(strs, cast.ToString(this.Line))
	strs = append(strs, this.Info)
	data = strings.Join(strs, "\t")
	return
}

func (this *StackErr) SetLevel(lvl int) {
	this.Level = lvl
}

func (this *StackErr) GetLevel() int {
	return this.Level
}
