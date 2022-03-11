package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"bufio"
	"github.com/hashicorp/yamux"
	"strconv"
	"strings"
	"time"
)

var proxytout = time.Millisecond * 1000 //timeout for wait magicbytes
var Sessions []*yamux.Session
var socksIp string

// listen for agents
func listenForAgents(address string, clients string, agentPassword string) error {
	var err, erry error
	var cer tls.Certificate
	var session *yamux.Session
	var ln net.Listener

	log.Printf("Will start listening for clients on %s", clients)
	log.Printf("Listening for agents on %s using TLS", address)
	cer, err = getRandomTLS(2048)
	if err != nil {
		log.Println(err)
		return err
	}
	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	ln, err = tls.Listen("tcp", address, config)
	if err != nil {
		log.Printf("Error listening on %s: %v", address, err)
		return err
	}
	var listenstr = strings.Split(clients, ":")
	portnum, errc := strconv.Atoi(listenstr[1])
	if errc != nil {
		log.Printf("Error converting listen str %s: %v", clients, errc)
	}
	socksIp = listenstr[0]

	for {
		conn, err := ln.Accept()
		conn.RemoteAddr()
		agentstr := conn.RemoteAddr().String()
		log.Printf("[%s] Got a connection from %v: ", agentstr, conn.RemoteAddr())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Errors accepting!")
		}

		reader := bufio.NewReader(conn)

		//read only 64 bytes with timeout=1-3 sec. So we haven't delay with browsers
		conn.SetReadDeadline(time.Now().Add(proxytout))
		statusb := make([]byte, 64)
		_, _ = io.ReadFull(reader, statusb)

		//Alternatively  - read all bytes with timeout=1-3 sec. So we have delay with browsers, but get all GET request
		//conn.SetReadDeadline(time.Now().Add(proxytout))
		//statusb,_ := ioutil.ReadAll(magicBuf)

		//log.Printf("magic bytes: %v",statusb[:6])
		//if hex.EncodeToString(statusb) != magicbytes {
		if string(statusb)[:len(agentPassword)] != agentPassword {
			//do HTTP checks
			log.Printf("Received request: %v", string(statusb[:64]))
			status := string(statusb)
			if strings.Contains(status, " HTTP/1.1") {
				httpresonse := "HTTP/1.1 301 Moved Permanently" +
					"\r\nContent-Type: text/html; charset=UTF-8" +
					"\r\nLocation: https://www.google.com/search?q=%D1%80%D1%83%D1%81%D1%81%D0%BA%D0%B8%D0%B9+%D0%B2%D0%BE%D0%B5%D0%BD%D0%BD%D1%8B%D0%B9+%D0%BA%D0%BE%D1%80%D0%B0%D0%B1%D0%BB%D1%8C+%D0%B8%D0%B4%D0%B8+%D0%BD%D0%B0%D1%85%D1%83%D0%B9" +
					"\r\nServer: Apache" +
					"\r\nContent-Length: 0" +
					"\r\nConnection: close" +
					"\r\n\r\n"

				conn.Write([]byte(httpresonse))
				conn.Close()
			} else {
				conn.Close()
			}

		} else {
			//magic bytes received.
			//disable socket read timeouts
			log.Printf("[%s] Got Client from %s", agentstr, conn.RemoteAddr())
			conn.SetReadDeadline(time.Now().Add(100 * time.Hour))
			//Add connection to yamux
			session, erry = yamux.Client(conn, nil)
			if erry != nil {
				log.Printf("[%s] Error creating client in yamux for %s: %v", agentstr, conn.RemoteAddr(), erry)
				continue
			}
			Sessions = append(Sessions, session)
			connection := addConnection(portnum, session, conn)
			go listenForClients(connection)
		}
	}
	return nil
}

// Catches local clients and connects to yamux
func listenForClients(connection *Connection) error {
	var ln net.Listener
	var address string
	var err error
	agentstr := connection.Conn.RemoteAddr().String()
	portinc := connection.Port
	for {
		address = fmt.Sprintf("%s:%d", socksIp, portinc)
		log.Printf("[%s] Waiting for clients on %s", agentstr, address)
		ln, err = net.Listen("tcp", address)
		if err != nil {
			log.Printf("[%s] Error listening on %s: %v", agentstr, address, err)
			portinc = portinc + 1
		} else {
			connection.Port = portinc
			break
		}
	}

	connection.Ln = ln

	for {
		log.Printf("wait for accept")
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("[%s] Error accepting on %s: %v", agentstr, address, err)
			return err
		}
		if connection.Session == nil {
			log.Printf("[%s] Session on %s is nil", agentstr, address)
			conn.Close()
			continue
		}
		log.Printf("[%s] Got client. Opening stream for %s", agentstr, conn.RemoteAddr())

		stream, err := connection.Session.Open()
		if err != nil {
			log.Printf("[%s] Error opening stream for %s: %v", agentstr, conn.RemoteAddr(), err)
			return err
		}

		// connect both of conn and stream

		go func() {
			log.Printf("[%s] Starting to copy conn to stream for %s", agentstr, conn.RemoteAddr())
			io.Copy(conn, stream)
			conn.Close()
			log.Printf("[%s] Done copying conn to stream for %s", agentstr, conn.RemoteAddr())
		}()
		go func() {
			log.Printf("[%s] Starting to copy stream to conn for %s", agentstr, conn.RemoteAddr())
			io.Copy(stream, conn)
			stream.Close()
			log.Printf("[%s] Done copying stream to conn for %s", agentstr, conn.RemoteAddr())
		}()
	}
}
