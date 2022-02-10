// bhg-scanner/scanner.go modified from Black Hat Go > CH2 > tcp-scanner-final > main.go
// Code : https://github.com/blackhat-go/bhg/blob/c27347f6f9019c8911547d6fc912aa1171e6c362/ch-2/tcp-scanner-final/main.go
// License: {$RepoRoot}/materials/BHG-LICENSE
// Useage:
// PortScanner requires a string parameter. It will scan ports 1 - 1024 of an address you give it. Use the syntax:
// PortScanner("address.org")
// or if you do not wish to enter an address (it will use a default value)
// PortScanner("")

package scanner

import (
	"fmt"
	"net"
	"sort"
	"time"
)

//TODO 3 : ADD closed ports; currently code only tracks open ports
var openports []int  // notice the capitalization here. access limited!
var closedports []int


func worker(ports, results chan int, address string) {
	for p := range ports {
		address := fmt.Sprintf(address + ":%d", p)    
		conn, err := net.DialTimeout("tcp", address, 1 * time.Second) // TODO 2 : REPLACE THIS WITH DialTimeout (before testing!)
		if err != nil { 
			results <- -1 * p
			continue
		}
		conn.Close()
		results <- p
	}
}

// for Part 5 - consider
// easy: taking in a variable for the ports to scan (int? slice? ); a target address (string?)?
// med: easy + return  complex data structure(s?) (maps or slices) containing the ports.
// hard: restructuring code - consider modification to class/object 
// No matter what you do, modify scanner_test.go to align; note the single test currently fails
func PortScanner(addr string) (int, int) {  

	ports := make(chan int, 100)   // TODO 4: TUNE THIS FOR CODEANYWHERE / LOCAL MACHINE
	results := make(chan int)
	if (addr == ""){
		for i := 0; i < cap(ports); i++ {
			go worker(ports, results,"scanme.nmap.org")
		}
	} else {
		for i := 0; i < cap(ports); i++ {
			go worker(ports, results,addr)
		}
	}
	
	

	go func() {
		for i := 1; i <= 1024; i++ {
			ports <- i
		}
	}()

	for i := 0; i < 1024; i++ {
		port := <-results
		if port > 0 {
			openports = append(openports, port)
		} else if port < 0 {
			closedports = append(closedports, -1 * port)
		}
	}

	close(ports)
	close(results)
	sort.Ints(openports)

	//TODO 5 : Enhance the output for easier consumption, include closed ports

	for _, port := range openports {
		fmt.Printf("%d: open\n", port)
	}
	for _, port := range closedports {
		fmt.Printf("%d: closed\n", port)
	}
	return len(openports), len(closedports)// TODO 6 : Return total number of ports scanned (number open, number closed); 
	//you'll have to modify the function parameter list in the defintion and the values in the scanner_test
}
