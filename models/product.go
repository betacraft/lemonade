package models

import (
	"bytes"
	"fmt"
	"github.com/AdRoll/goamz/s3"
	"github.com/disintegration/imaging"
	"github.com/rainingclouds/lemonades/aws"
	"github.com/rainingclouds/lemonades/db"
	"github.com/rainingclouds/lemonades/logger"
	"github.com/rainingclouds/lemonades/parsers"
	"gopkg.in/mgo.v2/bson"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
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

type Attribute struct {
	Name  string `json:"name" bson:"name"`
	Value string `json:"value" bson:"value"`
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

	Description string      `json:"description" bson:"description"`
	Specs       []Attribute `json:"specs" bson:"specs"`

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

func UpdateAllProductInfo() {
	pageNo := 0
	products := new([]Product)
	var i int
	err := db.MgFindPage(C_PRODUCT, &bson.M{}, pageNo, products)
	if err != nil {
		logger.Err("Error while updating price of products", err)
	}
	for len(*products) > 0 {
		for i = 0; i < len(*products); i++ {
			err = (*products)[i].UpdateInfo()
			if err != nil {
				logger.Err("Error while updating info", err)
				continue
			}
			err = UpdateProductInfo(&(*products)[i])
			if err != nil {
				logger.Err("Error while updating info", err)
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

func (p *Product) UpdateInfo() error {
	item, err := parsers.Parse(p.ProductLink)
	if err != nil {
		return err
	}
	p.Name = item.Name
	p.ProductLink = item.ProductLink
	p.ProductImage = item.ProductImage
	p.MainCategory = item.MainCategory
	p.SubCategory = item.SubCategory
	p.PriceValue = item.PriceValue
	p.PriceCurrency = item.PriceCurrency
	price := new(Price)
	price.Date = time.Now().UTC()
	price.PriceValue = item.PriceValue
	price.PriceCurrency = item.PriceCurrency
	p.PriceHistory = []Price{}
	p.PriceHistory = append(p.PriceHistory, *price)
	p.Description = item.Description
	p.Specs = []Attribute{}
	for key, value := range item.Attributes {
		p.Specs = append(p.Specs, Attribute{Name: key, Value: value})
	}
	return p.Update()
}

func (p *Product) UpdatePrice() error {
	price, err := FetchProductPrice(p.ProductLink)
	if err != nil {
		logger.Err("While fetching price", p, err)
		return err
	}
	if price.PriceValue == p.PriceValue {
		return nil
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
	go p.ProcessImage()
	return db.MgCreateStrong(C_PRODUCT, p)
}

func (p *Product) ProcessImage() {
	client := http.Client{}
	urlParts := strings.Split(p.ProductImage, "/")
	imgName := urlParts[len(urlParts)-1]
	logger.Debug("Image name is", imgName)
	imgResponse, err := client.Get(p.ProductImage)
	if err != nil {
		logger.Err("Error while downloading the product image", err)
		return
	}
	logger.Debug("Content length", imgResponse.ContentLength)
	imgData, err := ioutil.ReadAll(imgResponse.Body)
	if err != nil {
		logger.Err("Error while copying the product image", err)
		return
	}
	imgResponse.Body.Close()
	err = aws.Bucket().Put(fmt.Sprintf("product/%v/default/%v", p.Id.Hex(), imgName), imgData, "image/jpeg", s3.PublicRead, s3.Options{})
	if err != nil {
		logger.Err("Error while uploading the basic image", err)
		return
	}
	// resizing image to 300x300
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		logger.Err("Error while converting bytes to image", err)
		return
	}
	img300x300 := imaging.Resize(img, 300, 0, imaging.Lanczos)
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, img300x300, nil)
	err = aws.Bucket().Put(fmt.Sprintf("product/%v/300x300/%v", p.Id.Hex(), imgName), buf.Bytes(), "image/jpeg", s3.PublicRead, s3.Options{})
	if err != nil {
		logger.Err("Error while uploading the 300x300 image", err)
		return
	}
	p.ProductImage = aws.Bucket().URL(fmt.Sprintf("product/%v/300x300/%v", p.Id.Hex(), imgName))
	p.ProductImage = strings.Replace(p.ProductImage, "https", "http", -1)
	err = p.Update()
	if err != nil {
		logger.Err("Error while saving product while processing image", err)
		return
	}
	err = UpdateProductInfo(p)
	if err != nil {
		logger.Err("Error while saving product while processing image", err)
		return
	}
}

func (p *Product) Update() error {
	return db.MgUpdateStrong(C_PRODUCT, p.Id, p)
}

func GetProductByName(name string) (*Product, error) {
	product := new(Product)
	err := db.MgFindOneStrong(C_PRODUCT, &bson.M{"name": name}, product)
	return product, err
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
	item, err := parsers.Parse(rawUrl)
	if err != nil {
		return nil, err
	}
	price := new(Price)
	price.Date = time.Now().UTC()
	price.PriceValue = item.PriceValue
	price.PriceCurrency = item.PriceCurrency
	return price, nil
}

func FetchProductInfo(rawurl string) (*Product, error) {
	product, err := GetProductByProductLink(strings.Split(rawurl, "?")[0])
	if err != nil && err.Error() != "not found" {
		return nil, err
	}
	log.Println(product)
	if product.Id.Hex() != "" {
		return product, nil
	}
	item, err := parsers.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	product = new(Product)
	product.Name = item.Name
	product.ProductLink = item.ProductLink
	product.ProductImage = item.ProductImage
	product.MainCategory = item.MainCategory
	product.SubCategory = item.SubCategory
	product.PriceValue = item.PriceValue
	product.PriceCurrency = item.PriceCurrency
	price := new(Price)
	price.Date = time.Now().UTC()
	price.PriceValue = item.PriceValue
	price.PriceCurrency = item.PriceCurrency
	product.PriceHistory = []Price{}
	product.PriceHistory = append(product.PriceHistory, *price)
	product.Description = item.Description
	product.Specs = []Attribute{}
	for key, value := range item.Attributes {
		product.Specs = append(product.Specs, Attribute{Name: key, Value: value})
	}
	return product, nil
}
