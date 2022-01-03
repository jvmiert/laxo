Seems that currently Lazada and Tiki easily allow individual third parties to open accounts for creating applications. However, Shopee says it will only accept applications from real businesses.

# To-do
- Setup SQLC based on the migration file that I already made
  - https://docs.sqlc.dev/en/stable/tutorials/getting-started-postgresql.html
  - https://github.com/search?q=packages+schema++filename%3Asqlc.yaml&type=Code&ref=advsearch&l=&l=
- Setup Makefile for SQLC to automatically generate with correct version (use go tooling just like golang-migrate)
  - https://github.com/msal4/sdms-grpc/blob/2ce950e6e65f3bc91a04a971fde79694fc4d9bb5/Makefile
- Figure out the data model for products in each platform
  - Tiki
  - Shopee
  - Lazada
- Figure out product image sizes on each platform
  - Tiki
  - Shopee
  - Lazada

# Secrets

## Tiki
Application ID
5428082774690597

Application Secret
LROE3g1vd0yPQgcqOv6spa1z39Ep2jiS

# Image notes

## Tiki
"*For the best user experience, TIKI only display image have size greater than 500×500 pixel in the media gallery and lower than 700 width pixel inside description"
- https://open.tiki.vn/docs/docs/current/guides/tiki-theory-v2/product-v2/

## Lazada
"Use this API to upload a single image file to Lazada site. Allowed image formats are JPG and PNG. The maximum size of an image file is 1MB."
- https://open.lazada.com/doc/api.htm?spm=a2o9m.11193531.0.0.58e26bbe1KlJWh#/api?cid=5&path=/image/upload

## Shopee
"Image file. Max 2.0 MB each. Image format accepted: JPG, JPEG, PNG"
- https://open.shopee.com/documents?module=91&type=1&id=660&version=2

# Links
- https://open.tiki.vn/docs/docs/current/getting-started/
- https://open.lazada.com/
- https://open.shopee.com/

# Competitors

- https://ecomkey.asia/en/price-list
- https://www.anchanto.com
- https://www.crescodata.com/
