package main

import (
	"fmt"
	"strconv"
)

type dd struct {
	d  int    `dab`
	ar string `ccc`
}

func (d dd) String() string {
	return "ar" + d.ar + "d" + strconv.Itoa(d.d)
}

func main() {
	fmt.Printf(dd{}.String())
}
