package services

import (
	"../logger"
	"io/ioutil"
	"os"
	"encoding/json"
	"os/signal"
	"syscall"
)

type ServiceLocator struct {
}

func NewServiceLocator() *ServiceLocator{
	return &ServiceLocator{}
}

func (self *ServiceLocator) LoadConfig(configPath string) (map[string]interface{}, error) {

	file, err := ioutil.ReadFile(configPath)

	if nil != err {

		return nil, err
	}

	config := map[string]interface{}{}

	err = json.Unmarshal(file, &config)

	if (nil != err) {
		return nil, err
	}

	return config, nil

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