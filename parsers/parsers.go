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
	case "www.snapdeal.com":
		fallthrough
	case "snapdeal.com":
		return parseSnapdeal(uri)
	case "www.paytm.com":
		fallthrough
	case "paytm.com":
		return parsePaytm(uri)
	default:
		return nil, errors.New("Please use flipkart/amazon urls")
	}
}
