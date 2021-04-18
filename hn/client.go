package hn

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const apiBase = "https://hacker-news.firebaseio.com/v0"

// use Client to interact with HN API
type Client struct {
	apiBase string
}

// AutoInitiate client base URL if nothing passed by user
func (c *Client) defaultify() {
	if c.apiBase == "" {
		c.apiBase = apiBase
	}
}

func (c *Client) TopStories() (ids []int, err error) {
	c.defaultify()

	req_url := fmt.Sprintf("%s/topstories.json", c.apiBase)
	resp, err := http.Get(req_url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &ids); err != nil {
		return nil, err
	}

	return
}

func (c *Client) GetItem(id int) (item Item, err error) {
	c.defaultify()

	req_url := fmt.Sprintf("%s/item/%d.json", c.apiBase, id)
	resp, err := http.Get(req_url)

	if err != nil {
		return item, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return item, err
	}

	if err := json.Unmarshal(body, &item); err != nil {
		return item, err
	}
	return
}

// sample item from HN: https://hacker-news.firebaseio.com/v0/item/26681984.json
type Item struct {
	ID          int    `json:"id"`
	Kids        []int  `json:"kids"`
	Score       int    `json:"score"`
	Time        int    `json:"time"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	URL         string `json:"url"`
	Descendants int    `json:"descendants"`
	By          string `json:"by"`
}
