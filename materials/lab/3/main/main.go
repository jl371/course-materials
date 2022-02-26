// Build and Use this File to interact with the shodan package
// In this directory lab/3/shodan/main:
// go build main.go
// SHODAN_API_KEY=YOURAPIKEYHERE ./main <search term>

package main

import (
	"fmt"
	"log"
	"os"
	"encoding/json"
	"shodan/shodan"
	"strconv"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalln("Usage: main <searchterm> <pageamount>")
	}
	apiKey := os.Getenv("SHODAN_API_KEY")
	s := shodan.New(apiKey)
	info, err := s.APIInfo()
	if err != nil {
		log.Panicln(err)
	}
	fmt.Printf(
		"Query Credits: %d\nScan Credits:  %d\n\n",
		info.QueryCredits,
		info.ScanCredits)

	pa, err2 := strconv.Atoi(os.Args[2]) 

	if err2 != nil {
		log.Panicln(err)
	}

	for i:=0; i < pa; i++ {
		hostSearch, err := s.HostSearch(os.Args[1],i)
	if err != nil {
		log.Panicln(err)
	}
	
	fmt.Printf("Host Data Dump, page %d\n", i)
	for _, host := range hostSearch.Matches {
		fmt.Println("==== start ",host.IPString,"====")
		h,_ := json.Marshal(host)
		fmt.Println(string(h))
		fmt.Println("==== end ",host.IPString,"====")
	}


	fmt.Printf("IP, Port\n")

	for _, host := range hostSearch.Matches {
		fmt.Printf("%s, %d\n", host.IPString, host.Port)
	}
	}
	


}