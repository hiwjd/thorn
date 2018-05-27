# thorn

thorn is a tool to export local server port to public.

## get started

start `server` at a public server. `./server`

start `client` at the server in local network. `./client`

assume you need to ssh into local server which open sshd on port 22.

```
curl -XPOST http://public-server-ip:9991/openport --data "port=22&vport=2222&clientID=clientID001"
```

then

```
ssh -p 2222 user@public-server-ip
```

## how it works

we tell server what we want like 'I want to connect your port 2222 then you proxy it to local server(which presents as clientID) 22'.

then server tell client to connect port 22 on itself, an pipe the connection to the server port 2222.
