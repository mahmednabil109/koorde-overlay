#!/usr/bin/env bash
pgrep koord | xargs kill -9

# for p in $(pgrep koorde); do
#      echo $(ps  $p)
# done