package postgresql

import (
	"github.com/udbjqrmna/banana/db"
)

type Connection struct {
	url string
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
	return nil, nil
}
