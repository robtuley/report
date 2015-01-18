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

	go splunkStormForwarder(to.String(), accessKey)
}

func splunkStormForwarder(url string, accessKey string) {
	var wg sync.WaitGroup
	stopping := false

	client := &http.Client{}
	ticker := time.NewTicker(2 * time.Second)

	for {
		buffer := ""

	buffering:
		for {
			select {
			case json, more := <-channel.JsonEncoded:
				if !more {
					stopping = true
					break buffering
				}
				if len(buffer) > 0 {
					buffer += "\n"
				}
				buffer += json
			case <-ticker.C:
				if len(buffer) > 0 {
					break buffering
				}
			}
		}

		if len(buffer) > 0 {
			wg.Add(1)
			go func(data string) {
				defer wg.Done()

				req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(buffer)))
				if err != nil {
					log.Println("error:> ", err)
					return
				}
				req.SetBasicAuth("x", accessKey)
				req.Header.Set("Content-Type", "text/plain")

				resp, err := client.Do(req)
				if err != nil {
					log.Println("error:> ", err)
					return
				}

				decoder := json.NewDecoder(resp.Body)
				var msg map[string]interface{}
				err = decoder.Decode(&msg)
				if err != nil {
					log.Println("error:> ", err)
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
			ticker.Stop()
			wg.Wait()
			channel.Drain <- true
			return
		}
	}
}
