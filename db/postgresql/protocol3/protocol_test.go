package protocol3

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestEncodeStartupMessage(t *testing.T) {
	connectionUrl := `{"url":"127.0.0.1:5432","user":"postgres","password":"123456"}`
	para := make(map[string]string)
	if err := json.Unmarshal([]byte(connectionUrl), &para); err != nil {
		log.Error().Error(err).Msg("出现异常")
		return
	}

	log.Trace().Bytes("msg", EncodeStartupMessage(para)).Msg("生成的消息")
}

func TestTakeMessage(t *testing.T) {
	src := make([]byte, 200)[:0]
	src = append(src, "abce"...)
	src = append(src, "2"...)
	src = append(src, 0)
	src = append(src, "2fija1aabc"...)
	src = append(src, 0)
	src = append(src, "2fija"...)
	src = append(src, 0)

	s2 := src[198:199][:0]
	s2 = append(s2, 0)
	s2 = append(s2, '2')

	s, ind := takeContent(src, s2)
	fmt.Print(s + "\n")
	s, ind = takeContent(src[ind:], s2)
	fmt.Print(s)
}
