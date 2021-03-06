# Copyright (C) 2018 MizukiSonoko. All rights reserved.

test:
	go test -v -race ./...

cover:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...

example:
	cd example; go build  -o crowler crowler.go; ./crowler

build:
	go build example/app.go

clean:
	-rm goparse
	-rm goparse
	-rm example/crowler


.PHONY: test cover example clean