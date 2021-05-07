package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Route struct {
	Env     string `json:"env"`
	Country string `json:"country"`
	Url     string `json:"url"`
}

var client http.Client

func get(url string) map[string]interface{} {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "infomiho-checker")

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var value map[string]interface{}
	errJson := json.Unmarshal(body, &value)
	if errJson != nil {
		log.Fatal(err)
	}
	return value
}

func check(results chan string, route Route, key *string) {
	all := get(route.Url)
	value := all[*key]
	result := fmt.Sprintf("[%s %5s] %s = %v", route.Country, route.Env, *key, value)
	results <- result
}

func main() {
	client = http.Client{
		Timeout: time.Second * 5, // Timeout after 2 seconds
	}

	key := flag.String("key", "", "key in site params to check")

	flag.Parse()

	if *key == "" {
		log.Fatal("empty key, provide -key option")
	}

	jsonFile, err := os.Open("./routes.json")
	if err != nil {
		fmt.Println(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var routes []Route

	json.Unmarshal(byteValue, &routes)

	results := make(chan string)

	// Do something for each route
	for _, route := range routes {
		go check(results, route, key)
	}

	for i := 0; i < len(routes); i++ {
		result := <-results
		fmt.Println(result)
	}

	defer jsonFile.Close()
}
