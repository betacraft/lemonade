package models

import (
	"github.com/rainingclouds/lemonades/db"
	"testing"
)

func TestFetchProductInfo(t *testing.T) {
	db.InitMongo()
	_, err := FetchProductInfo("http://www.flipkart.com/vu-23-8-jl3-60-cm-23-8-led-tv/p/itme4zmg2jw5zx93?pid=TVSE4ZMF72WZWMMZ&srno=b_2&al=hB9J3RAl8wvhTvYOYdIwNwLSspt%2BtxYe5fPn1SCHzu%2FUIG3Rxz58p4aLq2lx4bRfFLwHQxVDMNU%3D&ref=2876720e-6c1f-4ff6-b3bb-da54f36e4de0")
	if err != nil {
		t.Error(err)
	}
}
