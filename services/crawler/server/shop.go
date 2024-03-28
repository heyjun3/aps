package server

import (
	"context"

	"connectrpc.com/connect"

	shopv1 "crawler/server/gen/shop/v1"
	"crawler/shop"
)

func convertShopsIntoShopv1Shops(shops []*shop.Shop) *shopv1.Shops {
	shopsv1 := []*shopv1.Shop{}
	for _, shop := range shops {
		shopv1 := shopv1.Shop{Id: shop.ID, SiteName: shop.SiteName, Name: shop.Name, Url: shop.URL, Interval: shop.Interval}
		shopsv1 = append(shopsv1, &shopv1)
	}
	return &shopv1.Shops{Shop: shopsv1}
}

func convertShopv1ShopsIntoShops(shops []*shopv1.Shop) []*shop.Shop {
	rakutenShops := []*shop.Shop{}
	for _, s := range shops {
		rakutenShop := &shop.Shop{ID: s.Id, SiteName: s.SiteName, Name: s.Name, URL: s.Url, Interval: s.Interval}
		rakutenShops = append(rakutenShops, rakutenShop)
	}
	return rakutenShops
}

type ShopServer struct{}

func (s *ShopServer) ShopList(ctx context.Context, req *connect.Request[shopv1.ShopListRequest]) (*connect.Response[shopv1.ShopListResponse], error) {
	repo := shop.ShopRepository{}
	shops, err := repo.GetAll(db, context.Background())
	if err != nil {
		return nil, err
	}
	shopsv1 := convertShopsIntoShopv1Shops(shops)
	res := connect.NewResponse(&shopv1.ShopListResponse{
		Shops: shopsv1,
	})
	return res, nil
}

func (s *ShopServer) CreateShop(ctx context.Context, req *connect.Request[shopv1.CreateShopRequest]) (*connect.Response[shopv1.CreateShopResponse], error) {
	shops := convertShopv1ShopsIntoShops(req.Msg.Shops.GetShop())
	repo := shop.ShopRepository{}
	err := repo.Save(db, ctx, shops)
	if err != nil {
		return nil, err
	}
	res := connect.NewResponse(&shopv1.CreateShopResponse{
		Shops: req.Msg.Shops,
	})
	return res, nil
}

func (s *ShopServer) DeleteShop(ctx context.Context, req *connect.Request[shopv1.DeleteShopRequest]) (*connect.Response[shopv1.DeleteShopResponse], error) {
	shops := []*shop.Shop{}
	for _, id := range req.Msg.Ids {
		shops = append(shops, &shop.Shop{ID: id})
	}
	repo := shop.ShopRepository{}
	err := repo.DeleteShops(context.Background(), db, shops)
	if err != nil {
		return nil, err
	}
	res := connect.NewResponse(&shopv1.DeleteShopResponse{
		Shops: convertShopsIntoShopv1Shops(shops),
	})

	return res, nil
}
