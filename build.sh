#!/bin/bash



case $1 in
    "win")
        CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build
        tar czf wdns.tar.gz wdns.exe etc/
        ;;
    "linux")
        CCGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
        tar czf wdns.tar.gz wdns etc/
    ;;
    "mac")
        CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build
        tar czf wdns.tar.gz wdns etc/

    ;;
    *)
        echo "Usage : build.sh win|linux|mac"
        exit
    ;;
esac

echo "make succeed, package wdns.tar.gz "


