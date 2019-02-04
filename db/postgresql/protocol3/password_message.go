package protocol3

type PasswordMessage struct {
	frontMessage
	Password string
}

func EncodePassword(passwd string) []byte {
	buf := make([]byte, 6+len(passwd))[:0]

	buf = append(buf, 'p')
	buf = AppendInt(buf, 5+len(passwd))
	buf = append(buf, passwd...)
	buf = append(buf, 0)

	return buf
}
