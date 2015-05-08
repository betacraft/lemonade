package models

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/rainingclouds/lemonades/db"
	"github.com/rainingclouds/lemonades/logger"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

//go:generate easytags Product.go bson
//go:generate easytags Product.go json

type Price struct {
	Date          time.Time `json:"date" bson:"date"`
	PriceValue    int64     `json:"price_value" bson:"price_value"`
	PriceCurrency string    `json:"price_currency" bson:"price_currency"`
}

type Product struct {
	Id      bson.ObjectId   `json:"id" bson:"_id,omitempty"`
	AddedBy []bson.ObjectId `json:"added_by" bson:"added_by"`

	ProductLink  string `json:"product_link" bson:"product_link"`
	ProductImage string `json:"product_image" bson:"product_image"`

	Name         string `json:"name" bson:"name"`
	MainCategory string `json:"main_category" bson:"main_category"`
	SubCategory  string `json:"sub_category" bson:"sub_category"`

	PriceValue    int64  `json:"price_value" bson:"price_value"`
	PriceCurrency string `json:"price_currency" bson:"price_currency"`

	PriceHistory []Price `json:"price_history" bson:"price_history"`

	Description string            `json:"description" bson:"description"`
	Attributes  map[string]string `json:"attributes" bson:"attributes"`

	Timestamp
}

func UpdateProductPrices() {
	pageNo := 0
	products := new([]Product)
	var i int
	err := db.MgFindPage(C_PRODUCT, &bson.M{}, pageNo, products)
	if err != nil {
		logger.Err("Error while updating price of products", err)
	}
	for len(*products) > 0 {
		for i = 0; i < len(*products); i++ {
			err = (*products)[i].UpdatePrice()
			if err != nil {
				logger.Err("While saving updated price", (*products)[i], err)
				continue
			}
		}
		pageNo = pageNo + 1
		err = db.MgFindPage(C_PRODUCT, &bson.M{}, pageNo, products)
		if err != nil {
			logger.Err("Error while updating price of products page no %d", pageNo, err)
			return
		}
	}
}

func (p *Product) UpdatePrice() error {
	price, err := FetchProductPrice(p.ProductLink)
	if err != nil {
		logger.Err("While fetching price", p, err)
		return err
	}
	p.PriceHistory = append(p.PriceHistory, *price)
	p.PriceValue = price.PriceValue
	err = p.Update()
	if err != nil {
		return err
	}
	return UpdateProductInfo(p)
}

func (p *Product) Create() error {
	p.Id = bson.NewObjectId()
	return db.MgCreateStrong(C_PRODUCT, p)
}

func (p *Product) Update() error {
	return db.MgUpdateStrong(C_PRODUCT, p.Id, p)
}

func GetProductById(id bson.ObjectId) (*Product, error) {
	product := new(Product)
	err := db.MgFindOneStrong(C_PRODUCT, &bson.M{"_id": id}, product)
	return product, err
}

func (p *Product) CreateOrUpdate() error {
	if p.Id.Hex() == "" {
		log.Println("Creating a new product")
		return p.Create()
	}
	return p.Update()
}

func GetProductByProductLink(link string) (*Product, error) {
	log.Println("Finding for", link)
	product := new(Product)
	err := db.MgFindOneStrong(C_PRODUCT, &bson.M{"product_link": link}, product)
	return product, err
}

func FetchProductPrice(rawUrl string) (*Price, error) {
	uri, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}
	log.Println(uri.Host)
	switch uri.Host {
	case "www.flipkart.com":
		fallthrough
	case "flipkart.com":
		var priceValue int64
		doc, err := goquery.NewDocument(rawUrl)
		if err != nil && err.Error() != "not found" {
			return nil, err
		}
		doc.Find(".seller-table-wrap").Each(func(i int, s *goquery.Selection) {
			dataConfig, ok := s.Attr("data-config")
			if ok && strings.Contains(dataConfig, "sellingPrice") {
				priceValue, _ = strconv.ParseInt(strings.Split(strings.Split(dataConfig, "\"sellingPrice\":")[1], ",")[0], 10, 64)
			}
		})
		if priceValue == 0 {
			doc.Find(".selling-price").Each(func(i int, s *goquery.Selection) {
				if i == 0 {
					val := strings.TrimSpace(strings.Split(s.Text(), "Rs.")[1])
					priceValue, _ = strconv.ParseInt(strings.Replace(val, ",", "", -1), 10, 64)
				}
			})
		}
		if priceValue == 0 {
			html, _ := doc.Html()
			return nil, errors.New("Error in parsing flipkart html" + html)
		}
		log.Println(fmt.Sprintf("%d Rs", priceValue))
		price := new(Price)
		price.Date = time.Now().UTC()
		price.PriceCurrency = "Rs"
		price.PriceValue = priceValue
		return price, nil
	case "amazon.in":
		fallthrough
	case "www.amazon.in":
		var priceValue int64
		doc, err := goquery.NewDocument(rawUrl)
		if err != nil && err.Error() != "not found" {
			return nil, err
		}
		log.Printf("Fething info")
		doc.Find("#fbt_item_data").Each(func(i int, s *goquery.Selection) {
			priceValue, _ = strconv.ParseInt(strings.Split(strings.Split(s.Text(), "buyingPrice\":")[1], ",")[0], 10, 64)
		})
		log.Println(fmt.Sprintf("%d Rs", priceValue))
		if priceValue == 0 {
			html, _ := doc.Html()
			return nil, errors.New("Error in parsing amazon.in html" + html)
		}
		price := new(Price)
		price.Date = time.Now().UTC()
		price.PriceCurrency = "Rs"
		price.PriceValue = priceValue
		return price, nil
	}
	return nil, errors.New("Illegal url")
}

func FetchProductInfo(rawurl string) (*Product, error) {
	uri, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	log.Println(uri.Host)
	switch uri.Host {
	case "www.flipkart.com":
		fallthrough
	case "flipkart.com":
		product, err := GetProductByProductLink(strings.Split(rawurl, "?")[0])
		if err != nil && err.Error() != "not found" {
			return nil, err
		}
		log.Println(product)
		if product.Id.Hex() != "" {
			return product, nil
		}
		product = new(Product)
		product.ProductLink = strings.Split(rawurl, "?")[0]
		doc, err := goquery.NewDocument(rawurl)
		if err != nil {
			return nil, err
		}
		doc.Find("meta").Each(func(i int, s *goquery.Selection) {
			property, ok := s.Attr("property")
			if ok {
				content, ok := s.Attr("content")
				if ok {
					switch property {
					case "og:image":
						product.ProductImage = content
						log.Println(content)
					}
				}
			}
		})
		doc.Find(".title").Each(func(i int, s *goquery.Selection) {
			itemprop, ok := s.Attr("itemprop")
			if ok {
				if itemprop == "name" {
					log.Println(s.Text())
					product.Name = s.Text()
				}
			}
		})
		doc.Find(".seller-table-wrap").Each(func(i int, s *goquery.Selection) {
			dataConfig, ok := s.Attr("data-config")
			if ok && strings.Contains(dataConfig, "sellingPrice") {
				product.PriceValue, _ = strconv.ParseInt(strings.Split(strings.Split(dataConfig, "\"sellingPrice\":")[1], ",")[0], 10, 64)
				product.PriceCurrency = "Rs"
			}
		})
		if product.PriceValue == 0 {
			doc.Find(".selling-price").Each(func(i int, s *goquery.Selection) {
				if i == 0 {
					val := strings.TrimSpace(strings.Split(s.Text(), "Rs.")[1])
					product.PriceValue, _ = strconv.ParseInt(strings.Replace(val, ",", "", -1), 10, 64)
					product.PriceCurrency = "Rs"
				}
			})
		}
		log.Println(fmt.Sprintf("%d Rs", product.PriceValue))
		doc.Find(".clp-breadcrumb").Children().Children().Children().Each(func(i int, s *goquery.Selection) {
			switch i {
			case 1:
				product.MainCategory = strings.TrimSpace(s.Text())
			case 2:
				product.SubCategory = strings.TrimSpace(s.Text())
			}
		})
		return product, nil
	case "www.amazon.in":
		fallthrough
	case "amazon.in":
		product, err := GetProductByProductLink(strings.Split(rawurl, "/ref")[0])
		if err != nil && err.Error() != "not found" {
			return nil, err
		}
		log.Println(product)
		if product.Id.Hex() != "" {
			return product, nil
		}
		product.ProductLink = strings.Split(rawurl, "/ref")[0]
		doc, err := goquery.NewDocument(rawurl)
		if err != nil {
			return nil, err
		}
		doc.Find("#productTitle").Each(func(i int, s *goquery.Selection) {
			log.Println(s.Text())
			product.Name = s.Text()
		})
		components := strings.Split(strings.Split(rawurl, "/ref")[0], "/")
		product.ProductImage = "http://images.amazon.com/images/P/" + components[len(components)-1] + ".jpg"
		log.Println(product.ProductImage)
		doc.Find("#fbt_item_data").Each(func(i int, s *goquery.Selection) {
			product.PriceValue, _ = strconv.ParseInt(strings.Split(strings.Split(s.Text(), "buyingPrice\":")[1], ",")[0], 10, 64)
			product.PriceCurrency = "Rs"
		})
		log.Println(fmt.Sprintf("%d Rs", product.PriceValue))
		doc.Find("#wayfinding-breadcrumbs_feature_div").Children().Children().Children().Children().Each(func(i int, s *goquery.Selection) {
			log.Println(strings.TrimSpace(s.Text()))
			switch i {
			case 1:
				product.MainCategory = strings.TrimSpace(s.Text())
			case 2:
				product.SubCategory = strings.TrimSpace(s.Text())
			}
		})
		return product, nil
	default:
		return nil, errors.New("Please use flipkart/amazon urls")
	}
	return nil, nil
}
