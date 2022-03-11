# SocksNet
Reverse SOCKS5 bot net

## dependencies
install https://github.com/mitchellh/gox 

## build
```
chmod +x genrelease.sh  

# to build all platform binaries  
./genrelease.sh

# to build only current relese    
./genrelease.sh dev 
```
client and server binary files will be stored in release folder 

## config options
copy .env.example to .env  
```
AGENT_PASSWORD=Secret{p}aSSword 
- password for authenticating new agents (clients who creates reverse socks5 tunnels) 

AGENT_PORT=:8443 
- server port on listen for agents connections

SOCKS_ADDRESS=127.0.0.1:1080 
- server local address for socks5 connection. each new agent will open new port (incrementing from selected) or will use free one from previusly closed connecion

CONTROL_DOMAIN=https://ruskiykorablidinahuy.today  
- when clien run without server ip flag (or when it lose connection) it will connect to control domain to get ip:port address of server 
(andpoint: {CONTROL_DOMAIN}/spoint.html).  

RPC_PORT=8484
RPC_USER=user
RPC_PASSWORD=passw0rd
- server RPC creds, comming soon...
```


## usage  
1. start server  
```chmod +x server && ./server ```
2. start client(s)    
```client.exe ```
3. use available socks5 proxy on server machine  
``` 
curl --socks5 localhost:1080 https://ifconfig.me/ 
```
