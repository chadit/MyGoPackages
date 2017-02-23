package Mws

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// ProductTrackings array of ProductTrackings
type ProductTrackings []ProductTracking

// Len -lenth of sort
func (p ProductTrackings) Len() int {
	return len(p)
}

// Less - sort by seller FeedbackCount then total amount in an attempt to get the buy box winner
func (p ProductTrackings) Less(i, j int) bool {
	if p[i].SellerFeedbackCount > p[j].SellerFeedbackCount {
		return true
	} else if p[i].TotalAmount > p[j].TotalAmount {
		return true
	}
	return false
}

// Swap - lll
func (p ProductTrackings) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// ProductTracking model
type ProductTracking struct {
	ID                           string    `json:"id"  bson:"_id" binding:"required"`
	DateCreated                  time.Time `json:"dateCreated" bson:"DateCreated" binding:"required"`
	MarketplaceID                string    `json:"marketplaceId" bson:"MarketplaceId"`
	Asin                         string    `json:"asin" bson:"Asin"`
	Upc                          string    `json:"UPC" bson:"Upc"`
	Domain                       string    `json:"domain" bson:"Domain"`
	Title                        string    `json:"title" bson:"Title"`
	Category                     string    `json:"category" bson:"Category"`
	Author                       string    `json:"author" bson:"Author"`
	Condition                    string    `json:"condition" bson:"Condition"`
	SubCondition                 string    `json:"subCondition" bson:"SubCondition"`
	PathName                     string    `json:"pathname" bson:"PathName"`
	ImageURL                     string    `json:"imageUrl" bson:"ImageUrl"`
	CurrencyCode                 string    `json:"currencyCode" bson:"CurrencyCode"`
	RegularAmount                float64   `json:"regularPrice" bson:"RegularAmount"`
	SaleAmount                   float64   `json:"salePrice" bson:"SaleAmount"`
	ShippingAmount               float64   `json:"shippingPrice" bson:"ShippingAmount"`
	TotalAmount                  float64   `json:"-" bson:"TotalAmount"`
	SalesRank                    int       `json:"salesRank" bson:"SalesRank"`
	SellerFeedbackCount          int       `json:"sellerFeedbackCount" bson:"SellerFeedbackCount"`
	SellerPositiveFeedbackRating string    `json:"sellerPositiveFeedbackRating" bson:"SellerPositiveFeedbackRating"`
	Count                        int       `json:"count" bson:"Count"`
	Channel                      string    `json:"channel" bson:"Channel"`
	IsSoldByAmazon               bool      `json:"isSoldByAmazon" bson:"IsSoldByAmazon"`
	IsBuyBoxEligible             bool      `json:"isBuyBoxEligible" bson:"IsBuyBoxEligible"`
	PackageLength                float64   `json:"packageLength" bson:"PackageLength"`
	PackageWidth                 float64   `json:"packageWidth" bson:"PackageWidth"`
	PackageHeight                float64   `json:"packageHeight" bson:"PackageHeight"`
	PackageWeight                float64   `json:"packageWeight" bson:"PackageWeight"`
	AmazonFees                   float64   `json:"amazonFees" bson:"-"`
	RetryCount                   int       `json:"retryCount" bson:"-"`
	Message                      string    `json:"message" bson:"-"`
}

// NewProductTracking gets a new object
func NewProductTracking(asin string) ProductTracking {
	newItem := ProductTracking{}
	newItem.InitProductTracking(asin)
	return newItem
}

// InitProductTracking set defaults
func (p *ProductTracking) InitProductTracking(asin string) {
	p.Asin = asin
	p.CurrencyCode = "USD"
	p.SalesRank = 9223372036854775807 // max number
	p.Channel = "Merchant"

	eventTime := time.Now().UTC()
	p.ID = uuid.NewV4().String()
	p.DateCreated = eventTime
}

// SetupSaveProductTracking updates the user object modified uers
func (p *ProductTracking) SetupSaveProductTracking() {
	eventTime := time.Now().UTC()
	if p.ID == "" {
		p.ID = uuid.NewV4().String()
	}

	if p.DateCreated.IsZero() {
		p.DateCreated = eventTime
	}
}

// ConvertToAmazonResult converts the product tracking item to an amazon data object
func (p *ProductTracking) ConvertToAmazonResult() AmazonResult {
	return AmazonResult{
		TimeStamp:      p.DateCreated.Unix(),
		Domain:         p.Domain,
		Title:          p.Title,
		RegularAmount:  p.RegularAmount,
		SaleAmount:     p.SaleAmount,
		ShippingAmount: p.ShippingAmount,
		PathName:       p.PathName,
		Message:        "",
		RetryCount:     0,

		AmazonFees:                   p.AmazonFees,
		Count:                        p.Count,
		Condition:                    p.Condition,
		SubCondition:                 p.SubCondition,
		Channel:                      p.Channel,
		IsSoldByAmazon:               p.IsSoldByAmazon,
		IsBuyBoxEligible:             p.IsBuyBoxEligible,
		SellerFeedbackCount:          p.SellerFeedbackCount,
		SellerPositiveFeedbackRating: p.SellerPositiveFeedbackRating,
		SalesRank:                    p.SalesRank,
	}
}

// ConvertToDefaultResult converts the product tracking item to an default data object
func (p *ProductTracking) ConvertToDefaultResult() AmazonResult {
	return AmazonResult{
		TimeStamp:      p.DateCreated.Unix(),
		Domain:         p.Domain,
		Title:          p.Title,
		RegularAmount:  p.RegularAmount,
		SaleAmount:     p.SaleAmount,
		ShippingAmount: p.ShippingAmount,
		PathName:       p.PathName,
		Message:        "",
		RetryCount:     0,

		AmazonFees:     p.AmazonFees,
		Count:          p.Count,
		Condition:      p.Condition,
		SubCondition:   p.SubCondition,
		Channel:        p.Channel,
		IsSoldByAmazon: p.IsSoldByAmazon,
	}
}

// CheckAndValidateProductPrices will massage the prices, if sale price is zero it will make it equal to regular price
// if regular is zero it will set it to sale
// it will also check if shipping is included in sale price
func (p *ProductTracking) CheckAndValidateProductPrices() {
	if p.SaleAmount > 0 {
		// Amazon started to add Shipping to thier Landed amount.  We want to keep sale and shipping seperate
		// so if this happens, we will subtrack shipping
		if p.SaleAmount == p.RegularAmount+p.ShippingAmount {
			p.SaleAmount = p.SaleAmount - p.ShippingAmount
			p.TotalAmount = p.SaleAmount
		} else {
			p.TotalAmount = p.SaleAmount + p.ShippingAmount
		}
	}

	if p.RegularAmount == 0 && p.SaleAmount > 0 {
		p.RegularAmount = p.SaleAmount
	}

	if p.SaleAmount == 0 && p.RegularAmount > 0 {
		p.SaleAmount = p.RegularAmount
	}
}

// AmazonResult returns history formated as an amazon payload
type AmazonResult struct {
	TimeStamp                    int64   `json:"timestamp" bson:"-"`
	Domain                       string  `json:"domain" bson:"-"`
	Title                        string  `json:"title" bson:"-"`
	RegularAmount                float64 `json:"regularPrice" bson:"-"`
	SaleAmount                   float64 `json:"salePrice" bson:"-"`
	ShippingAmount               float64 `json:"shippingPrice" bson:"-"`
	PathName                     string  `json:"pathname" bson:"-"`
	Message                      string  `json:"message" bson:"-"`
	RetryCount                   int     `json:"retryCount" bson:"-"`
	AmazonFees                   float64 `json:"amazonFees" bson:"-"`
	Count                        int     `json:"count" bson:"-"`
	Condition                    string  `json:"condition" bson:"-"`
	SubCondition                 string  `json:"subCondition" bson:"-"`
	Channel                      string  `json:"channel" bson:"-"`
	IsSoldByAmazon               bool    `json:"isSoldByAmazon" bson:"-"`
	IsBuyBoxEligible             bool    `json:"isBuyBoxEligible" bson:"-"`
	SellerFeedbackCount          int     `json:"sellerFeedbackCount" bson:"-"`
	SellerPositiveFeedbackRating string  `json:"sellerPositiveFeedbackRating" bson:"-"`
	SalesRank                    int     `json:"salesRank" bson:"-"`
}

// DefaultResults returns history formated as an default payload
type DefaultResults struct {
	TimeStamp      int     `json:"timestamp" bson:"-"`
	Domain         string  `json:"domain" bson:"-"`
	Title          string  `json:"title" bson:"-"`
	RegularAmount  float64 `json:"regularPrice" bson:"-"`
	SaleAmount     float64 `json:"salePrice" bson:"-"`
	ShippingAmount float64 `json:"shippingPrice" bson:"-"`
	PathName       string  `json:"pathname" bson:"-"`
	Message        string  `json:"message" bson:"-"`
	RetryCount     int     `json:"retryCount" bson:"-"`
}

type lowestPricedOffersAttribute struct {
	MarketplaceID string `json:"-MarketplaceID"`
	ASIN          string `json:"-ASIN"`
	ItemCondition string `json:"-ItemCondition"`
}

type amazonFeeOption struct {
	ReferralFeesPercent  float64
	VcfDomesticStandard  float64
	VcfDomesticExpedited float64
	VcfInternational     *float64
	MediaFG              bool
	VcfPerPound          *float64
	ShippingCreditID     int
	MinReferralFees1     bool
	MinReferralFees2     bool
}

func getDefaultamazonFeeOptions() map[string]amazonFeeOption {
	amazonitems := make(map[string]amazonFeeOption)
	amazonitems["Amazon Kindle"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Automotive Parts and Accessories"] = amazonFeeOption{ReferralFeesPercent: 12, VcfDomesticStandard: 0.65, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Baby Products (excluding baby apparel)"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.65, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Books"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 1.35, VcfDomesticExpedited: 1.35, VcfInternational: createFloat64Pointer(1.35), MediaFG: true, VcfPerPound: nil, ShippingCreditID: 2, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Camera and Photo"] = amazonFeeOption{ReferralFeesPercent: 8, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Cell Phone Accessories"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Consumer Electronics"] = amazonFeeOption{ReferralFeesPercent: 8, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Home & Garden (including Pet Supplies)**"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Industrial & Scientific (including Food Service and Janitorial & Sanitation)"] = amazonFeeOption{ReferralFeesPercent: 12, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Kindle Accessories"] = amazonFeeOption{ReferralFeesPercent: 25, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Office Products"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Personal Computers"] = amazonFeeOption{ReferralFeesPercent: 6, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Software & Computer Games"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 1.35, VcfDomesticExpedited: 1.35, VcfInternational: nil, MediaFG: true, VcfPerPound: nil, ShippingCreditID: 1, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Sporting Goods*"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Tires & Wheels"] = amazonFeeOption{ReferralFeesPercent: 10, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Tools & Home Improvement**"] = amazonFeeOption{ReferralFeesPercent: 12, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Toys"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Video & DVD"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.80, VcfDomesticExpedited: 0.80, VcfInternational: nil, MediaFG: true, VcfPerPound: nil, ShippingCreditID: 5, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Video Game Consoles"] = amazonFeeOption{ReferralFeesPercent: 8, VcfDomesticStandard: 1.35, VcfDomesticExpedited: 1.35, VcfInternational: nil, MediaFG: true, VcfPerPound: nil, ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Watches"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.45, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Unlocked Cell Phones"] = amazonFeeOption{ReferralFeesPercent: 8, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Any Other Products"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Amazon Devices"] = amazonFeeOption{ReferralFeesPercent: 25, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Apparel"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Art"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Art and Craft Supply"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Audible"] = amazonFeeOption{ReferralFeesPercent: 25, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Automotive Al"] = amazonFeeOption{ReferralFeesPercent: 12, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Baby Product"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: nil, ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Beauty"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["BISS"] = amazonFeeOption{ReferralFeesPercent: 12, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Book"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Car Audio or Theater"] = amazonFeeOption{ReferralFeesPercent: 12, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["CE"] = amazonFeeOption{ReferralFeesPercent: 8, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Classical"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Collectibles"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Digital Accessories 5"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Digital Device Accessory"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Digital Music Album"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Digital Software"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Digital Video Games"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["DVD"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 1.35, VcfDomesticExpedited: 1.35, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["eBook"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Entertainment Memorabilia"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Furniture"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Gift Card"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Gourmet"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["GPS or Navigation System"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Grocery"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Health and Beauty"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Hobby"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Home"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Home Improvement"] = amazonFeeOption{ReferralFeesPercent: 12, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Home Theater"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Jewelry"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: true}
	amazonitems["Kitchen"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Lawn & Patio"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Lighting"] = amazonFeeOption{ReferralFeesPercent: 12, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Luggage"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Major Appliances"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["MotorCycle"] = amazonFeeOption{ReferralFeesPercent: 12, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Movie"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Music"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 1.35, VcfDomesticExpedited: 1.35, VcfInternational: createFloat64Pointer(1.35), MediaFG: true, VcfPerPound: nil, ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Musical Instruments"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Network Media Player"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Office Electronics"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Office Product"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Pantry"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["PC Accessory"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Personal Computer"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Pet Products"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Photography"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Portable Audio Video"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Prescription Drugs"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Prestige Beauty"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Receiver or Amplifier"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Shoes"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Single Detail Page Misc"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Software"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Speakers"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Sports"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Television"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Toy"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["TV Series Episode Video on Demand"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.54, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Video"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Video Games"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 1.35, VcfDomesticExpedited: 1.35, VcfInternational: nil, MediaFG: true, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 1, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Watch"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: true}
	amazonitems["Wireless"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: true, MinReferralFees2: false}
	amazonitems["Wireless Phone"] = amazonFeeOption{ReferralFeesPercent: 8, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Wireless Phone Accessory"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}
	amazonitems["Wireless Plan Option"] = amazonFeeOption{ReferralFeesPercent: 15, VcfDomesticStandard: 0.45, VcfDomesticExpedited: 0.65, VcfInternational: nil, MediaFG: false, VcfPerPound: createFloat64Pointer(0.05), ShippingCreditID: 6, MinReferralFees1: false, MinReferralFees2: false}

	return amazonitems
}

func createFloat64Pointer(value float64) *float64 {
	newValue := value
	return &newValue
}

func getAmazonProductTypeNames() []string {
	return []string{
		"ACCESSORY",
		"ACCESSORY_OR_PART_OR_SUPPLY",
		"AMAZON_TABLET_ACCESSORY",
		"ANTENNA",
		"AUTO_ACCESSORY",
		"BATTERY",
		"CABLE_OR_ADAPTER",
		"CAMERA_OTHER_ACCESSORIES",
		"CAMERA_POWER_SUPPLY",
		"CE_ACCESSORY",
		"CE_CARRYING_CASE_OR_BAG",
		"COMPUTER_COOLING_DEVICE",
		"COMPUTER_INPUT_DEVICE",
		"COMPUTER_OUTPUT_DEVICE_ACCESSORY",
		"COMPUTER_SPEAKER",
		"COMPUTER_VIDEO_GAME_CONTOLLER",
		"CONSUMER_ELECTRONICS_PARTS",
		"DIGITAL_DEVICE_ACCESSORY",
		"ELECTRONIC_COMPONENT",
		"GPS_OR_NAVIGATION_ACCESSORY",
		"HOME_LIGHTING_ACCESSORY",
		"INSTRUMENT_PARTS_AND_ACCESSORIES",
		"MOTHERBOARD",
		"PHONE_ACCESSORY",
		"PORTABLE_AV_DEVICE",
		"SOUND_CARD",
		"SPEAKERS",
		"SYSTEM_POWER_DEVICE",
		"VIDEO_GAME_ACCESSORIES",
		"WIRELESS_ACCESSORY",
		"AUDIO_OR_VIDEO",
		"COMPUTER_COMPONENT",
		"COMPUTER_ADD_ON",
		"HEADPHONES",
		"RAM_MEMORY",
		"REMOTE_CONTROL",
	}
}
