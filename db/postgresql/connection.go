package postgresql

import (
	"fmt"
	"github.com/udbjqrmna/banana/db"
	"io"
	"net"
	"time"
)

const (
	inited        byte = 0x0  //对象构建完成，未进行任何动作
	tryConnection byte = 0x1  //正在尝试连接服务器当中
	connError     byte = 0xa  //连接出现无法继续的错误
	connected     byte = 0x10 //连接已经完成
	isReady       byte = 0x20 //已经被使用,此状态可进行操作
	querying      byte = 0x21 //sql语句请求正在发送，等待返回
	returning     byte = 0x25 //查询返回中，下一个查询还不能启动
	transaction   byte = 0x30 //事务已经启动
	tranquerying  byte = 0x31 //事务中的sql语句请求正在发送，等待返回
	tranReturning byte = 0x35 //事务中查询返回中，下一个查询还不能启动
	copyIning     byte = 0x40 //执行copy方法从客户端向服务器，待整理
	copyOuting    byte = 0x41 //执行copy方法从服务器向客户端，待整理
	closing       byte = 0xA0 //连接正在关闭中。
	closed        byte = 0xAF //连接已经被关闭
	unknownError  byte = 0xFF //连接已经被关闭
	//以下为连接参数使用的名称
	password = "password"
	user     = "user"
)

/*
Postgresql 的连接对象。
此对象将启动2个协程，
  一个接收服务器端发来的消息
  一个做为守护，处理从服务器与客户端发来的消息
客户端发送消息将使用调用者自己的协程进行处理
*/
type Connection struct {
	url      string
	connPara map[string]string //连接的一些参数，如连接地址，账号密码等等。
	status   byte

	conn net.Conn

	request    chan []byte
	response   chan []byte
	networkErr chan error
}

func (c *Connection) GetConnectUrl() string {
	return c.url
}

//创建一个连接的方法，此方法保证操作的哪一种类的数据库
func NewConnection(connectionUrl string) (db.Connection, error) {
	//TODO 创建一个连接的方法
	conn := &Connection{
		url:    connectionUrl,
		status: inited,
	}

	go conn.run()
	return conn, nil
}

//run　连接的协程方法，在此方法里实现接收双向来的信息并保证执行
func (c *Connection) run() {
	//有可能在执行过程当中出现异常后，重新开始连接
	for c.status == inited {
		if err := c.StartLink(); err != nil {
			log.Error().Error(err).Msg("无法连接对象,连接退出。")
			c.status = connError
			goto exitRun
		}

		go c.read() //启动接收服务器信息的协程

		for c.status > connected && c.status < closing {
			select {
			//收到客户端发来消息的处理
			case requestMsg := <-c.request:
				//做一个保护，如果一个空的消息发来。只进行记录
				if requestMsg != nil || len(requestMsg) <= 0 {
					log.Warn().Msg("丢掉一个收到的空消息。")
					break
				}

				if requestMsg[0] == connClose && len(requestMsg) == 1 {
					log.Info().Msg("收到关闭操作的消息。开始进行关闭")
					goto exitRun
				}

				if err := c.handleRequest(requestMsg); err != nil {
					c.status = inited //当发送消息出现异常时，将状态改回初始化，准备再次重新连接
				}

			//收到服务端发来消息的处理
			case repMsg := <-c.response:
				if repMsg != nil || len(repMsg) <= 0 {
					log.Warn().Msg("丢掉一个收到的服务端空消息。这个应该是不会出现的")
				}
				log.Trace().Msgf("收到服务端发来的消息。%v", repMsg)

				if err := c.handleResponse(repMsg); err != nil {
					log.Error().Error(err).Msg("收到一个无法操作的服务器消息。直接退出")
					c.status = unknownError
				}

			//收到网络异常时的处理
			case err := <-c.networkErr:
				log.Error().Error(err).Msgf("收到服务端接收的异常。")

				if err == io.EOF {
					c.status = inited
				}
				//TODO 收到其他服务器异常时的处理
			}
		}
	}

exitRun:
	c.destroy()
}

//接收服务器发来的消息的协程
func (c *Connection) read() {
	for c.status > connected && c.status < closing {
		if err := c.conn.SetReadDeadline(time.Now().Add(2 * time.Second)); err == nil {
			var buf = make([]byte, defaultBuf)
			if receiveCount, err := c.conn.Read(buf); err == nil {
				c.response <- buf[:receiveCount]
			} else {
				c.networkErr <- err
			}
		} else {
			c.networkErr <- err
			return
		}
	}
}

//开始启动连接。
func (c *Connection) StartLink() error {
	c.status = tryConnection

	// 开始连接，如果出现异常在过一段时间继续连接
	for c.status == tryConnection {
		conn, err := net.Dial(netProtocol, "127.0.0.1:5432")

		if err != nil {
			log.Warn().Error(err).Msg("连接出现异常。等待重试。")
			time.Sleep(connectionTryIdleTime * time.Second)
		} else {
			log.Debug().Msg("与postgresql网络连接成功，等待回应登录。")
			c.conn = conn
			c.status = connected
			return nil
		}
	}

	return db.NotReady("尝试请求中断，连接")
}

func (c *Connection) destroy() {
	//如果当前状态为网络异常。直接返回就好，不需要再做清理
	if c.status == connError {
		return
	}

	c.status = closing
	_ = c.conn.Close()
	c.status = closed
	//TODO　销毁的方法
	//TODO 事务需要在销毁时处理掉，如果有未结束的事务，应该调用rollback()方法
}

func (c *Connection) handleRequest(requestMsg []byte) error {
	log.Trace().Bytes("request", requestMsg).Msg("开始处理客户端消息")
	_, err := c.conn.Write(requestMsg)

	return err
}

func (c *Connection) handleResponse(responseMsg []byte) error {
	log.Trace().Bytes("rspMsg", responseMsg).Msg("开始处理服务端消息")

	switch responseMsg[0] {
	case AuthenticationKey:
		if err := c.handleAuthentication(responseMsg[1:]); err == nil {
			c.status = isReady
		} else {
			return err
		}
	}

	//TODO 取第一个值，判断类型，并且返回对应的结果
	return nil
}

/*
以下部分将是客户端调用的部分，与连接不在同一个协程内执行
*/

//返回当前连接是否已经在使用，当已经放入通道内准备给出时也将被赋值使用
func (c *Connection) Busy() bool {
	return c.status >= isReady
}

//向连接发送关闭消息
func (c *Connection) Close() {
	switch c.status {
	//这个时候是在进行连接的过程当中，还未连接上。可直接断开
	case tryConnection, connError:
		c.status = closed
	case closing, closed: //正在关闭或者已经关闭，不再发送关闭消息
		return
	}

	c.request <- ConnCloseMsg
	return
}

//TODO 以下方法需要实现
func (c *Connection) StopUse() {
	c.status = connected
	log.Debug().Msg(fmt.Sprintf("这个连接使用完了,%p", c))
	//TODO 事务在还回来后需要进行个判断，是否已经提交或者回滚
	return
}

//StartUse　开始使用，并进行标志
func (c *Connection) StartUse() error {
	if c.status >= isReady {
		return db.AlreadyInUse("连接")
	}

	if c.status < connected {
		return db.NotReady("连接")
	}

	log.Debug().Msgf("开始使用连接,%p", c)
	c.status = isReady
	return nil
}
