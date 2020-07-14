package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	articlesUrl = "https://storage.googleapis.com/aller-structure-task/articles.json"
	marketingUrl = "https://storage.googleapis.com/aller-structure-task/contentmarketing.json"
	port = ":8081"
	pattern = "/"
)

type Item struct {
	CerebroScore      float64 `json:"cerebro-score,omitempty"`
	CleanImage        string  `json:"cleanImage,omitempty"`
	CommercialPartner string  `json:"commercialPartner,omitempty"`
	HarvesterID       string  `json:"harvesterId,omitempty"`
	LogoURL           string  `json:"logoURL,omitempty"`
	Title             string  `json:"title,omitempty"`
	Type              string  `json:"type,omitempty"`
	URL               string  `json:"url,omitempty"`
}

type ContentmarketingJson struct {
	HttpStatus int64 `json:"httpStatus,omitempty"`
	Response   struct {
		Items []Item `json:"items,omitempty"`
	} `json:"response,omitempty"`
}

func getJsonByUrl(url string) ContentmarketingJson {
	req, err := http.Get(url)
	if err != nil {
		fmt.Println("could not get ", url)
		return ContentmarketingJson{}
	}
	//be sure that will be closed after all(must be closed)
	defer req.Body.Close()
	//read all bytes
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("unable to extract body")
	}

	res := ContentmarketingJson {}
	if err := json.Unmarshal(body, &res); err != nil {
		fmt.Println("unable to parse json bytes")
	}

	return res
}

func getCompiledJsons(firstUrl, secondUrl string) (out string){
	articles := getJsonByUrl(firstUrl).Response.Items
	marketing := getJsonByUrl(secondUrl).Response.Items

	//simulation of deque to pop easy
	var marketingDeque []Item
	for i := len(marketing) - 1; i >= 0; i-- {
		marketingDeque = append(marketingDeque, marketing[i])
	}

	//create Item with only Ad type value
	defaultMarketing := Item{}
	defaultMarketing.Type = "Ad"

	var  resultArr []Item
	for i := 0; i < len(articles); i++ {
		//append every article
		resultArr = append(resultArr, articles[i])
		//if already appended 5 articles need to append marketing item
		if (i + 1) % 5 == 0 {
			//append provided by second link marketing item if available and pop it from
			//marketingDeque, else if provided links are gone append defaultMarketing
			if len(marketingDeque) != 0 {
				top := len(marketingDeque) - 1
				resultArr = append(resultArr, marketingDeque[top])
				marketingDeque = marketingDeque[:top]
			} else  {
				resultArr = append(resultArr, defaultMarketing)
			}
		}
	}
	//serialize struct to bytes with intends
	response, err := json.MarshalIndent(resultArr, "", " ")
	if err != nil {
		fmt.Println("Error while parsing final json")
	}
	return string(response)
}

func home(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, getCompiledJsons(articlesUrl, marketingUrl))
}

func handleReq() {
	http.HandleFunc(pattern, home)
	log.Fatal(http.ListenAndServe(port, nil))
}

func main() {
	handleReq()
}