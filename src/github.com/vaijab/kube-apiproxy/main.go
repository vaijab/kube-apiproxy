package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

var (
	fleetEndpoint string
	unitName      string
	proxyListen   string
	apiPort       string
)

func init() {
	flag.StringVar(&fleetEndpoint, "fleet-endpoint", "unix:///run/fleet.sock", "fleet endpoint")
	flag.StringVar(&unitName, "unit-name", "kube-apiserver.service", "fleet unit name for kubernetes api server")
	flag.StringVar(&proxyListen, "proxy-listen", "localhost:8081", "proxy listen ip:port")
	flag.StringVar(&apiPort, "api-port", "8080", "kubernetes api port")
}

func main() {
	flag.Parse()
	ipChan := make(chan string)
	fleetClient, err := getClient(fleetEndpoint)
	if err != nil {
		log.Fatalln(err)
	}
	go func() {
		for {
			ip, err := getUnitMachineIP(fleetClient, unitName)
			if err != nil {
				log.Println(err)
			} else {
				ipChan <- ip
			}
			time.Sleep(time.Second * 3)
		}
	}()

	l, err := net.ResolveTCPAddr("tcp", proxyListen)
	listen, err := net.ListenTCP("tcp", l)
	if err != nil {
		log.Fatalln(err)
	}
	if listen == nil {
		fmt.Errorf("cannot bind to local socket.")
	}

	unitIP := ""
	go func() {
		prevIP := ""
		for {
			select {
			case ip := <-ipChan:
				if prevIP == ip {
					continue
				} else {
					prevIP = ip
					unitIP = ip
				}
			}
		}
	}()
	// TODO(vaijab) Would be great to handle OS signals
	for {
		// TODO(vaijab) it would probably better to use some locking here?
		if unitIP == "" {
			time.Sleep(time.Second * 1)
			continue
		}
		conn, err := listen.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		go proxy(conn, unitIP+":"+apiPort)
	}
}
