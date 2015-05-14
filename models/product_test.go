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
