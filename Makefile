# Copyright (C) 2018 MizukiSonoko. All rights reserved.

test:
	go test -v -race ./...

cover:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...

build:
	go build example/app.go

clean:
	-rm goparse


.PHONY: test cover
