package server

import (
	"context"

	"connectrpc.com/connect"

	"crawler/rakuten"
	shopv1 "crawler/server/gen/shop/v1"
)

func convertShopsIntoShopv1Shops(shops []rakuten.Shop) *shopv1.Shops {
	shopsv1 := []*shopv1.Shop{}
	for _, shop := range shops {
		shopv1 := shopv1.Shop{Id: shop.ID, SiteName: shop.SiteName, Name: shop.Name, Url: shop.URL, Interval: shop.Interval}
		shopsv1 = append(shopsv1, &shopv1)
	}
	return &shopv1.Shops{Shop: shopsv1}
}

type ShopServer struct{}

func (s *ShopServer) ShopList(ctx context.Context, req *connect.Request[shopv1.ShopListRequest]) (*connect.Response[shopv1.ShopListResponse], error) {
	logger.Info("ShopList", "status", "run")
	repo := rakuten.ShopRepository{}
	shops, err := repo.GetAll(db, context.Background())
	if err != nil {
		return nil, err
	}
	shopsv1 := convertShopsIntoShopv1Shops(shops)
	res := connect.NewResponse(&shopv1.ShopListResponse{
		Shops: shopsv1,
	})
	logger.Info("ShopList", "status", "done")
	return res, nil
}
