#!/bin/bash

echo "clean workspace..."
rm -rf wdns.exe
rm -rf wdns
echo "clean workspace success."

echo "start build package...."
case $1 in
    "win")
        CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go
        mv wdns.exe bin/
        ;;
    "linux")
        CCGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
        mv wdns-linux bin/
    ;;
    "mac")
        CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build main.go
        mv wdns-mac bin/
    ;;
    *)
        echo "Usage : build.sh win|linux|mac"
        exit
    ;;
esac

echo "build package success."


