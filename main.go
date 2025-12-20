package main

import (
	"distributed-configuration-management/application"
	"distributed-configuration-management/shared/driver"
	"distributed-configuration-management/shared/infrastructure/config"
	"flag"
	"fmt"
)

func main() {

	// test every config is ready
	config.ReadConfig()

	appMap := map[string]func() driver.RegistryContract{
		"appcontroller": application.NewAppController(),
		"appworker":     application.NewAppWorker(),
		"appagent":      application.NewAppAgent(),
	}
	flag.Parse()

	app, exist := appMap[flag.Arg(0)]
	if exist {
		driver.Run(app())
	} else {
		fmt.Println("You may try 'go run main.go <app_name>' :")
		for appName := range appMap {
			fmt.Printf(" - %s\n", appName)
		}
	}

}
