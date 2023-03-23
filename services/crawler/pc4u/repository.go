package pc4u

import (
	"context"

	"github.com/uptrace/bun"

	"crawler/scrape"
)

type Pc4uProductRepository struct{
	scrape.ProductRepository
}

func (r Pc4uProductRepository) GetByProductCodes(conn *bun.DB,
	ctx context.Context, codes ...string) (scrape.Products, error) {

	var pc4uProducts []Pc4uProduct
	err := conn.NewSelect().
		Model(&pc4uProducts).
		Where("product_code IN (?)", bun.In(codes)).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	var products scrape.Products
	for i := 0; i < len(pc4uProducts); i++ {
		products = append(products, &pc4uProducts[i])
	}
	return products, nil
}

func getProducts() (func(*bun.DB, context.Context, ...string)(scrape.Products, error)){
	var pc4uProducts []Pc4uProduct
	return func(conn *bun.DB, ctx context.Context, codes ...string) (scrape.Products, error){
		err := conn.NewSelect().
			Model(&pc4uProducts).
			Where("product_code IN (?)", bun.In(codes)).
			Scan(ctx)
		if err != nil {
			return nil, err
		}
	
		var products scrape.Products
		for i := 0; i < len(pc4uProducts); i++ {
			products = append(products, &pc4uProducts[i])
		}
		return products, nil
	}
}