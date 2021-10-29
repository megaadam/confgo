package main

import (
	//"flags" // use go-flags!!
	"bytes"
	"fmt"
	//"errors"
	"io"
	"net/http"
	"os"
	"strings"

	// "path"
	// "time"

	// positional args, cmd-completion, etc!
	// https://pkg.go.dev/github.com/jessevdk/go-flags
	"github.com/jessevdk/go-flags"

	//log "github.com/sirupsen/logrus"
	"log"
)

func main() {

	var opts struct {
		Verbose bool `short:"v" long:"verbose" description:"Always show full config tree"`
		Copy bool `short:"c" long:"copy" description:"Copy config from server1 -> server2"`
	}

	args, err := flags.ParseArgs(&opts, os.Args)

	if err != nil {
		panic(err)
	}

	servers := args[1:]

	firstServer := "localhost"
	if len(servers) == 0 {
		firstServer = "localhost"
	} else if len(servers) == 1 {
		firstServer = servers[0]
	} else {
		copyConfig(servers[0], servers[1])
		return
	}
	printConfig(firstServer, opts.Verbose)
}

func getUrl(server string) string {
	return fmt.Sprintf("http://%s:5000/config/__active/", server)
}

func printConfig(server string, verbose bool) error {
	conf := getConfig(server)

	if(verbose == false) {
		lines := strings.Split(conf, "\n")
		fmt.Println(strings.Join(lines[:9], "\n"))

		if(len(lines) >10) {
			fmt.Printf("[...] Omitted %v lines\n", len(lines) - 10)
		}
	} else {
	    fmt.Println(conf)
	}
	return nil
}

func getConfig(server string) string{
	url := getUrl(server)
	fmt.Println("Get config from:", url)
	resp, err := http.Get(url)

	if err != nil {
        log.Fatal(err)
    }

    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)

    if err != nil {

        log.Fatal(err)
    }
	return string(body)
}

func putConfig(server string, cfg string) {
	url := getUrl(server)
	fmt.Println("Put config to:", url)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer([]byte(cfg)))
    if err != nil {
        panic(err)
    }

    // set the request header Content-Type for json
    req.Header.Set("Content-Type", "application/json")
    // initialize http client
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
	fmt.Println(resp)
}

func copyConfig(server1 string, server2 string) error {
	fmt.Printf("Copy config from: %s to: %s ... \n", server1, server2)
	cfg := getConfig(server1)
	putConfig(server2, cfg)
	return nil
}
