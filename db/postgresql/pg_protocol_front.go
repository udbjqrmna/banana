package postgresql

func passwdMessage(pwd []byte) []byte {
	dst := make([]byte, len(pwd)+6)
	dst = append(dst, 'p')
	dst = AppendInt32(dst, int32(4+len(pwd)+1))

	dst = append(dst, pwd...)
	dst = append(dst, 0)

	return dst
}
