package parsers

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/url"
	"strconv"
	"strings"
)

func parseFlipkart(url *url.URL) (*Item, error) {
	item := new(Item)
	item.ProductLink = strings.Split(url.String(), "?")[0]
	doc, err := goquery.NewDocument(url.String())
	if err != nil {
		return nil, err
	}
	doc.Find(".title").Each(func(i int, s *goquery.Selection) {
		itemprop, ok := s.Attr("itemprop")
		if ok {
			if itemprop == "name" {
				log.Println(s.Text())
				item.Name = s.Text()
			}
		}
	})
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		property, ok := s.Attr("property")
		if ok {
			content, ok := s.Attr("content")
			if ok {
				switch property {
				case "og:image":
					item.ProductImage = content
					log.Println(content)
				}
			}
		}
	})
	doc.Find(".selling-price").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			val := strings.TrimSpace(strings.Split(s.Text(), "Rs.")[1])
			item.PriceValue, _ = strconv.ParseInt(strings.Replace(val, ",", "", -1), 10, 64)
			item.PriceCurrency = "Rs"
		}
	})
	if item.PriceValue == 0 {
		doc.Find(".seller-table-wrap").Each(func(i int, s *goquery.Selection) {
			dataConfig, ok := s.Attr("data-config")
			if ok && strings.Contains(dataConfig, "sellingPrice") {
				item.PriceValue, _ = strconv.ParseInt(strings.Split(strings.Split(dataConfig, "\"sellingPrice\":")[1], ",")[0], 10, 64)
				item.PriceCurrency = "Rs"
			}
		})
	}
	log.Println(fmt.Sprintf("%d Rs", item.PriceValue))
	doc.Find(".clp-breadcrumb").Children().Children().Children().Each(func(i int, s *goquery.Selection) {
		switch i {
		case 1:
			item.MainCategory = strings.TrimSpace(s.Text())
		case 2:
			item.SubCategory = strings.TrimSpace(s.Text())
		}
	})
	item.Attributes = map[string]string{}
	var key string
	var value string
	doc.Find(".specSection").Children().Children().Each(func(i int, s *goquery.Selection) {
		s.Children().Each(func(i int, s *goquery.Selection) {
			if s.Children().HasClass("groupHead") {
				return
			}
			s.Children().Each(func(i int, s *goquery.Selection) {
				if s.HasClass("specsKey") {
					key, _ = s.Html()
					key = strings.TrimSpace(key)
					log.Println(key)
				}
				if s.HasClass("specsValue") {
					if strings.TrimSpace(s.Children().Text()) != "" {
						return
					}
					value, _ = s.Html()
					value = strings.TrimSpace(value)
					log.Println(strings.TrimSpace(value))
				}
				item.Attributes[key] = value
			})
		})
	})
	return item, nil
}
