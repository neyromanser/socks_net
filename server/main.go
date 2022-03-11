package main

import (
	"log"
	"github.com/neyromanser/socks_net/helpers"
)

func main(){
	config := helpers.GetConfig(".")
	agentPassword := config.AgentPassword
	agentPort := config.AgentPort
	socksAddress := config.SocksAddress

	go checkStatus()
	log.Fatal(listenForAgents(agentPort, socksAddress, agentPassword))
}

