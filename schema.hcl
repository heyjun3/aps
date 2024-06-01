table "products" {
  schema = schema.crawler
  column "site_code" {
    null = false
    type = character_varying
  }
  column "shop_code" {
    null = false
    type = character_varying
  }
  column "product_code" {
    null = false
    type = character_varying
  }
  column "name" {
    null = true
    type = character_varying
  }
  column "jan" {
    null = true
    type = character_varying
  }
  column "price" {
    null = true
    type = bigint
  }
  column "url" {
    null = true
    type = character_varying
  }
  primary_key {
    columns = [column.site_code, column.shop_code, column.product_code]
  }
  partition {
    type    = LIST
    columns = [column.site_code]
  }
}
table "crawler" "schema_migrations" {
  schema = schema.crawler
  column "version" {
    null = false
    type = bigint
  }
  column "dirty" {
    null = false
    type = boolean
  }
  primary_key {
    columns = [column.version]
  }
}
table "ark_products" {
  schema = schema.public
  column "name" {
    null = true
    type = character_varying
  }
  column "jan" {
    null = true
    type = character_varying
  }
  column "price" {
    null = true
    type = bigint
  }
  column "shop_code" {
    null = false
    type = character_varying
  }
  column "product_code" {
    null = false
    type = character_varying
  }
  column "url" {
    null = true
    type = character_varying
  }
  primary_key {
    columns = [column.shop_code, column.product_code]
  }
  index "ark_products_product_code_idx" {
    columns = [column.product_code]
  }
}
table "asins_info" {
  schema = schema.public
  column "asin" {
    null = false
    type = character_varying
  }
  column "jan" {
    null = true
    type = character_varying
  }
  column "title" {
    null = true
    type = character_varying
  }
  column "quantity" {
    null = true
    type = bigint
  }
  column "modified" {
    null = true
    type = date
  }
  primary_key {
    columns = [column.asin]
  }
}
table "bomber_products" {
  schema = schema.public
  column "name" {
    null = true
    type = character_varying
  }
  column "jan" {
    null = true
    type = character_varying
  }
  column "price" {
    null = true
    type = bigint
  }
  column "shop_code" {
    null = false
    type = character_varying
  }
  column "product_code" {
    null = false
    type = character_varying
  }
  column "url" {
    null = true
    type = character_varying
  }
  primary_key {
    columns = [column.shop_code, column.product_code]
  }
}
table "buffalo_products" {
  schema = schema.public
  column "name" {
    null = true
    type = character_varying
  }
  column "jan" {
    null = true
    type = character_varying
  }
  column "price" {
    null = true
    type = bigint
  }
  column "shop_code" {
    null = false
    type = character_varying
  }
  column "product_code" {
    null = false
    type = character_varying
  }
  column "url" {
    null = true
    type = character_varying
  }
  primary_key {
    columns = [column.shop_code, column.product_code]
  }
}
table "current_prices" {
  schema = schema.public
  column "seller_sku" {
    null = false
    type = character_varying
  }
  column "amount" {
    null = false
    type = bigint
  }
  column "point" {
    null = false
    type = bigint
  }
  column "created_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "percent_point" {
    null = false
    type = bigint
  }
  primary_key {
    columns = [column.seller_sku]
  }
  foreign_key "current_prices_seller_sku_fkey" {
    columns     = [column.seller_sku]
    ref_columns = [table.inventories.column.seller_sku]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}
table "desired_prices" {
  schema = schema.public
  column "seller_sku" {
    null = false
    type = character_varying
  }
  column "amount" {
    null = false
    type = bigint
  }
  column "point" {
    null = false
    type = bigint
  }
  column "percent_point" {
    null = false
    type = bigint
  }
  column "created_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.seller_sku]
  }
  foreign_key "desired_prices_seller_sku_fkey" {
    columns     = [column.seller_sku]
    ref_columns = [table.inventories.column.seller_sku]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}
table "favoriteproduct" {
  schema = schema.public
  column "url" {
    null = false
    type = character_varying
  }
  column "jan" {
    null = false
    type = character_varying
  }
  column "cost" {
    null = true
    type = integer
  }
  primary_key {
    columns = [column.url, column.jan]
  }
}
table "hikaritv_products" {
  schema = schema.public
  column "name" {
    null = true
    type = character_varying
  }
  column "jan" {
    null = true
    type = character_varying
  }
  column "price" {
    null = true
    type = bigint
  }
  column "shop_code" {
    null = false
    type = character_varying
  }
  column "product_code" {
    null = false
    type = character_varying
  }
  column "url" {
    null = true
    type = character_varying
  }
  primary_key {
    columns = [column.shop_code, column.product_code]
  }
}
table "ikebe_product" {
  schema = schema.public
  column "name" {
    null = true
    type = character_varying
  }
  column "jan" {
    null = true
    type = character_varying
  }
  column "price" {
    null = true
    type = bigint
  }
  column "shop_code" {
    null = false
    type = character_varying
  }
  column "product_code" {
    null = false
    type = character_varying
  }
  column "url" {
    null = true
    type = character_varying
  }
  primary_key {
    columns = [column.shop_code, column.product_code]
  }
}
table "inactivestock" {
  schema = schema.public
  column "SKU" {
    null = false
    type = character_varying
  }
  column "asin" {
    null = true
    type = character_varying
  }
  primary_key {
    columns = [column.SKU]
  }
}
table "inventories" {
  schema = schema.public
  column "asin" {
    null = false
    type = character_varying
  }
  column "fnsku" {
    null = false
    type = character_varying
  }
  column "seller_sku" {
    null = false
    type = character_varying
  }
  column "condition" {
    null = false
    type = character_varying
  }
  column "product_name" {
    null = false
    type = character_varying
  }
  column "quantity" {
    null = false
    type = integer
  }
  column "created_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "fulfillable_quantity" {
    null = false
    type = bigint
  }
  primary_key {
    columns = [column.seller_sku]
  }
}
table "kaago_products" {
  schema = schema.public
  column "name" {
    null = true
    type = character_varying
  }
  column "jan" {
    null = true
    type = character_varying
  }
  column "price" {
    null = true
    type = bigint
  }
  column "shop_code" {
    null = false
    type = character_varying
  }
  column "product_code" {
    null = false
    type = character_varying
  }
  column "url" {
    null = true
    type = character_varying
  }
  primary_key {
    columns = [column.shop_code, column.product_code]
  }
}
table "keepa_products" {
  schema = schema.public
  column "asin" {
    null = false
    type = character_varying
  }
  column "sales_drops_90" {
    null = true
    type = integer
  }
  column "created" {
    null = true
    type = date
  }
  column "modified" {
    null = true
    type = date
  }
  column "price_data" {
    null = true
    type = jsonb
  }
  column "rank_data" {
    null = true
    type = jsonb
  }
  column "render_data" {
    null = true
    type = jsonb
  }
  column "drops_ma_7" {
    null = true
    type = integer
  }
  primary_key {
    columns = [column.asin]
  }
  index "ix_keepaproducts_asin" {
    columns = [column.asin]
  }
}
table "lowest_prices" {
  schema = schema.public
  column "seller_sku" {
    null = false
    type = character_varying
  }
  column "amount" {
    null = false
    type = bigint
  }
  column "point" {
    null = false
    type = bigint
  }
  column "created_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "percent_point" {
    null = false
    type = bigint
  }
  primary_key {
    columns = [column.seller_sku]
  }
  foreign_key "lowest_prices_seller_sku_fkey" {
    columns     = [column.seller_sku]
    ref_columns = [table.inventories.column.seller_sku]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}
table "murauchi_products" {
  schema = schema.public
  column "name" {
    null = true
    type = character_varying
  }
  column "jan" {
    null = true
    type = character_varying
  }
  column "price" {
    null = true
    type = bigint
  }
  column "shop_code" {
    null = false
    type = character_varying
  }
  column "product_code" {
    null = false
    type = character_varying
  }
  column "url" {
    null = true
    type = character_varying
  }
  primary_key {
    columns = [column.shop_code, column.product_code]
  }
}
table "mws_products" {
  schema = schema.public
  column "asin" {
    null = false
    type = character_varying
  }
  column "filename" {
    null = false
    type = character_varying
  }
  column "title" {
    null = true
    type = character_varying
  }
  column "jan" {
    null = true
    type = character_varying
  }
  column "unit" {
    null = true
    type = bigint
  }
  column "price" {
    null = true
    type = bigint
  }
  column "cost" {
    null = true
    type = bigint
  }
  column "fee_rate" {
    null = true
    type = double_precision
  }
  column "shipping_fee" {
    null = true
    type = bigint
  }
  column "profit" {
    null = true
    type = bigint
    as {
      expr = "((((price - (cost * unit)))::double precision - (((price)::double precision * fee_rate) * (1.1)::double precision)) - (shipping_fee)::double precision)"
      type = STORED
    }
  }
  column "profit_rate" {
    null = true
    type = double_precision
    as {
      expr = "(((((price - (cost * unit)))::double precision - (((price)::double precision * fee_rate) * (1.1)::double precision)) - (shipping_fee)::double precision) / (price)::double precision)"
      type = STORED
    }
  }
  column "created_at" {
    null = true
    type = timestamp
  }
  column "url" {
    null = true
    type = character_varying
  }
  primary_key {
    columns = [column.asin, column.filename]
  }
  index "is_filename" {
    columns = [column.filename]
  }
  index "ix_asin" {
    columns = [column.asin]
  }
  index "ix_profit" {
    columns = [column.profit]
  }
}
table "netsea_products" {
  schema = schema.public
  column "name" {
    null = true
    type = character_varying
  }
  column "jan" {
    null = true
    type = character_varying
  }
  column "price" {
    null = true
    type = bigint
  }
  column "shop_code" {
    null = false
    type = character_varying
  }
  column "product_code" {
    null = false
    type = character_varying
  }
  column "url" {
    null = true
    type = character_varying
  }
  primary_key {
    columns = [column.shop_code, column.product_code]
  }
}
table "netsea_shops" {
  schema = schema.public
  column "name" {
    null = true
    type = character_varying
  }
  column "shop_id" {
    null = false
    type = character_varying
  }
  primary_key {
    columns = [column.shop_id]
  }
}
table "nojima_products" {
  schema = schema.public
  column "name" {
    null = true
    type = character_varying
  }
  column "jan" {
    null = true
    type = character_varying
  }
  column "price" {
    null = true
    type = bigint
  }
  column "shop_code" {
    null = false
    type = character_varying
  }
  column "product_code" {
    null = false
    type = character_varying
  }
  column "url" {
    null = true
    type = character_varying
  }
  primary_key {
    columns = [column.shop_code, column.product_code]
  }
}
table "pc4u_products" {
  schema = schema.public
  column "name" {
    null = true
    type = character_varying
  }
  column "jan" {
    null = true
    type = character_varying
  }
  column "price" {
    null = true
    type = bigint
  }
  column "shop_code" {
    null = false
    type = character_varying
  }
  column "product_code" {
    null = false
    type = character_varying
  }
  column "url" {
    null = true
    type = character_varying
  }
  primary_key {
    columns = [column.shop_code, column.product_code]
  }
  index "ix_pc4uproducts_product_code" {
    columns = [column.product_code]
  }
}
table "product_master" {
  schema = schema.public
  column "date" {
    null = true
    type = integer
  }
  column "name" {
    null = true
    type = character_varying
  }
  column "asin" {
    null = true
    type = character_varying
  }
  column "jan" {
    null = true
    type = character_varying
  }
  column "sku" {
    null = false
    type = character_varying
  }
  column "fnsku" {
    null = true
    type = character_varying
  }
  column "danger_class" {
    null = true
    type = character_varying
  }
  column "sell_price" {
    null = true
    type = integer
  }
  column "cost_price" {
    null = true
    type = integer
  }
  primary_key {
    columns = [column.sku]
  }
}
table "rakuten_products" {
  schema = schema.public
  column "name" {
    null = true
    type = character_varying
  }
  column "jan" {
    null = true
    type = character_varying
  }
  column "price" {
    null = true
    type = bigint
  }
  column "shop_code" {
    null = false
    type = character_varying
  }
  column "product_code" {
    null = false
    type = character_varying
  }
  column "url" {
    null = true
    type = character_varying
  }
  primary_key {
    columns = [column.shop_code, column.product_code]
  }
}
table "run_service_histories" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "shop_name" {
    null = false
    type = character_varying
  }
  column "url" {
    null = false
    type = character_varying
  }
  column "status" {
    null = false
    type = character_varying
  }
  column "started_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "ended_at" {
    null = true
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
}
table "public" "schema_migrations" {
  schema = schema.public
  column "version" {
    null = false
    type = bigint
  }
  column "dirty" {
    null = false
    type = boolean
  }
  primary_key {
    columns = [column.version]
  }
}
table "shops" {
  schema = schema.public
  column "id" {
    null = false
    type = character_varying
  }
  column "site_name" {
    null = true
    type = character_varying
  }
  column "name" {
    null = true
    type = character_varying
  }
  column "url" {
    null = true
    type = character_varying
  }
  column "interval" {
    null = true
    type = character_varying
  }
  primary_key {
    columns = [column.id]
  }
}
table "spapi_fees" {
  schema = schema.public
  column "asin" {
    null = false
    type = character_varying
  }
  column "fee_rate" {
    null = true
    type = double_precision
  }
  column "ship_fee" {
    null = true
    type = bigint
  }
  column "modified" {
    null = true
    type = date
  }
  primary_key {
    columns = [column.asin]
  }
  foreign_key "spapi_fees_asin_fkey" {
    columns     = [column.asin]
    ref_columns = [table.asins_info.column.asin]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
}
table "spapi_prices" {
  schema = schema.public
  column "asin" {
    null = false
    type = character_varying
  }
  column "price" {
    null = true
    type = bigint
  }
  column "modified" {
    null = true
    type = date
  }
  primary_key {
    columns = [column.asin]
  }
  foreign_key "spapi_prices_asin_fkey" {
    columns     = [column.asin]
    ref_columns = [table.asins_info.column.asin]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
}
table "stock" {
  schema = schema.public
  column "sku" {
    null = false
    type = character_varying
  }
  column "home_stock_count" {
    null = true
    type = integer
  }
  column "fba_stock_count" {
    null = true
    type = integer
  }
  primary_key {
    columns = [column.sku]
  }
}
table "super_product_details" {
  schema = schema.public
  column "product_code" {
    null = false
    type = character_varying
  }
  column "set_number" {
    null = false
    type = integer
  }
  column "shop_code" {
    null = true
    type = character_varying
  }
  column "price" {
    null = true
    type = bigint
  }
  column "jan" {
    null = true
    type = character_varying
  }
  primary_key {
    columns = [column.product_code, column.set_number]
  }
  foreign_key "super_product_details_product_code_fkey" {
    columns     = [column.product_code]
    ref_columns = [table.super_products.column.product_code]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
}
table "super_products" {
  schema = schema.public
  column "name" {
    null = true
    type = character_varying
  }
  column "product_code" {
    null = false
    type = character_varying
  }
  column "price" {
    null = true
    type = bigint
  }
  column "shop_code" {
    null = true
    type = character_varying
  }
  column "url" {
    null = true
    type = character_varying
  }
  primary_key {
    columns = [column.product_code]
  }
}
table "super_shops" {
  schema = schema.public
  column "name" {
    null = true
    type = character_varying
  }
  column "shop_id" {
    null = false
    type = character_varying
  }
  primary_key {
    columns = [column.shop_id]
  }
}
schema "crawler" {
}
schema "public" {
  comment = "standard public schema"
}
