package parsers

import (
	"log"
	"testing"
)

// func TestSnapdeal(t *testing.T) {
// 	item, err := Parse("http://www.snapdeal.com/product/infocus-m330/682019088263?")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	log.Println(item)
// }

func TestPaytm(t *testing.T) {
	item, err := Parse("https://paytm.com/shop/p/wd-elements-2-5-inch-1-tb-external-hard-disk-black-cmplxwd_elements_1tb_black_168")
	if err != nil {
		t.Error(err)
	}
	log.Println(item)
}
