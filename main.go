package main

import (
	"flag"

	"github.com/TedMANNERs/rps-go-master/master"
)

func main() {
	azureURLPtr := flag.String("project", "https://techweek2018functionapp.azurewebsites.net/api", "Azure API url")
	portPtr := flag.Int("port", 8080, "The port to use for hosting")
	networkMatchAddrPtr := flag.String("matchAddress", "192.168.201", "The network address used for finding the correct network adapter (e.g. 192.168.201 for 192.168.201.x)")
	flag.Parse()
	master.Start(*azureURLPtr, *portPtr, *networkMatchAddrPtr)
}
