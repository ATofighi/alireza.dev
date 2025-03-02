package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Test struct {
	Enabled bool
	Codes   []string "json:\"code_ha\""
}

func main() {
	// data, err := json.Marshal(Test{
	// 	Enabled: false,
	// 	Codes: []string{
	// 		"asdasd",
	// 		"vsdfsdf",
	// 	},
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }
	var data Test
	fmt.Println(data)
	err := json.Unmarshal([]byte(`{"Enabled":false,"code_ha":["asdasd","vsdfsdf"]}`), &data)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(data)
}
