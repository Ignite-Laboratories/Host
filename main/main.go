package main

import (
	"fmt"
	"github.com/ignite-laboratories/core"
	"github.com/ignite-laboratories/host"
)

func main() {
	for core.Alive {
		fmt.Println(host.Mouse.GetCoordinates())
	}
}
