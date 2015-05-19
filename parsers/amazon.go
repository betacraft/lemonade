package parsers

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/url"
	"strconv"
	"strings"
)

func parseAmazon(url *url.URL) (*Item, error) {
	item := new(Item)
	item.ProductLink = strings.Split(url.String(), "/ref")[0]
	doc, err := goquery.NewDocument(url.String())
	if err != nil {
		return nil, err
	}
	doc.Find("#productTitle").Each(func(i int, s *goquery.Selection) {
		log.Println(s.Text())
		item.Name = s.Text()
	})
	components := strings.Split(strings.Split(url.String(), "/ref")[0], "/")
	item.ProductImage = "http://images.amazon.com/images/P/" + components[len(components)-1] + ".jpg"
	log.Println(item.ProductImage)
	doc.Find("#fbt_item_data").Each(func(i int, s *goquery.Selection) {
		item.PriceValue, _ = strconv.ParseInt(strings.Split(strings.Split(s.Text(), "buyingPrice\":")[1], ",")[0], 10, 64)
		item.PriceCurrency = "Rs"
	})
	log.Println(fmt.Sprintf("%d Rs", item.PriceValue))
	doc.Find("#wayfinding-breadcrumbs_feature_div").Children().Children().Children().Children().Each(func(i int, s *goquery.Selection) {
		log.Println(strings.TrimSpace(s.Text()))
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
	done := false
	doc.Find(".pdClearfix").Children().Children().Children().Children().Each(func(i int, s *goquery.Selection) {
		s.Children().Each(func(i int, s *goquery.Selection) {
			if s.Children().HasClass("groupHead") || done {
				return
			}
			s.Children().Each(func(i int, s *goquery.Selection) {
				if done {
					return
				}
				if s.HasClass("label") {
					key, _ = s.Html()
					key = strings.TrimSpace(key)
					log.Println(key)
				}
				if s.HasClass("value") {
					value, _ = s.Html()
					value = strings.TrimSpace(value)
					log.Println(value)
				}
				item.Attributes[key] = value
				if strings.TrimSpace(key) == "Customer Reviews" {
					log.Println("done")
					done = true
				}
			})
		})
	})
	return item, nil
}
