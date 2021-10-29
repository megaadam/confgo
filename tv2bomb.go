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
	"strings"
	"strconv"
	"time"

	// positional args, cmd-completion, etc!
	// https://pkg.go.dev/github.com/jessevdk/go-flags
	// "github.com/jessevdk/go-flags"


	//log "github.com/sirupsen/logrus"

)

const (
	BASE_URL = "http://localhost:5000"
	SUB_URL = "/config/__active/services/liveIngest/channels/popup_%d/state"
	SUB_URL_CHAN = "/config/__active/services/liveIngest/channels/%s/state"
	SUB_URL_REFSYS = "/config/__active/services/liveIngest/channels/%s/state"
	CLI_ARG = "services.liveIngest.channels.popup_%d.state"
	CONFCLI = false
)

type Cfg struct {
    State   string      `json:"state"`
}

func getAllChannels() []string {
	all := []string {"play_1", "popup_14", "popup_8", "tv2_news",
	"play_1_25_no_scte", "popup_14_25_no_scte", "popup_8_25_no_scte", "tv2_news_25_no_scte",
	"play_2", "popup_15", "popup_9", "tv2_nord",
	"play_2_25_no_scte", "popup_15_25_no_scte", "popup_9_25_no_scte", "tv2_nord_25_no_scte",
	"play_3", "popup_2", "tv2_bornholm", "tv2_oest",
	"play_3_25_no_scte", "popup_2_25_no_scte", "tv2_bornholm_25_no_scte", "tv2_oest_25_no_scte",
	"popup_1", "popup_3", "tv2_charlie", "tv2_oestjylland",
	"popup_10", "popup_3_25_no_scte", "tv2_charlie_25_no_scte", "tv2_oestjylland_25_no_scte",
	"popup_10_25_no_scte", "popup_4", "tv2_fri", "tv2_sport",
	"popup_11", "popup_4_25_no_scte", "tv2_fri_25_no_scte", "tv2_sport_25_no_scte",
	"popup_11_25_no_scte", "popup_5", "tv2_fyn", "tv2_sport_x",
	"popup_12", "popup_5_25_no_scte", "tv2_fyn_25_no_scte", "tv2_sport_x_25_no_scte",
	"popup_12_25_no_scte", "popup_6", "tv2_lorry", "tv2_syd",
	"popup_1_25_no_scte", "popup_6_25_no_scte", "tv2_lorry_25_no_scte", "tv2_syd_25_no_scte",
	"popup_13", "popup_7", "tv2_midtvest", "tv2_zulu",
	"popup_13_25_no_scte", "popup_7_25_no_scte", "tv2_midtvest_25_no_scte", "tv2_zulu_25_no_scte",}

	return all
}

func getChannelList(count int) []string {
	var channelList []string

	for i := 1; i <= count; i++ {
		ch := fmt.Sprintf("polsat243-%d", i)
		channelList = append(channelList, ch)
	}

	return channelList
}



func getUrlsRefSys(count int) []string {
	var urls []string

	for i := 1; i <= count; i++ {
		ch := fmt.Sprintf("polsat243-%d", i)
		url := fmt.Sprintf(BASE_URL + SUB_URL_REFSYS, ch)

		urls = append(urls, url)
	}
	return urls
}



func channelsToCheck() []string {
	return []string {"popup_1", "popup_2", "popup_3", "popup_4", "popup_5", "popup_6", "popup_7", }
}

func chanToCheck(a string, channels []string) bool {
    for _, b := range channels {
        if b == a {
            return true
        }
    }
    return false
}

func main() {
	counts := os.Args[1]

	count, _ := strconv.Atoi(counts)

	if count == 0 {
		count = 10
	}

	channels := getChannelList(count)

	fmt.Println(channels)
	urls := getUrls(channels)
	//urls := getAllUrls() // All "known" channels
	args := getCliArgs()
	_ = urls

	iterCount := 0
	for {
		fmt.Println("\n\nIteration: ", iterCount)
		iterCount++
		dt := time.Now()
		fmt.Println("time: ", dt.Format("2006-01-02 15:04:05"))
		setChannels(urls, args, "catchup")
		checkConfCli("Catchup", channels)
		time.Sleep(time.Second * 10)


		fmt.Println("time: ", time.Now().Format("2006-01-02 15:04:05"))
		setChannels(urls, args, "enabled")
		checkConfCli("Enabled", channels)

		time.Sleep(15 * time.Second)
	}
}

func checkConfCli(expected string, channels []string) {
	cmd := exec.Command("ew-live-ingest-tool", "-l")
	stdout, err := cmd.Output()
	if(err != nil) {
		fmt.Println(err)
		return
	}

		fmt.Println("---")
	var allChannels []string
	allChannels = strings.Split(string(stdout), "\n")
	for _, channel := range allChannels {
		fields:= strings.FieldsFunc(channel, splitFn)

		if len(fields) > 1 && chanToCheck(fields[0], channels) && fields[1] != expected {
			fmt.Println("ERROR: ", fields[0], "Expected:", expected, "Actual:", fields[1])
			os.Exit(-1)

		}
	}


}

func splitFn(c rune) bool {
	return c == ' '
}

func getUrls(channels []string) []string {
	var urls []string
	for _, channel := range channels {
		url := BASE_URL + fmt.Sprintf(SUB_URL_CHAN, channel)
		urls = append(urls, url)
	}

	return urls
}

func getAllUrls() []string {
	var urls []string
	for _, chanx := range getAllChannels() {
		url := BASE_URL + fmt.Sprintf(SUB_URL_CHAN, chanx)
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
		fmt.Printf("setChannels() REST [%d channels]: %s", len(urls), state)

		for _, url := range urls  {
			wg.Add(1)
			go setChannel(url, state, &wg)
		}

	}

	wg.Wait()


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

	start := time.Now()
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

	fmt.Printf("Request duration: %s\n", time.Since(start))

}