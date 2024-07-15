package repeate

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type RepeateItem struct {
	bun.BaseModel `bun:"crawler.repeate_items"`
	ID            uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	URL           string    `bun:"url,notnull"`
	Jan           *string   `bun:"jan"`
}
