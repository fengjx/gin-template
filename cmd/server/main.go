package main

import (
	_ "gin-template/internal/app/bootstrap"
	"gin-template/internal/app/command"
)

func main() {
	command.Execute()
}
