package main

import (
	"fmt"
	"valida/application/config"
	"valida/application/global"
	"valida/infrastructure/api"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Infoln(fmt.Sprintf("Application Version : %s", global.BUILD_VERSION))

	config.Init()

	// forever := make(chan bool)

	// go func() {
	// 	//Call Service : For example api.Serve()
	// 	api.Serve()
	// }()

	// <-forever

	api.Serve()
}
