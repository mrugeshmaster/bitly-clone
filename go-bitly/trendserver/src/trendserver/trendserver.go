package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"os"
)

// GCP  
var trend_server = os.Getenv("trend_server")

// var trend_server = "localhost:9090"

type links_details struct {
	Shorturl_code string `json:"shorturl_code"`
	Accessed_at   string `json:"accessed_at"`
}

type reqBody struct {
	Listoflinks []links_details `json:"listoflinks"`
}

var post_links_array []links_details
var put_links_array []links_details

func addToListOfLinks(shorturl_code string) {

	links_det := links_details{
		Shorturl_code: shorturl_code,
		Accessed_at:   time.Now().Format(" Mon Jan _2 2006 15:04:05 "),
	}
	fmt.Println(links_det)

	resp, err := http.Get("http://" + trend_server + "/api/listoflinks")
	if err != nil || resp.StatusCode != 200 {
		fmt.Println(links_det)
		post_links_array = append(post_links_array, links_det)

		reqBody, _ := json.Marshal(map[string][]links_details{
			"listoflinks": post_links_array,
		})

		resp, err := http.Post("http://"+trend_server+"/api/listoflinks", "application/json", bytes.NewBuffer(reqBody))
		if err != nil || resp.StatusCode != 200 {
			fmt.Println("Failed to create document")
		} else {
			fmt.Println("Successfully created document")
		}
		log.Fatal(err)
	}

	respBody := map[string][]links_details{}

	if err != nil || resp.StatusCode != 200 {
		fmt.Println("Failed to fetch document")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &respBody)

	fmt.Println("LongURL from Cache : ", respBody["listoflinks"])

	put_links_array = respBody["listoflinks"]

	links_det = links_details{
		Shorturl_code: shorturl_code,
		Accessed_at:   time.Now().Format(" Mon Jan _2 2006 15:04:05 "),
	}

	put_links_array = append(put_links_array, links_det)

	reqB := reqBody{Listoflinks: put_links_array}

	byteArray, err := json.Marshal(reqB)

	req, _ := http.NewRequest("PUT", "http://"+trend_server+"/api/listoflinks", bytes.NewBuffer(byteArray))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	if _, err := client.Do(req); err != nil {
		log.Fatal(err)
	}
}