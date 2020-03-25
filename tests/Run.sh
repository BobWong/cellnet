#!/usr/bin/env bash
export CURRDIR=`pwd`
cd ../../../../..
export GOPATH=`pwd`
cd ${CURRDIR}

set -e

go test -v .

# 编译例子
trap 'rm -rf examplebin' EXIT
mkdir -p examplebin
go build -p 4 -v -o ./examplebin/echo github.com/bobwong89757/cellnet/examples/echo
go build -p 4 -v -o ./examplebin/echo github.com/bobwong89757/cellnet/examples/chat/client
go build -p 4 -v -o ./examplebin/echo github.com/bobwong89757/cellnet/examples/chat/server
go build -p 4 -v -o ./examplebin/echo github.com/bobwong89757/cellnet/examples/fileserver
go build -p 4 -v -o ./examplebin/echo github.com/bobwong89757/cellnet/examples/websocket


