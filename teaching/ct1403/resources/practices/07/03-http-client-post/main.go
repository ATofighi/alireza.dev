package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type MyRequest struct {
	X string
	Y int
}

func main() {
	data, _ := json.Marshal(MyRequest{
		X: "salam",
		Y: 10,
	})
	resp, err := http.Post("https://httpbin.org/post", "application/json", bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\n\n===== RESPONSE")
	for k, vs := range resp.Header {
		fmt.Printf("%s: %d, %+v\n", k, len(vs), vs)
	}
	fmt.Println("=========================")
	sc := bufio.NewScanner(resp.Body)
	for sc.Scan() {
		fmt.Println(sc.Text())
	}
}
