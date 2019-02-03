package db

//AlreadyInUse 指明已经在使用
type AlreadyInUse string

func (e AlreadyInUse) Error() string {
	return string(e) + "已经在使用"
}

//NotReady 还未就绪，当下无法使用
type NotReady string

func (e NotReady) Error() string {
	return string(e) + "还未就绪"
}

//Unknown 未知的错误
type Unknown string

func (e Unknown) Error() string {
	return "未知的错误：" + string(e)
}
