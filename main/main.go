package main

import (
	"fmt"
	"net/http"
	"time"
)

type UrlResponse struct {
	URL      string
	isActive bool
}

const DEFAULT_URL_1 string = "https://vnexpress.net/"
const DEFAULT_URL_2 string = "https://viblo.asia/p/crawl-du-lieu-voi-golang-V3m5WbLglO7"
const DEFAULT_URL_3 string = "https://viblo.asia/p/crawl-data-trong-golang-voi-goquery-LzD5dNoEZjY"
const DEFAULT_URL_4 string = "https://dev.to/vianeltxt/how-to-build-a-web-scraper-using-golang-with-colly-18lh"
const DEFAULT_URL_5 string = "https://tuoitre.vn/"

func main() {
	fakeData := generateFakeData()
	start := time.Now()

	pingWithLimitedGoroutines(fakeData)
	//pingWithConcurrent(fakeData)
	//pingWithoutConcurrent(fakeData)

	duration := time.Since(start)
	fmt.Printf("total: %f seconds\n", duration.Seconds())

}

func generateFakeData() []string {
	urls := []string{DEFAULT_URL_1, DEFAULT_URL_2, DEFAULT_URL_3, DEFAULT_URL_4, DEFAULT_URL_5}
	fakeData := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		index := i % 5
		fakeData[i] = urls[index]
	}
	return fakeData
}

func pingWithLimitedGoroutines(fakeData []string){
	maxGoroutines := 5
	guard := make(chan struct{}, maxGoroutines)
	resultChannel := make(chan UrlResponse)

	for i := 0; i < len(fakeData); i++ {
		guard <- struct{}{} // would block if guard channel is already filled
		go func(url string) {
			checkWebsite(url, resultChannel)
			<-guard
		}(fakeData[i])
	}

	for i:=0; i < len(fakeData); i++{
		result := <-resultChannel
		fmt.Printf("%s -- status: %v\n", result.URL, result.isActive)
	}
}

func pingWithConcurrent(fakeData []string) {
	resultChannel := make(chan UrlResponse)
	for _, url := range fakeData {
		checkWebsite(url, resultChannel)
	}

	for i := 0; i < len(fakeData); i++ {
		result := <-resultChannel
		fmt.Printf("%s -- status: %v\n", result.URL, result.isActive)
	}
}

func checkWebsite(url string, c chan UrlResponse) {
	go func() {
		resp, err := http.Get(url)
		if resp != nil {
			c <- UrlResponse{url, true}
		} else if err != nil {
			c <- UrlResponse{url, false}
		}
	}()
}

func pingWithoutConcurrent(fakeData []string){
	for i := 0; i < len(fakeData); i++ {
		result, err := http.Get(fakeData[i])
		if result != nil{
			fmt.Printf("%s -- status: %v\n", fakeData[i], true)
		} else if err != nil{
			fmt.Printf("%s -- status: %v\n", fakeData[i], false)
		}

	}
}
