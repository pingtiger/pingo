package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/docopt/docopt-go"
	"github.com/robinjmurphy/pingo"
)

const Version = "1.0.0"
const Usage = `Pingo.

Usage:
  pingo <hostname> <port> [-i <interval>] [-t <timeout>] [-r <region>] [-ps]
  pingo -h | --help
  pingo -v | --version

Options:
  -h --help                 Show this screen.
  -v --version              Show version.
  -i --interval <interval>  Ping interval in seconds. [default: 60]
  -t --timeout <timeout>    Time in milliseconds to wait for a TCP response. [default: 1000]
  -r --region <region>      AWS region.
  -p --publish              Publish results to CloudWatch (requires AWS environment variables to be set.)
  -s --silent               Disable all output.`

func main() {
	arguments, _ := docopt.Parse(Usage, nil, true, Version, false)
	hostname := arguments["<hostname>"].(string)
	port, _ := strconv.Atoi(arguments["<port>"].(string))
	interval, _ := strconv.Atoi(arguments["--interval"].(string))
	timeout, _ := strconv.Atoi(arguments["--timeout"].(string))
	region := arguments["--region"]
	publish := arguments["--publish"].(bool)
	silent := arguments["--silent"].(bool)
	host := pingo.Host{hostname, port}
	handlers := []pingo.Handler{}
	if publish {
		if region == nil {
			fmt.Println("No AWS region specified")
			os.Exit(1)
		}
		handlers = append(handlers, pingo.NewCloudWatchHandler(region.(string)))
	}
	if !silent {
		handlers = append(handlers, pingo.LoggingHandler)
	}
	for {
		pingo.Ping(host, time.Duration(timeout)*time.Millisecond, handlers)
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
