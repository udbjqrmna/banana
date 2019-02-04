package protocol3

const (
	ProtocolVersionNumber = 0x00030000 // 3.0
	//sslRequestNumber      = 80877103

	url      = "url"
	user     = "user"
	password = "password"
	database = "database"
)

func EncodeStartupMessage(para map[string]string) []byte {
	dst := make([]byte, 64)[:0]

	dst = append(dst, 0, 0, 0, 0)
	dst = AppendInt32(dst, ProtocolVersionNumber)

	for k, v := range para {
		if !(k == url || k == password) {
			dst = append(dst, k...)
			dst = append(dst, 0)
			dst = append(dst, v...)
			dst = append(dst, 0)
		}
	}

	dst = append(dst, 0)

	SetInt(dst, len(dst))

	return dst
}
