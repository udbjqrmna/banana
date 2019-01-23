package main

import (
	"github.com/udbjqrmna/banana/db"
	"github.com/udbjqrmna/banana/db/postgresql"
	"github.com/udbjqrmna/onelog"
	"os"
)

var log = onelog.New(&onelog.Stdout{Writer: os.Stdout}, onelog.TraceLevel, &onelog.JsonPattern{}).AddRuntime(&onelog.Caller{})

func main() {
	var c db.Connection = &postgresql.Connection{}
	c.New("abc", tempH)
}

var tempH = func(c db.Connection) error {
	log.Debug().String("返回的值", c.GetConnectUrl()).Msg("这是一个测试")
	return nil
}
