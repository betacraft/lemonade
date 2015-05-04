package models

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/rainingclouds/lemonades/db"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/url"
	"strconv"
	"strings"
)

//go:generate easytags Product.go bson
//go:generate easytags Product.go json

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

	Description string            `json:"description" bson:"description"`
	Attributes  map[string]string `json:"attributes" bson:"attributes"`

	Timestamp
}

func (p *Product) Create() error {
	p.Id = bson.NewObjectId()
	return db.MgCreateStrong(C_PRODUCT, p)
}

func (p *Product) Update() error {
	return db.MgUpdateStrong(C_PRODUCT, p.Id, p)
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
