package parsers

import (
	"log"
	"testing"
)

// func TestAmazonIn(t *testing.T) {
// 	item, err := Parse("http://www.amazon.in/Nokia-105-Black-Color/dp/B00CZ50ZQW/ref=sr_1_1?s=electronics&ie=UTF8&qid=1432050591&sr=1-1")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	log.Println(item)
// }

// func TestFlipkart(t *testing.T) {
// 	item, err := Parse("http://www.flipkart.com/spice-smart-pulse-m-9010-smartwatch/p/itmeyfaq7gjgq4kf?pid=SMWEYF9SAZKGSG93&icmpid=reco_pp_personalhistoryFooter_storage_na_1")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	log.Println(item)
// }

func TestSnapdeal(t *testing.T) {
	item, err := Parse("http://www.snapdeal.com/product/infocus-m330/682019088263?")
	if err != nil {
		t.Error(err)
	}
	log.Println(item)
}

// func TestPaytm(t *testing.T) {
// 	item, err := Parse("https://paytm.com/shop/p/wd-elements-2-5-inch-1-tb-external-hard-disk-black-cmplxwd_elements_1tb_black_168")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	log.Println(item)
// }
