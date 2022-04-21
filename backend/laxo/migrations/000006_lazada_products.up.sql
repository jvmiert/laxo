BEGIN;
CREATE TABLE IF NOT EXISTS products_lazada(
  id CHAR(26) DEFAULT ulid_create() NOT NULL PRIMARY KEY,
  lazada_id BIGINT NOT NULL,
  lazada_primary_category BIGINT NOT NULL,
  created TIMESTAMP WITH TIME ZONE NOT NULL,
  updated TIMESTAMP WITH TIME ZONE NOT NULL,
  status VARCHAR(128) NOT NULL,
  sub_status VARCHAR(128) NOT NULL,
  shop_id CHAR(26) NOT NULL,
  CONSTRAINT fk_product_shop_lazada FOREIGN KEY(shop_id) REFERENCES shops(id)
);

CREATE TABLE IF NOT EXISTS products_attribute_lazada(
  id CHAR(26) DEFAULT ulid_create() NOT NULL PRIMARY KEY,
  name TEXT,
  short_description TEXT,
  description TEXT,
  brand TEXT,
  model TEXT,
  headphone_features TEXT,
  bluetooth TEXT,
  warranty_type TEXT,
  warranty TEXT,
  hazmat TEXT,
  expire_date TEXT,
  brand_classification TEXT,
  ingredient_preference TEXT,
  lot_number TEXT,
  units_hb TEXT,
  fmlt_skincare TEXT,
  quantitative TEXT,
  skincare_by_age TEXT,
  skin_benefit TEXT,
  skin_type TEXT,
  user_manual TEXT,
  country_origin_hb TEXT,
  color_family TEXT,
  fragrance_family TEXT,
  source TEXT,
  product_id CHAR(26) NOT NULL,
  CONSTRAINT fk_attribute_products_lazada FOREIGN KEY(product_id) REFERENCES products_lazada(id)
);

CREATE TABLE IF NOT EXISTS products_sku_lazada(
  id CHAR(26) DEFAULT ulid_create() NOT NULL PRIMARY KEY,
  status VARCHAR(128),
  quantity INTEGER,
  seller_sku VARCHAR(128),
  shop_sku VARCHAR(128),
  url TEXT,
  color_family VARCHAR(64),
  price INTEGER,
  available INTEGER,
  sku_id BIGINT,
  package_content TEXT,
  package_width VARCHAR(64),
  package_weight VARCHAR(64),
  package_length VARCHAR(64),
  package_height VARCHAR(64),
  special_price VARCHAR(64),
  special_to_time TIMESTAMP WITH TIME ZONE,
  special_from_time TIMESTAMP WITH TIME ZONE,
  special_from_date TIMESTAMP WITH TIME ZONE,
  special_to_date TIMESTAMP WITH TIME ZONE,
  product_id CHAR(26) NOT NULL,
  CONSTRAINT fk_sku_products_lazada FOREIGN KEY(product_id) REFERENCES products_lazada(id)
);
COMMIT;
