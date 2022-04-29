# SKU
A description of the sku field meaning.

## SellerSku
This field is set by the customer and can be 200 characters long.

## ShopSku
This field is set by Lazada and identifies the shop specific SKU.

## SkuId
No documentation found but it seems a global identifier for the generated SKU in integer form.


# Variations
For whatever reason the variations are not clearly returned from the product listing API. In order to understand
what variation attributes are present on the item, we need to call the item detail API it seems.

https://open.lazada.com/doc/doc.htm?spm=a2o9m.11193535.0.0.4fa438e4OGpzVa#?nodeId=10557&docId=108146
