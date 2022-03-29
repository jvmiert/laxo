BEGIN;
CREATE TABLE IF NOT EXISTS shops(
  id CHAR(26) DEFAULT ulid_create() NOT NULL PRIMARY KEY,
  user_id CHAR(26) NOT NULL,
  shop_name VARCHAR(300) UNIQUE NOT NULL,
  created TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  last_update TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_shop_user FOREIGN KEY(user_id) REFERENCES users(id)
);
COMMIT;
