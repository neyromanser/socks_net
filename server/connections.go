package main

import (
	"github.com/hashicorp/yamux"
	"log"
	"net"
	"strings"
	"time"
)

type Connection struct {
	Country 		string
	Ip 				string
	Port 			int
	SessionIndex 	int
	Session 		*yamux.Session
	Conn 			net.Conn
	IsOnline 		bool
	Ln 				net.Listener
	//Cmd     		chan string
}

var Connections []*Connection

func (c *Connection) Close() {
	log.Printf("run Stop method")
	c.IsOnline = false
	c.Ln.Close()
	//c.Cmd <- "stop"
}

func checkStatus(){
	for {
		for _, c := range Connections {
			if c.Session.IsClosed() && c.IsOnline {
				c.Close()
			}else{
				log.Printf("session active %s:%d (%s) - %t", c.Ip, c.Port, c.Country, c.IsOnline)
			}
		}
		time.Sleep(time.Second * 5)
	}
}

func addConnection(portNum int, session *yamux.Session, conn net.Conn) (connection *Connection) {
	ipSting := strings.Split(conn.RemoteAddr().String(), ":")
	i := 0
	for _, connection := range Connections {
		if connection.Session.IsClosed() {
			connection.IsOnline = false
			break
		}
		i++
	}
	sockPort := portNum + i
	connection = &Connection{
		Ip: ipSting[0],
		Port: sockPort,
		SessionIndex: len(Sessions) - 1,
		Session: session,
		IsOnline: true,
		Country: getCountry(ipSting[0]),
		Conn: conn,
		//Cmd: make(chan string),
	}
	if i >= len(Connections) {
		Connections = append(Connections, connection)
	}else{
		Connections[i] = connection
	}
	return
}