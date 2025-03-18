package main

import (
	"github.com/GomdimApps/lcme"
	"github.com/NightGhost404/SWAP_AI/App/files"
	"github.com/NightGhost404/SWAP_AI/App/server"
)

func main() {
	lcme.ScaleFork(func() {
		// Inicia o servidor
		server.StartServer()
	})
	// Inicia a rotina de limpeza de arquivos antigos
	files.StartCleanupRoutine()
}
