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

func TestFlipkart(t *testing.T) {
	item, err := Parse("http://www.flipkart.com/htc-one-max/p/itmdqrpcg9tzkrs7")
	if err != nil {
		t.Error(err)
	}
	log.Println(item)
}

// func TestSnapdeal(t *testing.T) {
// 	item, err := Parse("http://www.snapdeal.com/product/infocus-m330/682019088263?")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	log.Println(item)
// }

// func TestPaytm(t *testing.T) {
// 	item, err := Parse("https://paytm.com/shop/p/flow-bluetooth-speakers-for-4-1-home-theater-system-computer-tv-usb-mmc-remote-fm-CMPLXFLOW_FLBT41BL_HOMETHEATERSPKR_BLACK")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	log.Println(item)
// }
