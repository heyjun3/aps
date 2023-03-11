package scrape

import (
	"context"
	"database/sql"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

func NewDBconnection(dsn string) (*sql.DB, error) {
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error("open error", err)
		return nil, err
	}
	return conn, nil
}

type Product interface {
	GenerateMessage(filename string) ([]byte, error)
	GetProductCode() string
	GetJan() string
	GetURL() string
	IsValidJan() bool
	SetJan(string)
	Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns boil.Columns, insertColumns boil.Columns) error
}

type Products []Product

func (p Products) getProductCodes() []string {
	var codes []string
	for _, pro := range p {
		codes = append(codes, pro.GetProductCode())
	}
	return codes
}

func (p Products) MapProducts(products Products) Products{
	mapped := map[string]Product{}
	for _, v := range products {
		code := v.GetProductCode()
		mapped[code] = v
	}

	for _, v := range p {
		product, exist := mapped[v.GetProductCode()]
		if !exist {
			continue
		}
		v.SetJan((product).GetJan())
	}
	return p
}
