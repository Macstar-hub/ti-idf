#!/bin/bash 
run_house_crawler () {
    while true
    do  
        go run house_crawler/house_crawler.go
        if [ $? -eq 1 ]; then
            echo "Stopping script as condition met."
            exit 0
        fi
        sleep 120
    done
}


run_house_crawler



