package parsers

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/rainingclouds/lemonades/logger"
	"log"
	"net/url"
	"strconv"
	"strings"
)

func parseSnapdeal(url *url.URL) (*Item, error) {
	item := new(Item)
	item.ProductLink = strings.Split(url.String(), "/ref")[0]
	doc, err := goquery.NewDocument(url.String())
	if err != nil {
		return nil, err
	}
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		property, ok := s.Attr("property")
		if ok {
			content, ok := s.Attr("content")
			if ok {
				switch property {
				case "og:title":
					item.Name = content
					log.Println(content)
				case "og:image":
					item.ProductImage = content
					log.Println(content)
				}
			}
		}
	})
	doc.Find("#selling-price-id").Each(func(i int, s *goquery.Selection) {
		log.Println(s.Text())
		item.PriceValue, err = strconv.ParseInt(s.Text(), 10, 64)
		if err != nil {
			logger.Err("Error while parsing", url.String(), err)
		}
	})
	item.PriceCurrency = "Rs"
	doc.Find(".containerBreadcrumb").Children().Children().Each(func(i int, s *goquery.Selection) {

		switch i {
		case 1:
			log.Println("Main", s.Text())
			item.MainCategory = strings.TrimSpace(s.Text())
		case 2:
			log.Println("Sub", s.Text())
			item.SubCategory = strings.TrimSpace(s.Text())
		}
	})
	return item, nil
}
