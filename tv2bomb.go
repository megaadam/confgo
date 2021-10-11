package main

import (
	//"flags" // use go-flags!!
	//"fmt"
	"bytes"
	"io"
	"net/http"
	//"os"
	"fmt"
	"encoding/json"

	"os"
	"sync"
	"time"

	// positional args, cmd-completion, etc!
	// https://pkg.go.dev/github.com/jessevdk/go-flags

	//log "github.com/sirupsen/logrus"

)

const (
	BASE_URL = "http://adam-3003-1-24-2:5000"
	SUB_URL = "/config/__active/services/liveIngest/channels/popup_%d/state"
)

type Cfg struct {
    State   string      `json:"state"`
}

func main() {
	urls := getUrls()
	_ = urls
	for {
		dt := time.Now()
		fmt.Println("Current date and time is: ", dt.Format("2006-01-02 15:04:05"))
		setChannels(urls, "enabled")
		// time.Sleep(15 * time.Second)
		checkChannels(urls, "enabled")

		setChannels(urls, "catchup")
		// time.Sleep(15 * time.Second)
		checkChannels(urls, "catchup")

	}
}

func getUrls() []string {
	var urls []string
	for i := 1; i <= 7; i++ {
		url := BASE_URL + fmt.Sprintf(SUB_URL, i)
		urls = append(urls, url)
	}

	return urls
}

func checkChannels(urls []string, expected string) {
	var wg sync.WaitGroup


	fmt.Print("checkChannels() ", expected)
	for _, url := range urls  {
		wg.Add(1)
		go checkChannel(url, expected, &wg)
	}
	wg.Wait()
	fmt.Print("\n")
}

func setChannels(urls []string, state string) {
	var wg sync.WaitGroup

	fmt.Print("setChannels() ", state)
	for _, url := range urls  {
		wg.Add(1)
		go setChannel(url, state, &wg)



	}

	wg.Wait()
	fmt.Print("\n")
}

func checkChannel(url, expected string, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	fmt.Print(".")
	body, err := io.ReadAll(resp.Body)

	if err != nil {

		fmt.Println(err)
	}

	var dat map[string]interface{}
	err = json.Unmarshal(body, &dat)
	state := dat["state"]
	fail := false
	if(state != expected) {
		fmt.Println("\n\n URL:", url, state, " != ", expected)
		fail = true
	}

	if fail {
		os.Exit(-1)
	}
}




func setChannel(url, state string, wg *sync.WaitGroup) {
	defer wg.Done()

	data := &Cfg {
		State: state}

	buff, err := json.Marshal(data)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer([]byte(buff)))
    if err != nil {
        panic(err)
    }

    // set the request header Content-Type for json
    req.Header.Set("Content-Type", "application/json")
    // initialize http client
    client := &http.Client{}
    _, err = client.Do(req)
    if err != nil {
        panic(err)
    }
	fmt.Print(".")

}