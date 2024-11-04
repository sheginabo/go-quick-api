package main

import (
	initModule "github.com/sheginabo/go-quick-api/init"
)

func main() {
	initProcess := initModule.NewMainInitProcess("./")
	initProcess.Run()
}
