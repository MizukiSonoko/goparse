# Copyright (C) 2018 MizukiSonoko. All rights reserved.

test:
	vgo test -v -race ./...

cover:
	vgo test -race -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: test cover
