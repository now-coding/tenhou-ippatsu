package main

import (
	"compress/gzip"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	files := get_log_file_list()
	for _, file := range files {
		html := get_html_by_file(file)
		ids := get_paifu_ids_from_html(html)

		for _, id := range ids {
			paifu := get_paifu(id)
			log.Println(paifu)
			break
		}
		break
	}
}

func get_paifu(id string) string {
	url := "https://tenhou.net/0/log/?" + id
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}

func get_paifu_ids_from_html(html *goquery.Document) []string {
	ids := []string{}

	html.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			r := regexp.MustCompile(`log=([\w\-]+)`)
			matches := r.FindStringSubmatch(href)
			if len(matches) > 0 {
				ids = append(ids, matches[1])
			}
		}
	})

	return ids
}

func get_html_by_file(file string) *goquery.Document {
	url := "https://tenhou.net/sc/raw/dat/" + file
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	reader, err := gzip.NewReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func get_log_file_list() []string {
	res, err := http.Get("https://tenhou.net/sc/raw/list.cgi")
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	r := regexp.MustCompile(`\w+\.html\.gz`)
	files := []string{}
	for _, matches := range r.FindAllStringSubmatch(string(body), -1) {
		files = append(files, matches[0])
	}

	return files
}
