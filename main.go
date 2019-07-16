/*
share-dir is a very simple static file server in go
Usage:
    -d=".":    the directory of static files to host
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func validInternalIP(ipAddress string) bool {
	ipAddress = strings.Trim(ipAddress, " ")

	re, _ := regexp.Compile(`^192.(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){2}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	if re.MatchString(ipAddress) {
		return true
	}
	return false
}

func main() {

	rand.Seed(time.Now().UnixNano())

	port := strconv.Itoa(1000 + rand.Intn(8999))
	directory := flag.String("d", ".", "the directory of static file to host")
	flag.Parse()

	interfaces, err := net.Interfaces()

	if err != nil {
		log.Print(err)
		return
	}

	for _, i := range interfaces {
		byNameInterface, err := net.InterfaceByName(i.Name)

		if err != nil {
			log.Println(err)
		}

		addresses, err := byNameInterface.Addrs()
		for _, v := range addresses {
			ip, _, _ := net.ParseCIDR(v.String())
			//fmt.Printf("Sharing %s on http://%v:%s\n", *directory, ip[:strings.IndexByte(ip, '/')], port)
			if validInternalIP(ip.String()) {
				fmt.Printf("Sharing %s on http://%v:%s\n", *directory, ip, port)
			}
		}
		//	}
	}

	http.Handle("/", http.FileServer(http.Dir(*directory)))
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
