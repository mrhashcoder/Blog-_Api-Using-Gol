package main

//for sending http post request to server for creating new blog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func makeReq() {

	body := map[string]interface{}{
		"Id":          "4",
		"Author":      "hashCoder",
		"Title":       "Gawds-Induction-Project",
		"Description": "Gawds Is Love",
		"content":     "Anuj sir please induct krlo pakka saare kaam karunga",
	}

	bytesRepresentation, err := json.Marshal(body)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post("http://localhost:1234/createBlog", "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(resp)
	fmt.Println("reach till here")

}

func main() {
	makeReq()
}
