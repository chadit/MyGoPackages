package Mws

import (
	"fmt"
	"sort"
	"testing"
)

func TestGetDefaultAmazonFeeCount(t *testing.T) {
	expected := 91
	defaultAmazonItems := getDefaultamazonFeeOptions()
	actual := len(defaultAmazonItems)
	if actual != expected {
		// t.Errorf("Test failed, expected: '%s', got:  '%s'", expectedString, actualString)
		t.Errorf("expected: '%d', got:  '%d'", expected, actual)
	}
}

func TestToyAmazonFeeOptions(t *testing.T) {
	expected := float64(15)
	defaultAmazonItems := getDefaultamazonFeeOptions()
	actual := defaultAmazonItems["Toy"]
	if actual.ReferralFeesPercent != expected {
		t.Errorf("expected: '%g', got:  '%g'", expected, actual.ReferralFeesPercent)
	}
}

func Test_ProductTracking_Sort(t *testing.T) {
	var items ProductTrackings
	item1 := ProductTracking{RegularAmount: 20.95, SaleAmount: 20.95, ShippingAmount: 14, TotalAmount: 34.95, SalesRank: 4594, SellerFeedbackCount: 22}
	item2 := ProductTracking{RegularAmount: 34.99, SaleAmount: 34.99, ShippingAmount: 0, TotalAmount: 34.99, SalesRank: 4594, SellerFeedbackCount: 27385}
	item3 := ProductTracking{RegularAmount: 24.99, SaleAmount: 24.99, ShippingAmount: 10.59, TotalAmount: 35.58, SalesRank: 4594, SellerFeedbackCount: 10224}
	item4 := ProductTracking{RegularAmount: 34.06, SaleAmount: 34.06, ShippingAmount: 5.95, TotalAmount: 40.01, SalesRank: 4594, SellerFeedbackCount: 83}

	items = append(items, item4)
	items = append(items, item3)
	items = append(items, item2)
	items = append(items, item1)

	sort.Sort(items)

	for i := range items {
		fmt.Println(items[i])
	}
	if items[0].TotalAmount != 34.99 {
		t.Errorf("expected: '%g', got:  '%g'", 34.99, items[0].TotalAmount)
	}

}
