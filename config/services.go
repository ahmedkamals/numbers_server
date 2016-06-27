package config

import (
	"../logger"
	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"
	"os/signal"
	"syscall"
)

type ServiceLocator struct {
}

func (*ServiceLocator) LoadConfig(configPath string) map[string]interface{} {

	file, err := ioutil.ReadFile(configPath)

	if nil != err {
		// Todo: Use logger and be brave, don't exit.
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}

	config := map[string]interface{}{}

	err = json.Unmarshal(file, &config)

	if (nil != err) {
		// Todo: Use logger and stop freaking out.
		panic(err)
	}

	return config

}

func (*ServiceLocator) Logger() *logger.Log{
	return logger.NewLogger()
}

func (*ServiceLocator) BlockIndefinitely() {

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	println("Blocking indefinitely...")
	<-sigc
	println("Bye Bye!")
}