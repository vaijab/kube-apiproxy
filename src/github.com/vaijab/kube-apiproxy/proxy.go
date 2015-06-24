package main

import (
	"io"
	"log"
	"net"
	"sync"
)

func proxy(lconn *net.TCPConn, remoteAddr string) {
	var wg sync.WaitGroup
	wg.Add(2)
	r, err := net.ResolveTCPAddr("tcp", remoteAddr)
	rconn, err := net.DialTCP("tcp", nil, r)
	if err != nil {
		log.Println(err)
		return
	}
	if rconn == nil {
		log.Println("failed to connect to remote socket.")
	}
	go copyBytes(lconn, rconn, &wg)
	go copyBytes(rconn, lconn, &wg)
	wg.Wait()
	lconn.Close()
	rconn.Close()
}

func copyBytes(dest, src *net.TCPConn, wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := io.Copy(dest, src)
	if err != nil {
		log.Printf("I/O error: %v\n", err)
	}
	dest.CloseWrite()
	src.CloseRead()
}
