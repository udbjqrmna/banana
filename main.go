package main

import (
	"fmt"
	"github.com/udbjqrmna/banana/db"
)

func main() {
	var us = db.GetDBConnPool()
	us.Db = "abc"
	fmt.Println(us)
}
