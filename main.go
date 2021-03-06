package main

import (
	"strconv"
	"os"
	"./app"
	"./services"
)

const (
	MAX_WORKERS = 4
	MAX_QUEUE = 4
)

func main()  {

	serviceLocator := services.NewServiceLocator()

	config, err := serviceLocator.LoadConfig("config/main.json")

	if (nil != err) {
		serviceLocator.Logger().Critical(err.Error())
		os.Exit(0)
	}

	baseConfig, backendTimeout := parseConfigs(&config)
	baseConfig["timeout"] = strconv.Itoa(backendTimeout)

	dispatcher := app.NewDispatcher(MAX_WORKERS, MAX_QUEUE, baseConfig)
	dispatcher.Run()

	serviceLocator.BlockIndefinitely()
}


func parseConfigs(config *map[string]interface{}) (map[string]string, int){

	baseConfig := map[string]string{}
	var backendTimeout int

	for key, value := range *config {
		switch key {
		case "base":
			data := value.(map[string]interface{})
			for k, base := range data {

				switch base.(type) {
				case float64:
					baseConfig[k] = strconv.Itoa(int(base.(float64)))
				case string:
					baseConfig[k] = base.(string)
					break
				}
			}
			break
		case "backends":
			data := value.(map[string]interface{})
			backendTimeout = int(data["timeout"].(float64))
			break
		}
	}

	return baseConfig, backendTimeout
}