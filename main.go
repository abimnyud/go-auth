package main

import (
	"abimanyu.dev/go-auth/model"
	"abimanyu.dev/go-auth/routes"
)

func main() {
	model.Setup()
	routes.Setup()
}
