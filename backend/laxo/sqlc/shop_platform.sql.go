// Code generated by sqlc. DO NOT EDIT.
// source: shop_platform.sql

package sqlc

import (
	"context"
)

const getPlatformsByShopID = `-- name: GetPlatformsByShopID :many
SELECT id, shop_id, platform_name, created FROM shops_platforms
WHERE shop_id = $1
ORDER BY shops_platforms.platform_name
`

func (q *Queries) GetPlatformsByShopID(ctx context.Context, shopID string) ([]ShopsPlatform, error) {
	rows, err := q.db.Query(ctx, getPlatformsByShopID, shopID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ShopsPlatform
	for rows.Next() {
		var i ShopsPlatform
		if err := rows.Scan(
			&i.ID,
			&i.ShopID,
			&i.PlatformName,
			&i.Created,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
