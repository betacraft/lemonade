package parsers

import (
	"errors"
	"log"
	"net/url"
)

func Parse(rawUrl string) (*Item, error) {
	uri, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}
	log.Println(uri.Host)
	switch uri.Host {
	case "www.flipkart.com":
		fallthrough
	case "flipkart.com":
		return parseFlipkart(uri)
	case "www.amazon.in":
		fallthrough
	case "amazon.in":
		return parseAmazon(uri)
	default:
		return nil, errors.New("Please use flipkart/amazon urls")
	}
}
