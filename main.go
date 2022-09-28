package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const (
	css = `href=["|'][^"]+\.css["|']`
	js  = `src=['|"][^"]+\.js['|"]`
)

func main() {

	os.RemoveAll("result")
	if err := os.MkdirAll("result", os.ModePerm); err != nil {
		log.Fatal(err)
		return
	}

	var link string
	fmt.Println("What website do you want to eat?")
	if _, err := fmt.Scan(&link); err != nil {
		fmt.Println("invalid website link")
	}

	body, _ := getFromLink(link)
	ioutil.WriteFile("./result/initial-dump", body, 0600)

	bodyStr := string(body)
	index := bodyStr

	// css
	cssRegex, err := regexp.Compile(css)
	matchCss := cssRegex.FindAllString(bodyStr, 100)

	for idx, match := range matchCss {
		match = match[6 : len(match)-1]
		fmt.Println("Match: ", match, " Error: ", err)
		cssBody, err := getFromLink(match)
		if err != nil {
			continue
		}
		path := fmt.Sprintf(`styles/style_%d.css`, idx)
		if len(matchCss) < 2 {
			path = `styles/style.css`
		}
		totalPath := fmt.Sprintf("./result/%s", path)

		if err := os.MkdirAll("result/styles", os.ModePerm); err != nil {
			log.Fatal(err)
		}
		ioutil.WriteFile(totalPath, cssBody, 0600)

		relPath := fmt.Sprintf(`./%s`, path)
		index = strings.ReplaceAll(index, match, relPath)
	}

	//js
	jsRegex, err := regexp.Compile(js)
	matchJs := jsRegex.FindAllString(bodyStr, 100)

	for idx, match := range matchJs {
		match = match[5 : len(match)-1]
		fmt.Println("Match: ", match, " Error: ", err)
		jsBody, err := getFromLink(match)
		if err != nil {
			continue
		}
		path := fmt.Sprintf(`script/script_%d.js`, idx)
		if len(matchJs) < 2 {
			path = `script/script.css`
		}
		totalPath := fmt.Sprintf("./result/%s", path)

		if err := os.MkdirAll("result/script", os.ModePerm); err != nil {
			log.Fatal(err)
		}
		ioutil.WriteFile(totalPath, jsBody, 0600)

		relPath := fmt.Sprintf(`./%s`, path)
		index = strings.ReplaceAll(index, match, relPath)
	}

	ioutil.WriteFile("./result/index.html", []byte(index), 0600)
}

func getFromLink(link string) ([]byte, error) {
	resp, err := http.Get(link)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	return body, err
}
