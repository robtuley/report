package report

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"
)

func init() {
	go stdout(jsonEventChannel)
}

func stdout(in chan string) {
	for {
		log.Println("json:>", <-in)
	}
}

func splunkStorm(in chan string, apiUrl string, projectId string, accessKey string) {
	log.Println("reporting:> splunkstorm")
	log.Println("url:> ", apiUrl)
	log.Println("project:> ", projectId)

	var to *url.URL
	to, err := url.Parse(apiUrl)
	if err != nil {
		log.Println("error:> ", err)
		return
	}
	params := url.Values{}
	params.Add("project", projectId)
	params.Add("sourcetype", "json_auto_timestamp")
	to.RawQuery = params.Encode()

	client := &http.Client{}
	timeout := make(chan bool, 1)
	var buffer *bytes.Buffer

	for {
		buffer = bytes.NewBuffer([]byte(<-in))

		go func() {
			time.Sleep(2 * time.Second)
			timeout <- true
		}()

	buffering:
		for {
			select {
			case json := <-in:
				buffer.WriteString("\n")
				buffer.WriteString(json)
			case <-timeout:
				break buffering
			}
		}

		go func(data *bytes.Buffer) {
			req, err := http.NewRequest("POST", to.String(), buffer)
			req.SetBasicAuth("x", accessKey)
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				log.Println("error:> ", err.Error())
				return
			}

			decoder := json.NewDecoder(resp.Body)
			var msg map[string]interface{}
			err = decoder.Decode(&msg)
			if err != nil {
				log.Println("error:> ", err.Error())
				return
			}

			_, ok := msg["bytes"]
			if ok {
				log.Println("sent:> ", msg["bytes"])
			} else {
				log.Println("error:>", msg)
			}

		}(buffer)
	}
}
