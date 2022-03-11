package main

import (
	"log"
)

func main(){
	config := GetConfig("..")
	agentPassword := config.AgentPassword
	agentPort := config.AgentPort
	socksAddress := config.SocksAddress

	go checkStatus()
	log.Fatal(listenForAgents(agentPort, socksAddress, agentPassword))
}

