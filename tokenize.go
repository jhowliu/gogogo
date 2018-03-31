package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

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

func httpGet(sentence string, c chan string) {
	var result string
	var response *http.Response

	req := buildRequest("GET", "https://lingbot-api.lingtelli.com", "/segment/simplesegment")
	req = queryEncoded(req, sentence)

	response, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println("ERROR")
		result = "ERROR"
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)
		result = string(body)
	} else {
		result = "ERROR"
	}
	c <- string(result)
}

func asyncHttpGet(sentences []string) []string {
	responses := []string{}

	ch := make(chan string, len(sentences))

	for _, s := range sentences {
		go httpGet(s, ch)
	}

	for {
		select {
		case body := <-ch:
			responses = append(responses, body)
			if len(responses) == len(sentences) {
				return responses
			}
		}
	}
}

func main() {
	sentences := []string{
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
		"今天吃牛排嗎",
	}
	fmt.Println("START TOKENIZING")
	start_t := time.Now()
	results := asyncHttpGet(sentences)
	end_t := time.Now()

	for _, ele := range results {
		fmt.Println(ele)
	}

	fmt.Println("Tokenize for", len(results), "sentences takes", end_t.Sub(start_t))
}
