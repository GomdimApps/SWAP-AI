package main

import (
	"github.com/GomdimApps/lcme"
	"github.com/NightGhost404/SWAP_AI/App/server"
)

func main() {
	lcme.ScaleFork(func() {
		server.StartServer()
	})
}
