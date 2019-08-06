#!/bin/sh

VERSION=0.5.1

DIR=$(cd $(dirname $0); pwd)
cd $DIR

rm -rf bin/*
go build -v -o bin/ultproxy -a -tags netgo -installsuffix netgo cmd/ultproxy/main.go

tar cvzf ultproxy-$VERSION.tar.gz bin README.md LICENSE

