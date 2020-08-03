package logger

import (
	"errors"
	"fmt"
	"testing"
)

func TestNewError(t *testing.T) {
	e := NewError("error reason 1", PARAM_EMAIL)
	fmt.Println("error:", e.Error()) //
	fmt.Println("errInfo", e.ErrorInfo())
	fmt.Println("detail", e.Detail())
	fmt.Println("format", e.Format())
	fmt.Println("stack", e.Stack())
}

func TestNewErrorWithError(t *testing.T) {
	err := errors.New("this is error reason")
	e := NewError(err, PARAM_EMAIL)
	fmt.Println("error:", e.Error())      //10400|邮箱错误
	fmt.Println("errInfo", e.ErrorInfo()) //this is error reason
	fmt.Println("detail", e.Detail())     //(error_test.go:20)this is error reason
	fmt.Println("format", e.Format())     //10400    邮箱错误        error_test.go   20      this is error reason
	fmt.Println("stack", e.Stack())
}

func TestNewErrorRpcx(t *testing.T) {
	//rpc server
	serverE := NewError("mysql error", PARAM_EMAIL)
	err := errors.New(serverE.Error())
	//rpc client revieve err
	e := NewError(err) //PARAM_EMAIL
	fmt.Println("error:", e.Error())
	fmt.Println("errInfo", e.ErrorInfo())
	fmt.Println("detail", e.Detail())
	fmt.Println("format", e.Format())
	fmt.Println("stack", e.Stack())
}

func TestNewErrorRpcx2(t *testing.T) {
	//rpc client call error
	err := errors.New("network eof")
	//rpc client revieve err
	e := NewError(err) //SYSTEM_DEFAUT
	fmt.Println("error:", e.Error())
	fmt.Println("errInfo", e.ErrorInfo())
	fmt.Println("detail", e.Detail())
	fmt.Println("format", e.Format())
	fmt.Println("stack", e.Stack())
}
