#!/bin/bash

fail=0
while true; do
    echo building
    go build -o game
    if test $? -ne 0 ; then

        fail=$((fail + 1))
        if test $fail -gt 10 ; then 
           exit 
        fi

        sleep 1
        continue
    fi
    fail=0

    echo running
    EXIT_ON_MODIFY=1 ./game
    if test $? -ne 0 ; then
        exit
    fi
done