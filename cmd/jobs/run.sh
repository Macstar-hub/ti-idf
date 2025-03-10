#!/bin/bash 
run_house_crawler () {
    while true
    do  
        go run web_crawler/house_crawler/house_crawler.go
        if [ $? -eq 1 ]; then
            echo "Stopping script as condition met."
            break
            # exit 0
        fi
        sleep 120
    done
}

run_price_analyse () {
    cd /Users/Shared/scripts.dir/ && bash /Users/Shared/scripts.dir/reset-tunnel.sh
    export  http_proxy='http://localhost:9443'
    export  https_proxy='http://localhost:9443'
    cd /Users/Shared/codes.dir/go.dir/git.dir/ti-idf/cmd/jobs
    go run price_analyzes/price_analyzes.go
    if [ $? -eq 1 ]; then
        echo "Price analyzing has problem."
        unset http_proxy https_proxy
        exit 1
    fi
    unset http_proxy https_proxy
}

run_price_hist () {
    python3 price_analyzes/priceBarchart.py
    if [ $? -eq 1 ]; then
        echo "Price histogram plotting has problem."
        exit 1
    fi
}

run_house_crawler
run_price_analyse
run_price_hist



