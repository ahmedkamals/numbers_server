package app

import (
	"flag"
	"net/http"
	"log"
	"net/url"
	"encoding/json"
	httpProtocol "../communication/protocols/http"
	"strconv"
	"fmt"
	"time"
)

var (
	// A buffered channel that we can send fetched numbers on.
	MergeQueue chan []int
	MergedData []int
	IsMergeDone chan bool
)

type NumbersServer struct {
}

func NewNumbersServer() *NumbersServer{
	return &NumbersServer{}
}

func (self *NumbersServer) Start(config map[string]string) {

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

func (self *NumbersServer)startMergeChanel(timeout int) {

	go func() {
		for {
			select {
			case items := <-MergeQueue:
				MergedData = append(MergedData, items...)

			// Giving extra 100ms for processing
			case <- time.After(time.Millisecond * time.Duration(timeout - 100)):
				fmt.Println("timed out", MergedData)
				IsMergeDone <- true
				return
			}
		}
	}()
}

func (*NumbersServer) respond(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Locking on Aggregated data
	numbers := <-AggregatedData
	fmt.Println("finally", numbers)
	json.NewEncoder(w).Encode(map[string]interface{}{"Numbers": numbers})
}
