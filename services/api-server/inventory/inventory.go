package inventory

import "github.com/uptrace/bun"

type Inventory struct {
	bun.BaseModel   `bun:"table:inventries"`
	Asin            string `json:"asin" bun:"asin"`
	FnSku           string `json:"fnSku" bun:"fnsku"`
	SellerSku       string `json:"sellerSku" bun:"seller_sku"`
	Condition       string `json:"condition" bun:"condition"`
	LastUpdatedTime string `json:"lastUpdatedTime"`
	ProductName     string `json:"productName" bun:"product_name"`
	TotalQuantity   int    `json:"totalQuantity" bun:"quantity"`
	Price           *int   `bun:"price"`
	Point           *int   `bun:"point"`
	LowestPrice     *int   `bun:"lowest_price"`
	LowestPoint     *int   `bun:"lowest_point"`
}
