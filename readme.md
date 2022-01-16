# API access
Seems that currently Lazada and Tiki easily allow individual third parties to open accounts for creating applications. However, Shopee says it will only accept applications from real businesses.

# To-do
- Need to figure out the category/attribute options for both Tiki and Lazada. See if they are compatible or how can map between them.

# Tech to-do
- Create frontend login / session storage
- Create a general API error structure and implement easy logic to consistently return on every endpoint

# Thoughts
It seems for SKU creation, the best approach is to create a logical data model for the service. This data model will then be the source of truth from which we create translation functions for each individual platform. The tricky thing will be that it seems Lazada has predefined SKU attributes. They do not allow arbitrary SKU creation. It seems Tiki also follows this strategy.

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

# Third party project links
- https://docs.google.com/spreadsheets/d/1CU_IFc_jVjZ7PLc6UEWmPNszufvYhhrVAGpJUemNU40/edit#gid=0
