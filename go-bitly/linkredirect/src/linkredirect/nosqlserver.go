package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"os"
)

// GCP
var cache_server = os.Getenv("cache_server")

// AWS
// var cache_server = "bitly-e66492f0c65474e2.elb.us-west-2.amazonaws.com:9090"

// Local
// var cache_server = "localhost:9001"

func cleanCache() {

	cacheEntries := []map[string]interface{}{}

	client := &http.Client{}

	forever := make(chan bool)

	time.Sleep(5 * time.Minute)

	resp, err := http.Get("http://" + cache_server + "/api")
	if err != nil || resp.StatusCode != 200 {
		fmt.Println("Fetching Cache Entries failed")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &cacheEntries)
	for i := range cacheEntries {
		if cacheEntries[i]["key"] != "listoflinks" && cacheEntries[i]["message"] != "deleteFlag" {
			// fmt.Println("cacheEntries : ", cacheEntries[i])
			delreq, err := http.NewRequest(http.MethodDelete, "http://"+cache_server+"/api/"+cacheEntries[i]["key"].(string)+"", nil)

			delResp, err := client.Do(delreq)

			if err != nil || delResp.StatusCode != 200 {
				fmt.Println("Delete cache failed")
			}
			fmt.Println("cacheEntries Deleted: ", cacheEntries[i])
		}
	}
	<-forever
}

func writeInCache(shortURL string, longURL string) int {

	reqBody, _ := json.Marshal(map[string]string{
		"longURL": longURL,
	})

	resp, err := http.Post("http://"+cache_server+"/api/"+shortURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil || resp.StatusCode != 200 {
		fmt.Println("Failed to create document")
	} else {
		fmt.Println("Successfully created document")
	}
	return resp.StatusCode

}

func getLongURLFromCache(shortlink_code string) string {

	respBody := map[string]string{}

	resp, err := http.Get("http://" + cache_server + "/api/" + shortlink_code)
	if err != nil || resp.StatusCode != 200 {
		fmt.Println("Failed to fetch document")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &respBody)
	fmt.Println("LongURL from Cache : ", respBody["longURL"])
	return respBody["longURL"]
}
