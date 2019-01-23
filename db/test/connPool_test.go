package test

import (
	. "github.com/udbjqrmna/banana/db"
	"github.com/udbjqrmna/banana/db/postgresql"
	"testing"
)

func TestNewPool(t *testing.T) {
	NewDefaultPool("", 20, 3, postgresql.CreateConnection)
}

func TestGetPool(t *testing.T) {
	NewDefaultPool("", 20, 10, postgresql.CreateConnection)

	//Log().InfoMsg("23")
	Log().InfoMsg(GetPool().Name)
}
