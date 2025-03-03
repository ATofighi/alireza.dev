package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// https://httpbin.org/status/400
	req, err := http.NewRequest("GET", "https://httpbin.org/get", nil)
	req.Header.Set("Test-Header", "testtsett")
	if err != nil {
		log.Fatal(err)
	}
	req.Write(os.Stdout)
	fmt.Println("==========")
	fmt.Println("Host:", req.Host)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n\n===== RESPONSE")
	fmt.Println("status code = ", resp.StatusCode)
	for k, vs := range resp.Header {
		fmt.Printf("%s: %d, %+v\n", k, len(vs), vs)
	}
	fmt.Println("=========================")
	sc := bufio.NewScanner(resp.Body)
	for sc.Scan() {
		fmt.Println(sc.Text())
	}
}
