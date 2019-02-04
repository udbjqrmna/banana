package test

import (
	"fmt"
	. "github.com/udbjqrmna/banana/db"
	"github.com/udbjqrmna/banana/db/postgresql"
	"github.com/udbjqrmna/onelog/log"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func CreatePool(max, core uint8) *ConnPool {
	if pool, err := NewDefaultPool("连接字串", max, core, postgresql.NewConnection, 5); err != nil {
		log.Trace().Msg("出现异常：" + err.Error())
		return nil
	} else {
		return pool
	}
}

func TestGetPool(t *testing.T) {
	CreatePool(3, 2)

	log.Info().Msg(GetPool().Name)
}

func TestGetPoolByName(t *testing.T) {
	CreatePool(6, 2)
	if pool, err := GetPoolByName("abcdefg"); err != nil {
		log.Error().Error(err).Msg("未找到指定名称的对象")
	} else {
		log.Info().Msg(pool.Name)
	}
	if pool, err := GetPoolByName(Default); err != nil {
		log.Error().Error(err).Msg("未找到指定名称的对象")
	} else {
		log.Info().Msg(pool.Name)
	}
}

func TestPoolRun(t *testing.T) {
	CreatePool(10, 1)

	pool := GetPool()
	log.Trace().Msg("1")
	pool.GetConnect()

	log.Trace().Msg("2")
	pool.GetConnect()
	log.Trace().Msg("3")
	pool.GetConnect()
	log.Trace().Msg("4")
	pool.GetConnect()
	pool.GetConnect()
	pool.Close()

	//
	//
	//pool.ReturnConnection(conn3)
	//log.Trace().Msg("5")
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
	CreatePool(5, 2)
	pool := GetPool()
	g := sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		g.Add(1)
		go func() {
			//log.Trace().Msg(fmt.Sprintf("开始执行方法"))
			time.Sleep(time.Duration(rand.Intn(6)) * time.Second)
			conn := pool.GetConnect()
			log.Trace().Msgf("得到一个连接:%p", conn)

			time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
			log.Trace().Msgf("开始还回连接:%p", conn)
			pool.ReturnConnection(conn)
			g.Done()
		}()
	}
	time.Sleep(20 * time.Second)

	for i := 0; i < 200; i++ {
		g.Add(1)
		go func() {
			//log.Trace().Msg(fmt.Sprintf("开始执行方法"))
			time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
			conn := pool.GetConnect()
			log.Trace().Msg(fmt.Sprintf("得到一个连接:%p", conn))

			time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
			log.Trace().Msg(fmt.Sprintf("开始还回连接:%p", conn))
			pool.ReturnConnection(conn)
			g.Done()
		}()
	}

	g.Wait()
	pool.Close()

	time.Sleep(3 * time.Second)
}

func TestFree(t *testing.T) {
	pool := CreatePool(10, 2)
	g := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		g.Add(1)
		go func() {
			conn := pool.GetConnect()
			log.Trace().Msgf("得到一个连接:%p", conn)

			time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)
			log.Trace().Msgf("开始还回连接:%p", conn)
			pool.ReturnConnection(conn)
			g.Done()
		}()
	}

	g.Wait()

	time.Sleep(30 * time.Second)
	pool.Close()
	time.Sleep(2 * time.Second)
}
