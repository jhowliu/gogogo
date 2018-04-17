package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Body struct {
	Tokens []string `json:"segmentresult"`
}

type Job struct {
	Sentence string
	Tokens   []string
}

func buildRequest(method, host, path string) *http.Request {
	uri := host + path
	req, _ := http.NewRequest(method, uri, nil)

	return req
}

func queryEncoded(req *http.Request, sentence string) *http.Request {
	q := req.URL.Query()
	q.Add("q", sentence)

	req.URL.RawQuery = q.Encode()

	return req
}

func tokenize(language string, job *Job) {
	var req *http.Request

	if language == "tw" {
		req = buildRequest("GET", "http://192.168.10.108:3013", "/simplesegment")
	} else {
		req = buildRequest("GET", "http://localhost:3008", "/simplesegment")
	}

	req = queryEncoded(req, job.Sentence)

	response, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(err)
	}

	defer response.Body.Close()

	if response.StatusCode == 200 {
		var body Body
		_ = json.NewDecoder(response.Body).Decode(&body)
		fmt.Println(body.Tokens)
		job.Tokens = body.Tokens
	}
}

func dispatcher(language string, numOfWorkers int, jobs chan *Job) {
	var workers []chan struct{} = make([]chan struct{}, numOfWorkers)

	// running workers
	for i := 0; i < numOfWorkers; i++ {
		workers[i] = worker(language, jobs)
	}

	// wait for workers finished
	for i := 0; i < numOfWorkers; i++ {
		<-workers[i]
		fmt.Printf("Worker %d finished\n", i)
	}
}

func worker(language string, jobs chan *Job) chan struct{} {
	var end chan struct{} = make(chan struct{}, 1)
	go func() {
		for true {
			job, ok := <-jobs
			if !ok {
				break
			}

			tokenize(language, job)
		}
		end <- struct{}{}
	}()

	return end
}

func Tokenize(sentences []string, language string, numOfWorkers int) [][]string {
	count := len(sentences)

	var jobs []Job
	var results [][]string
	var jobChannel chan *Job = make(chan *Job, count)

	for _, s := range sentences {
		jobs = append(jobs, Job{Sentence: s})
	}

	for i, _ := range jobs {
		jobChannel <- &jobs[i]
	}

	close(jobChannel)

	start_t := time.Now()
	dispatcher(language, numOfWorkers, jobChannel)
	end_t := time.Now()

	fmt.Println("Tokenize for", count, "sentences takes", end_t.Sub(start_t))

	for _, job := range jobs {
		results = append(results, job.Tokens)
	}

	return results
}
