package main

import (
	"context"

	"github.com/ShuaibKhan786/yours/cmd/yyd/backend"
	"github.com/ShuaibKhan786/yours/cmd/yyd/channel"
	"github.com/ShuaibKhan786/yours/cmd/yyd/gui"
)

func main() {
	//root context in which cancel will called when window close
	rootCtx, rootCancel := context.WithCancel(context.Background())

	c := channel.InitChannel()

	backend.InitBackend(rootCtx, c)

	yyd := gui.NewYYD(rootCtx, rootCancel, c)
	yyd.InitGUI()
}
