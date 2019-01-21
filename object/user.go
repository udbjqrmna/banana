package object

/*
用户的操作对象，此对象对应业务功能当中的用户对象。
此对象提供包括账号密码验证、操作权限、查看权限
*/
type User struct {
	name     string
	password string
}

//Name 返回用户名称
func (u *User) Name() string {
	if u == nil {
		return EmptyString
	}

	return u.name
}
