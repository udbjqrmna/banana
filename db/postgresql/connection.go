package postgresql

import (
	"fmt"
	"github.com/udbjqrmna/banana/db"
)

type Connection struct {
	url   string
	busy  bool
	count int
}

//create 创建一个新连接并返回
func (c *Connection) New(connectionUrl string, h db.HandleConnected) (db.Connection, error) {
	//TODO 创建连接
	var newC db.Connection = &Connection{url: connectionUrl}

	if h != nil {
		if err := h(newC); err != nil {
			return nil, err
		}
	}

	return newC, nil
}

func (c *Connection) GetConnectUrl() string {
	return c.url
}

//创建一个连接的方法，此方法保证操作的哪一种类的数据库
func CreateConnection(connectionUrl string) (db.Connection, error) {
	//TODO 创建一个连接的方法
	return &Connection{
		url:  connectionUrl,
		busy: false,
	}, nil
}

//TODO 以下方法需要实现
func (c *Connection) StopUse() {
	c.busy = false
	c.count++
	log.Debug().Msg(fmt.Sprintf("这个连接使用完了,%p,使用了%v次", c, c.count))
	//TODO
	return
}
func (c *Connection) StartUse() {
	//TODO
	log.Debug().Msg(fmt.Sprintf("开始使用连接,%p", c))
	c.busy = true
	return
}
func (c *Connection) Busy() bool {
	//TODO
	return c.busy
}
func (c *Connection) Close() {
	//TODO
	return
}
