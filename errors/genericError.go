package errors

//NotNil 指明不能为nil异常
type NotNil string

func (e NotNil) Error() string {
	return string(e) + "不能为空"
}

//NotUnderstand 不被理解的值的异常
type NotUnderstand string

func (e NotUnderstand) Error() string {
	return "不能理解的值:" + string(e)
}

//IsFull 已满
type IsFull string

func (e IsFull) Error() string {
	return string(e) + " 已满。无法再增加"
}

//UnableToConnect 无法连接服务器
type UnableToConnect string

func (e UnableToConnect) Error() string {
	return "无法连接服务器，连接信息：" + string(e)
}

//UnknownMistake 未知错误
type UnknownMistake string

func (e UnknownMistake) Error() string {
	return "未知错误。 " + string(e)
}

//AlreadyExisting 已经存在
type AlreadyExisting string

func (e AlreadyExisting) Error() string {
	return "错误，指定的定义已经存在: " + string(e)
}

//MisTake 一个错误
type MisTake string

func (e MisTake) Error() string {
	return "出现了一个错误：" + string(e)
}
