package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/karanbirsingh7/go-hackernews/hn"
)

var (
	port       int
	numStories int
	wg         sync.WaitGroup
)

func parseFlags() {
	flag.IntVar(&port, "port", 8080, "port to start web server on")
	flag.IntVar(&numStories, "num_stories", 30, "number of top stories to display")
	flag.Parse()
}

func main() {
	// parse Flags
	parseFlags()

	// main tempalte
	tpl := template.Must(template.ParseFiles("./index.gohtml"))

	// ROUTES
	http.HandleFunc("/", homeHandler(numStories, tpl))

	// start server
	bindAddress := fmt.Sprintf(":%d", port)
	fmt.Println("Server running on", bindAddress)
	log.Fatal(http.ListenAndServe(bindAddress, nil))
}

func homeHandler(numStories int, tpl *template.Template) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		startTime := time.Now()
		// call function to get all items
		stories, err := getTopStories()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		endTime := time.Since(startTime)
		data := templateData{
			Stories: stories,
			Time:    endTime,
		}

		// template the html
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})

}

func getTopStories() (stories []item, err error) {
	var client hn.Client

	items, err := client.TopStories()
	if err != nil {
		errMessage := fmt.Sprintf("Failed to laod top stories - %v", err)
		return nil, errors.New(errMessage)
	}

	// make struct to store results from goroutine
	type result struct {
		item item
		err  error
	}
	// initialize a new channel
	resultsCh := make(chan result)

	// get each item info
	for i := 0; i < numStories; i++ {
		// loop through each item ID: ex 2435
		go func(id int) {
			hnItem, err := client.GetItem(id)
			if err != nil {
				resultsCh <- result{
					err: err,
				}
			}
			// give results back to channel
			resultsCh <- result{
				item: parseHNItem(hnItem),
			}
		}(items[i])
	}

	// pull all values from channel
	var results []result
	for i := 0; i < numStories; i++ {
		results = append(results, <-resultsCh)
	}

	// convert back to stories type
	for _, res := range results {
		if res.err != nil {
			continue
		}

		if isStoryLink(res.item) {
			stories = append(stories, res.item)
		}
	}

	return stories, nil
}

// check if story or job posting or else
func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

// item is same as hn.Item with additional host field
type item struct {
	hn.Item
	Host string
}

type templateData struct {
	Stories []item
	Time    time.Duration
}
