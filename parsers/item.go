package parsers

//go:generate easytags Item.go bson
//go:generate easytags Item.go json

type Item struct {
	ProductLink  string `json:"product_link" bson:"product_link"`
	ProductImage string `json:"product_image" bson:"product_image"`

	Name         string `json:"name" bson:"name"`
	MainCategory string `json:"main_category" bson:"main_category"`
	SubCategory  string `json:"sub_category" bson:"sub_category"`

	PriceValue    int64  `json:"price_value" bson:"price_value"`
	PriceCurrency string `json:"price_currency" bson:"price_currency"`

	Description string            `json:"description" bson:"description"`
	Attributes  map[string]string `json:"attributes" bson:"attributes"`
}
