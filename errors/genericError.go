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
