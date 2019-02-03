package db

import (
	"errors"
	logD "github.com/udbjqrmna/banana/db/log"
)

const (
	EmptyString = ""        //EmptyString 一个空字符串的值，不是nil
	Default     = "default" //定义一个默认的名称
)

//定义此包的全局对象
var (
	log = logD.Log() //日志对象

	NotAllowOperation = errors.New("此操作不被允许")      //NotAllowOperation 不允许的操作
	pools             = make(map[string]*ConnPool) //保存池对象的map
)
