package main

import (
	"flag"

	"github.com/TedMANNERs/rps-go-master/master"
)

func main() {
	azureUrlPtr := flag.String("project", "https://techweek2018functionapp.azurewebsites.net/api", "Azure API url")
	flag.Parse()
	master.Start(*azureUrlPtr)
}
