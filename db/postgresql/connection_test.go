package postgresql

import (
	"github.com/udbjqrmna/banana/db/postgresql/pgproto3"
	"io"
	"net"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestConn(t *testing.T) {
	startupMsg := pgproto3.StartupMessage{
		ProtocolVersion: pgproto3.ProtocolVersionNumber,
		Parameters:      make(map[string]string),
	}

	startupMsg.Parameters["user"] = "postgres"
	startupMsg.Parameters["database"] = "test"
	startupMsg.Parameters["replication"] = "false"

	conn, err := net.Dial("tcp", "127.0.0.1:5432")
	if err != nil {
		log.Error().Error(err).Msg("连接出现错误。")
	} else {
		var startupMsg = startupMsg.Encode(nil)
		log.Trace().Bytes("Value", startupMsg).Msg("正常连接上了")
		conn.Write(startupMsg)
	}

	buf := make([]byte, 4096)
	//sql := pgproto3.Query{String: "begin",}
	//conn.Write(sql.Encode(nil))
	//sql = pgproto3.Query{String: "select * from test limit 1",}
	//conn.Write(sql.Encode(nil))
	//sql = pgproto3.Query{String: "insert into test(a, b, c) values(1,2,3);",}
	//conn.Write(sql.Encode(nil))
	//sql = pgproto3.Query{String: "commit",}
	//conn.Write(sql.Encode(nil))
	for {
		cnt, err := conn.Read(buf)
		if err != nil {
			log.Trace().Msgf("Fail to read data, %s\n", err)
			break
		}

		log.Trace().Bytes("Value", buf[0:cnt]).Msgf("收到消息：%s", buf[0:cnt])
	}
}

func TestConn2(t *testing.T) {
	startupMsg := pgproto3.StartupMessage{
		ProtocolVersion: pgproto3.ProtocolVersionNumber,
		Parameters:      make(map[string]string),
	}

	startupMsg.Parameters["user"] = "udbjqr"
	startupMsg.Parameters["database"] = "test"

	conn, err := net.Dial("tcp", "127.0.0.1:5432")
	if err != nil {
		log.Error().Error(err).Msg("连接出现错误。")
		return
	}

	buf := make([]byte, 4096)
	for {
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		cnt, err := conn.Read(buf)
		if err != nil {
			log.Trace().Msgf("Fail to read data, %T\n", err)
			if err == io.EOF {
				log.Trace().Msg("连接已经断开")
				return
			}
		}
		netErr, ok := err.(net.Error)
		if !ok {
			return
		}

		if netErr.Timeout() {
			log.Trace().Msg("timeout")
		}

		opErr, ok := netErr.(*net.OpError)
		if !ok {
			return
		}

		switch t := opErr.Err.(type) {
		case *net.DNSError:
			log.Trace().Msgf("net.DNSError:%+v", t)
			return
		case *os.SyscallError:
			log.Trace().Msgf("os.SyscallError:%+v", t)
			if errno, ok := t.Err.(syscall.Errno); ok {
				switch errno {
				case syscall.ECONNREFUSED:
					log.Trace().Msg("connect refused")
				case syscall.ETIMEDOUT:
					log.Trace().Msg("timeout1234")
				}
			}
		}

		log.Trace().String("make", "1234").Msg(string(buf[0:cnt]))
	}
}

func TestConn3(t *testing.T) {
	if conn, err := NewConnection("127.0.0.1:5432"); err != nil {
		log.Trace().Msg("错误")
	} else {
		time.Sleep(10 * time.Second)
		conn.Close()
	}

}
