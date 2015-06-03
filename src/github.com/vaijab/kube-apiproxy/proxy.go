package main

import (
	"io"
	"log"
	"net"
)

func proxy(lconn net.Conn, remoteAddr string) {
	r, err := net.ResolveTCPAddr("tcp", remoteAddr)
	rconn, err := net.DialTCP("tcp", nil, r)
	if err != nil {
		log.Println(err)
		return
	}
	if rconn == nil {
		log.Println("failed to connect to remote socket.")
	}
	go io.Copy(lconn, rconn)
	go io.Copy(rconn, lconn)
}
