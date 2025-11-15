package main

import (
	internal "auth-service/internal/application"
)

func main() {
	server := internal.CreateServer()
	server.MountRoutes()
	server.Start()
}
