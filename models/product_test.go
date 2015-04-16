package models

import (
	"github.com/rainingclouds/lemonades/db"
	"testing"
)

// func TestFetchProductInfoFlipkart(t *testing.T) {
// 	db.InitMongo()
// 	_, err := FetchProductInfo("http://www.flipkart.com/vu-23-8-jl3-60-cm-23-8-led-tv/p/itme4zmg2jw5zx93?pid=TVSE4ZMF72WZWMMZ&srno=b_2&al=hB9J3RAl8wvhTvYOYdIwNwLSspt%2BtxYe5fPn1SCHzu%2FUIG3Rxz58p4aLq2lx4bRfFLwHQxVDMNU%3D&ref=2876720e-6c1f-4ff6-b3bb-da54f36e4de0")
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

func TestFetchProductInfoAmazon(t *testing.T) {
	db.InitMongo()
	_, err := FetchProductInfo("http://www.amazon.in/gp/product/B00T9N0Y9G/ref=s9_ri_gw_g147_i3?pf_rd_m=A1VBAL9TL5WCBF&pf_rd_s=center-5&pf_rd_r=1NRBCR365NMM42C5X34E&pf_rd_t=101&pf_rd_p=525630627&pf_rd_i=1320006031")
	if err != nil {
		t.Error(err)
	}
}
