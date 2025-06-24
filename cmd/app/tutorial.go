package app

import (
	"golang-tutorial/config"
	"golang-tutorial/internal/server"
)

func Start() {
	config.Load()
	server.Run()
}
