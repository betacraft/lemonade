package parsers

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/url"
	"strings"
)

const (
	CATALOG_BASE_URL = "https://catalog.paytm.com/v1/p/"
)

func parsePaytm(url *url.URL) (*Item, error) {
	item := new(Item)
	item.ProductLink = strings.Split(url.String(), "?")[0]
	components := strings.Split(url.String(), "/")
	catalogUrl := CATALOG_BASE_URL + components[len(components)-1]
	log.Println(catalogUrl)
	doc, err := goquery.NewDocument(catalogUrl)
	if err != nil {
		return nil, err
	}
	data := map[string]interface{}{}
	if err := json.Unmarshal([]byte(doc.Text()), &data); err != nil {
		panic(err)
	}
	item.Name = data["name"].(string)
	item.ProductImage = data["image_url"].(string)
	item.MainCategory = data["vertical_label"].(string)
	item.SubCategory = ""
	item.PriceValue = int64(data["offer_price"].(float64))
	item.PriceCurrency = "Rs"
	item.Attributes = map[string]string{}
	for key, value := range data["attributes"].(map[string]interface{}) {
		item.Attributes[key] = fmt.Sprintf("%v", value)
	}
	return item, nil
}
