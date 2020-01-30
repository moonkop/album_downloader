package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	urlFile   = "C:\\Users\\Admin\\Downloads\\urls_drawings.txt"
	path      = "C:\\Users\\Admin\\Desktop\\rbq\\download\\drawings\\"
	proxyUrl  = "http://localhost:1080"
	parallels = 20
)

var downloading = make(chan int, parallels)

func main() {
	urls, err := ioutil.ReadFile(urlFile)
	if err != nil {
		return
	}
	str := string(urls)
	arr := strings.SplitN(str, "\n", -1)
	r, err := regexp.Compile("[^/]+$")
	for _, url := range arr {
		path := path + r.FindString(url)
		downloading <- 1
		go downloadFileIfNotExist(path, url)

	}
	//fmt.Println("all finished")
	fmt.Scan(&str)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

var purl, err = url.Parse(proxyUrl)

var PTransport = &http.Transport{
	Proxy: http.ProxyURL(purl),
}

func downloadFileIfNotExist(filepath string, url string) (err error) {
	defer func() { <-downloading }()

	if fileExists(filepath) {
		return
	}
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	client := http.Client{
		Transport: PTransport,
		Timeout:   20 * time.Second,
	}
	fmt.Println("downloading", url)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.27 Safari/537.36`)

	resp, err := client.Do(req)
	// Get the data
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	fmt.Println("finished", url)

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		<-downloading
		return err
	}

	return nil
}
