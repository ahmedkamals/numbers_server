package main

import (
	"./app"
	"./services"
	"strconv"
)

const (
	MAX_WORKERS = 4
	MAX_QUEUE = 4
)

func main()  {

	serviceLocator := &services.ServiceLocator{}
	//logger := serviceLocator.Logger()

	config := serviceLocator.LoadConfig("config/main.json")
	baseConfig, backendTimeout := parseConfigs(&config)

	dispatcher := app.NewDispatcher(MAX_WORKERS, MAX_QUEUE, baseConfig, backendTimeout)
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