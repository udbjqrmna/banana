package db

import "errors"

const (
	//EmptyString 一个空字符串的值，不是nil
	EmptyString string = ""
	//定义一个默认的名称
	Default string = "default"
)

//定义一些异常
var (
	//NotAllowOperation 不允许的操作
	NotAllowOperation = errors.New("此操作不被允许")
)
