#!/usr/bin/env bash

for (( i=10; i<=400; i+=5 )); do
     go run ./node_lookup.go -node-num $i -lookups 16 2> "log_out_$i.log"
     sleep .5
     ../clean.sh
     ./summerize.py 
     echo "Done $i"
     sleep 4
done