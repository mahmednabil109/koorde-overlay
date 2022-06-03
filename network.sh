#!/bin/bash

# run the bootstrap node
go run main.go -port 8080 &
sleep 1


# create the other nodes
port=8080
dport=16585
for ((i=1; i<$1; i++)); do
     ((port+=1))
     ((dport+=1))
     ./bin/koorde-overlay -port $port -first 0 -debug $2 -dport $dport -bootstrap 127.0.0.1:8080 &
     sleep 1
done
