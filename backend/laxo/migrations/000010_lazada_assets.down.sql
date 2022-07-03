BEGIN;
ALTER TABLE assets_lazada DROP CONSTRAINT fk_assets_assets_lazada;

DROP TABLE IF EXISTS assets_lazada;
COMMIT;
