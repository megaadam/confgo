package main

import (
	//"flags" // use go-flags!!
	//"fmt"
	"bytes"
	"os/exec"
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
	// "github.com/jessevdk/go-flags"


	//log "github.com/sirupsen/logrus"

)

const (
	BASE_URL = "http://localhost:5000"
	SUB_URL = "/config/__active/services/liveIngest/channels/popup_%d/state"
	CLI_ARG = "services.liveIngest.channels.popup_%d.state"
	CONFCLI = false
)

type Cfg struct {
    State   string      `json:"state"`
}

func main() {
	urls := getUrls()
	args := getCliArgs()
	_ = urls
	for {
		dt := time.Now()
		fmt.Println("time: ", dt.Format("2006-01-02 15:04:05"))
		setChannels(urls, args, "disabled")
		// time.Sleep(15 * time.Second)
		time.Sleep(time.Second * 10)
		//checkChannels(urls, "enabled")


		fmt.Println("time: ", time.Now().Format("2006-01-02 15:04:05"))
		setChannels(urls, args, "catchup")
		// time.Sleep(15 * time.Second)
		time.Sleep(time.Millisecond * 1500)
		os.Exit(0)
		//checkChannels(urls, "catchup")

	}
}

func getUrls() []string {
	var urls []string
	for i := 7; i >= 1; i-- {
		url := BASE_URL + fmt.Sprintf(SUB_URL, i)
		urls = append(urls, url)
	}

	return urls
}

func getCliArgs() []string {
	var args []string
	for i := 1; i <= 7; i++ {
		arg := fmt.Sprintf(CLI_ARG, i)
		args = append(args, arg)
	}

	return args
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

func setChannels(urls, cliArgs []string, state string) {
	var wg sync.WaitGroup


	if CONFCLI {
		fmt.Print("setChannels() [confcli]: ", state)

		for _, arg := range cliArgs  {
			wg.Add(1)
			go setChannel(arg, state, &wg)
		}
	} else {
		fmt.Print("setChannels() REST: ", state)

		for _, url := range urls  {
			wg.Add(1)
			go setChannel(url, state, &wg)
		}

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




func setChannel(urlArg, state string, wg *sync.WaitGroup) {
	defer wg.Done()

	if CONFCLI {
		cmd := exec.Command("confcli", urlArg, state)
		stdout, err := cmd.Output()
		_ = err
		_ = stdout
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Print(".")
		}
	} else {
		data := &Cfg {
			State: state}

		buff, err := json.Marshal(data)
		req, err := http.NewRequest(http.MethodPut, urlArg, bytes.NewBuffer([]byte(buff)))
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

}