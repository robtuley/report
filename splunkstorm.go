package report

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

func SplunkStorm(apiUrl string, projectId string, accessKey string) {
	log.Println("reporting:> splunkstorm")
	log.Println("url:> ", apiUrl)
	log.Println("project:> ", projectId)
	go splunkStormForwarder(apiUrl, projectId, accessKey)
}

func splunkStormForwarder(apiUrl string, projectId string, accessKey string) {
	var wg sync.WaitGroup
	stopping := false

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
		buffer = bytes.NewBuffer(nil)

		go func() {
			time.Sleep(2 * time.Second)
			timeout <- true
		}()

	buffering:
		for {
			select {
			case json, more := <-channel.JsonEncoded:
				if !more {
					stopping = true
					break buffering
				}
				buffer.WriteString(json)
				buffer.WriteString("\n")
			case <-timeout:
				break buffering
			}
		}

		if buffer.Len() > 0 {
			wg.Add(1)
			go func(data *bytes.Buffer) {
				defer wg.Done()

				req, err := http.NewRequest("POST", to.String(), buffer)
				req.SetBasicAuth("x", accessKey)
				req.Header.Set("Content-Type", "text/plain")

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

		if stopping {
			wg.Wait()
			channel.Drain <- true
			return
		}
	}
}
