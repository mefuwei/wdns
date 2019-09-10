#!/bin/bash

echo "clean workspace..."
rm -rf wdns.exe
rm -rf wdns
echo "clean workspace success."

echo "start build package...."
case $1 in
    "win")
        CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o wdns.exe main.go
        mv wdns.exe bin/
        ;;
    "linux")
        CCGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o wdns-linux main.go
        mv wdns-linux bin/
    ;;
    "mac")
        CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o wdns-mac main.go
        mv wdns-mac bin/
    ;;
    *)
        echo "build all system to bin/"
        echo "build windows system to bin/wdns.exe"
        CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o wdns.exe main.go
        mv wdns.exe bin/

        echo "build linux system to bin/wdns-linux"
        CCGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o wdns-linux main.go
        mv wdns-linux bin/

        echo "build mac system to bin/wdns-mac"
        CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o wdns-mac main.go
        mv wdns-mac bin/
        exit
    ;;
esac

echo "build package success."


