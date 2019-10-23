// Copyright (C) 2018,2019 MizukiSonoko. All rights reserved.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	goparse "github.com/MizukiSonoko/goparse/parse"
)

func main() {
	resp, err := http.Get("http://example.com")
	if err != nil {
		log.Fatalf("Get failed err:%s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ReadAll failed err:%s", err)
	}

	var title, charset string
	err = goparse.Parse("%s<title>%s</title>%s", string(body)).InsertOnly(1, &title)
	if err != nil {
		log.Fatalf("Parse failed err:%s", err)
	}
	fmt.Printf("title is %s\n", title)

	err = goparse.Parse("%s<meta charset=\"%s\" />%s", string(body)).InsertOnly(1, &charset)
	if err != nil {
		log.Fatalf("Parse failed err:%s", err)
	}
	fmt.Printf("charset is %s\n", charset)
}
