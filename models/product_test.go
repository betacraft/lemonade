package models

import (
	"github.com/rainingclouds/lemonades/db"
	"testing"
)

func TestFetchProductInfoFlipkart(t *testing.T) {
	db.InitMongo()
	_, err := FetchProductInfo("http://www.flipkart.com/dell-3440-latitude-intel-core-i5-35-56-cm-500-gb-hdd-4-ddr3-linux-ubuntu-notebook/p/itmdwzfzhtyh5ubb?pid=COMDWZFZMTH5WYXR&ref=L%3A4691174953981854897&srno=p_1&query=dell+latitude&otracker=from-search")
	if err != nil {
		t.Error(err)
	}
}

// func TestFetchProductInfoAmazon(t *testing.T) {
// 	db.InitMongo()
// 	_, err := FetchProductInfo("http://www.amazon.in/gp/product/B00T9N0Y9G/ref=s9_ri_gw_g147_i3?pf_rd_m=A1VBAL9TL5WCBF&pf_rd_s=center-5&pf_rd_r=1NRBCR365NMM42C5X34E&pf_rd_t=101&pf_rd_p=525630627&pf_rd_i=1320006031")
// 	if err != nil {
// 		t.Error(err)
// 	}
// }
