package test

import (
	"fmt"
	. "github.com/udbjqrmna/banana/db"
	"github.com/udbjqrmna/banana/db/postgresql"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestNewPool(t *testing.T) {
	if _, err := NewDefaultPool("", 10, 1, postgresql.CreateConnection); err != nil {
		Log().Trace().Msg("出现异常：" + err.Error())
	}
}

func TestGetPool(t *testing.T) {
	TestNewPool(nil)

	//Log().InfoMsg("23")
	Log().Info().Msg(GetPool().Name)
}

func TestPoolRun(t *testing.T) {
	TestNewPool(nil)

	pool := GetPool()
	Log().Trace().Msg("1")
	pool.GetConnect()

	Log().Trace().Msg("2")
	pool.GetConnect()
	Log().Trace().Msg("3")
	pool.GetConnect()
	Log().Trace().Msg("4")
	pool.GetConnect()
	pool.GetConnect()
	pool.Close()

	//
	//
	//pool.ReturnConnection(conn3)
	//Log().Trace().Msg("5")
	//conn5 :=pool.GetConnect()
	//
	//pool.ReturnConnection(conn1)
	//pool.ReturnConnection(conn2)
	//pool.ReturnConnection(conn5)
	//pool.ReturnConnection(conn4)
	//
	//time.Sleep(2 * time.Second)
	//
	//
	//
	//time.Sleep(2 * time.Second)
}

func TestPoolRun2(t *testing.T) {
	TestNewPool(nil)
	pool := GetPool()
	g := sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		g.Add(1)
		go func() {
			//Log().Trace().Msg(fmt.Sprintf("开始执行方法"))
			time.Sleep(time.Duration(rand.Intn(6)) * time.Second)
			conn := pool.GetConnect()
			Log().Trace().Msg(fmt.Sprintf("得到一个连接:%p", conn))

			time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
			Log().Trace().Msg(fmt.Sprintf("开始还回连接:%p", conn))
			pool.ReturnConnection(conn)
			g.Done()
		}()
	}
	time.Sleep(10 * time.Second)

	for i := 0; i < 200; i++ {
		g.Add(1)
		go func() {
			//Log().Trace().Msg(fmt.Sprintf("开始执行方法"))
			time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
			conn := pool.GetConnect()
			Log().Trace().Msg(fmt.Sprintf("得到一个连接:%p", conn))

			time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
			Log().Trace().Msg(fmt.Sprintf("开始还回连接:%p", conn))
			pool.ReturnConnection(conn)
			g.Done()
		}()
	}

	g.Wait()
	pool.Close()

	time.Sleep(3 * time.Second)
}
