BEGIN;
CREATE TABLE IF NOT EXISTS social_media_platforms(
  id INT NOT NULL PRIMARY KEY,
  platform_name TEXT
);

INSERT INTO social_media_platforms (id, platform_name)
VALUES
  (1, 'tiktok'),
  (2, 'facebook'),
  (3, 'zalo'),
  (4, 'instagram'),
  (5, 'telegram'),
  (6, 'signal'),
  (7, 'snapchat'),
  (8, 'youtube'),
  (9, 'whatsapp')
;

CREATE TABLE IF NOT EXISTS orders(
  id CHAR(26) DEFAULT ulid_to_string(ulid_generate()) NOT NULL PRIMARY KEY,
  selling_price NUMERIC(21,2),
  cost_price NUMERIC(21,2),
  discount_percent NUMERIC(19,4),
  discount_sum NUMERIC(19,4),
  discount_reason TEXT,
  shop_id CHAR(26) NOT NULL,
  created TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

  --@TODO: figure this out
  sales_platform VARCHAR(300),
  sales_social_media_id INT,
  CONSTRAINT order_platform_name CHECK(sales_platform IN ('lazada', 'tiki', 'shopee')),
  CONSTRAINT fk_social_media_platform_order FOREIGN KEY(sales_social_media_id) REFERENCES social_media_platforms(id),

  CONSTRAINT fk_shop_order FOREIGN KEY(shop_id) REFERENCES shops(id)
);

CREATE TABLE IF NOT EXISTS order_product(
  product_id CHAR(26) NOT NULL,
  order_id CHAR(26) NOT NULL,
  discount_percent NUMERIC(19,4),
  discount_sum NUMERIC(19,4),
  quantity INT,
  appearance_order INT,
  discount_reason TEXT,

  PRIMARY KEY(product_id, order_id),
  CONSTRAINT fk_product_order_product FOREIGN KEY(product_id) REFERENCES products(id),
  CONSTRAINT fk_order_order_product FOREIGN KEY(order_id) REFERENCES orders(id)
);

CREATE TABLE IF NOT EXISTS customers(
  id CHAR(26) DEFAULT ulid_to_string(ulid_generate()) NOT NULL PRIMARY KEY,
  shop_id CHAR(26) NOT NULL,
  full_name TEXT,

  created TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP WITH TIME ZONE,
  CONSTRAINT fk_shop_customer FOREIGN KEY(shop_id) REFERENCES shops(id)
);

CREATE TABLE IF NOT EXISTS customers_addresses(
  customer_id CHAR(26) NOT NULL,
  address TEXT,
  ward TEXT,
  district TEXT,
  city TEXT,
  country TEXT,

  CONSTRAINT fk_customer_customer_addresses FOREIGN KEY(customer_id) REFERENCES customers(id)
);

CREATE TABLE IF NOT EXISTS customers_phones(
  customer_id CHAR(26) NOT NULL,
  type TEXT,
  note TEXT,
  phone TEXT,

  CONSTRAINT fk_customer_customer_phones FOREIGN KEY(customer_id) REFERENCES customers(id)
);

CREATE TABLE IF NOT EXISTS customers_socials(
  platform_id INT NOT NULL,
  customer_id CHAR(26) NOT NULL,

  PRIMARY KEY(platform_id, customer_id),
  CONSTRAINT fk_customer_customer_social FOREIGN KEY(customer_id) REFERENCES customers(id),
  CONSTRAINT fk_social_platform_customer_social FOREIGN KEY(platform_id) REFERENCES social_media_platforms(id)
);

COMMIT;
