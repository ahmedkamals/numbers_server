package app

import (
	"flag"
	"net/http"
	"log"
	"net/url"
	"encoding/json"
	httpProtocol "../communication/protocols/http"
	"strconv"
)

type NumbersServer struct {
}

func NewNumbersServer() *NumbersServer{
	return &NumbersServer{}
}

func (self *NumbersServer) Start(config map[string]string, backendTimeout int) {

	listenAddr := flag.String("http.addr", ":" + config["port"], "http listen address")
	flag.Parse()

	http.HandleFunc(config["path"], handler(self))

	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

func handler(self *NumbersServer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		url := r.URL.Query()

		jobs := self.getJobsFromUrl(url)

		PushToChanel(NewJobCollection(jobs))

		self.respond(w)
	}
}

func (*NumbersServer) buildJob(id, method, host, path string) *Job{

	protocol := httpProtocol.NewProtocol(
		&http.Client{},
	)

	payload := NewPayload(method, protocol, host, path)
	return NewJob(id, payload)
}

func (*NumbersServer) respond(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]interface{}{"Numbers": []int{1, 2, 3}})
}

func (self *NumbersServer) getJobsFromUrl(urlValues url.Values) []Job {

	jobs := []Job{}

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
