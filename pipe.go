// netln allows to redirect incoming network connections to another
// address, like a transparent proxy.
//
// Usage: netln [-s maxspeed] remote-addr local-addr
//
// maxspeed is in bits per second. Forms like "33.6K" or "1M" are
// recognized.
package main

import (
	"errors"
	"flag"
	"github.com/gaswelder/iorate"
	"io"
	"log"
	"net"
)

func usage() {
	log.Fatal("Usage: netln [-s <max speed>] <remote addr> <local addr>")
}

func main() {
	// Parse the command line
	var maxSpeedString string
	var listenAddr, connectAddr string
	flag.StringVar(&maxSpeedString, "s", "none", "maximum transfer speed")
	flag.Parse()
	connectAddr = flag.Arg(0)
	listenAddr = flag.Arg(1)
	if listenAddr == "" || connectAddr == "" {
		usage()
	}
	maxSpeed, err := parseSpeed(maxSpeedString)
	if err != nil {
		log.Fatal(err)
	}

	// Bind to the local address and process clients.
	in, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s -> %s\n", listenAddr, connectAddr)
	for {
		client, err := in.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go processClient(client, connectAddr, maxSpeed)
	}
}

func processClient(client net.Conn, connectAddr string, maxSpeed int64) {
	log.Println(client.RemoteAddr(), "connected")

	// Connect to the target
	server, err := net.Dial("tcp", connectAddr)
	if err != nil {
		log.Println(err)
		client.Close()
		return
	}

	// If maxSpeed is given, add throttles to the writers.
	var cwriter, swriter io.Writer
	if maxSpeed != -1 {
		cwriter = iorate.NewWriter(client, iorate.Rate(maxSpeed)*iorate.Bps)
		swriter = iorate.NewWriter(server, iorate.Rate(maxSpeed)*iorate.Bps)
	} else {
		cwriter = client
		swriter = server
	}

	// Pass the data in both directions
	schan := getReadChan(server)
	cchan := getReadChan(client)
	var data []byte
	var n int
	sent := 0
	received := 0
	for err == nil {
		ok := false
		select {
		case data, ok = <-schan:
			n, err = write(data, cwriter)
			sent += n

		case data, ok = <-cchan:
			n, err = write(data, swriter)
			received += n
		}
		if !ok {
			err = errors.New("Channel closed")
		}
	}

	if err != nil {
		log.Println(err)
	}

	client.Close()
	server.Close()

	log.Printf("%s disconnected; %d sent, %d received.\n", client.RemoteAddr(),
		sent, received)
}

func getReadChan(r net.Conn) chan []byte {
	buf := make([]byte, 4096)
	c := make(chan []byte, 0)
	go (func() {
		for {
			n, err := r.Read(buf)
			if err != nil {
				close(c)
				break
			}
			c <- buf[:n]
		}
	})()
	return c
}

func write(data []byte, w io.Writer) (int, error) {
	n := len(data)
	sent := 0
	s := 0
	var err error

	for sent < n && err == nil {
		s, err = w.Write(data[sent:])
		sent += s
	}
	if sent < n && err == nil {
		err = errors.New("Data truncated")
	}
	return sent, err
}
