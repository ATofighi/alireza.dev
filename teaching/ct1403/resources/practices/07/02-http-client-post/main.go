package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func main() {
	resp, err := http.PostForm("https://httpbin.org/post", url.Values{
		"input1": []string{"value1", "value2"},
		"input2": []string{"value3"},
	})
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
