package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

func getenv(name string) (string, error) {
	v := os.Getenv(name)
	if v == "" {
		return v, errors.New("no environment variable: " + name)
	}
	return v, nil
}

func getRSS(rssFeed string) ([]string, error) {
	if rssFeed == "" {
		return []string{""}, errors.New("no feeds present")
	}
	return strings.Split(rssFeed, ";"), nil
}

func makeReadme(filename string) error {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://jethromay.dev/rss.xml")
	if err != nil {
		log.Fatalf("error getting feed: %v", err)
	}
	// Get latest items
	rssItem := feed.Items[0]

	// Unwrap Markdown content
	content, err := ioutil.ReadFile("static.md")
	if err != nil {
		log.Fatalf("cannot read file: %v", err)
		return err
	}
	stringyContent := string(content)

	date := time.Now().Format("2 Jan 2006")

	// Combine static and dymanic content
	blog := "- ✨ Read my latest blog post: **[" + rssItem.Title + "](" + rssItem.Link + ")**"
	updated := "Last updated on " + date + "."
	data := fmt.Sprintf("%s%s\n\n%s\n", stringyContent, blog, updated)

	// Prepare file 
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	return file.Sync()
}

func main() {

	makeReadme("../README.md")

}
