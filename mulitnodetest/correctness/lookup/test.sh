#!/usr/bin/env bash

for (( i=3; i<=8; i+=1 )); do
     go run ./main.go -node-num $((1<<$i)) -lookups 200 2> "log_out_$i.log"
     sleep .5
     ../../clean.sh
     ../../summerize.py 
     echo "Done $((1<<$i))"
     sleep 2
done