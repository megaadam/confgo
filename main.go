package main

import (
	//"flags" // use go-flags!!
	"fmt"
	"errors"
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

func printConfig(server string, verbose bool) error {
	confUrl := fmt.Sprintf("http://%s:5000/config/__active/", server)

	fmt.Println("Config from: ", confUrl)
	resp, err := http.Get(confUrl)

	if err != nil {
        log.Fatal(err)
    }

    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)

    if err != nil {

        log.Fatal(err)
    }

	if(verbose == false) {
		lines := strings.Split(string(body), "\n")
		fmt.Println(strings.Join(lines[:9], "\n"))

		if(len(body) >10) {
			fmt.Printf("[...] Omitted %v lines\n", len(lines) - 10)
		}
	} else {
	    fmt.Println(string(body))
	}
	return nil
}

func copyConfig(server1 string, server2 string) error {
	fmt.Printf("Copy config from: %s to: %s ... \n", server1, server2)
	fmt.Println("You wish!!")
	return errors.New("Not implemented.")
}
