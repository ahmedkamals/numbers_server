package app

import (
	"flag"
	"net/http"
	"log"
	"net/url"
	"encoding/json"
	"strconv"
	httpProtocol "../communication/protocols/http"
	"../queue"
	"../processing"
	"../services"
)

type NumbersServer struct {
	config map[string]string
}

func NewNumbersServer(config map[string]string) *NumbersServer{
	return &NumbersServer{config}
}

func (self *NumbersServer) Start() {

	listenAddr := flag.String("http.addr", ":" + self.config["port"], "http listen address")
	flag.Parse()

	http.HandleFunc(self.config["path"], handler(self))

	aggregator := processing.NewAggregator()

	timeout, err := strconv.Atoi(self.config["timeout"])
	serviceLocator := services.NewServiceLocator()

	if nil != err {
		serviceLocator.Logger().Error(err.Error())
	}

	go aggregator.MonitorNewData(timeout)
	go aggregator.Aggregate(processing.NewOperator())

	// Should be last line as it is a blocking.
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

func handler(self *NumbersServer) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		url := r.URL.Query()

		jobs := self.getJobsFromUrl(url)

		queue.PushToChanel(queue.NewJobCollection(jobs))

		self.respond(w)
	}
}

func (*NumbersServer) buildJob(id, method, host, path string) *queue.Job{

	protocol := httpProtocol.NewProtocol(
		&http.Client{},
	)

	payload := queue.NewPayload(method, protocol, host, path)
	return queue.NewJob(id, payload)
}

func (self *NumbersServer) getJobsFromUrl(urlValues url.Values) []queue.Job {

	jobs := []queue.Job{}

	for index, item := range urlValues["u"] {

		urlScheme, err := url.Parse(item)

		if (nil != err) {

			log.Fatal(err.Error())
		}

		job := self.buildJob(strconv.Itoa(index), http.MethodGet, urlScheme.Host, urlScheme.Path)
		jobs = append(jobs, *job)
	}

	return jobs
}

func (*NumbersServer) respond(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Locking on Aggregated data
	numbers := <-processing.AggregationQueue

	serviceLocator := services.ServiceLocator{}
	serviceLocator.Logger().Info("finally", numbers)

	json.NewEncoder(w).Encode(map[string]interface{}{"Numbers": numbers})
}
