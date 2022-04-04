# API access
Seems that currently Lazada and Tiki easily allow individual third parties to open accounts for creating applications. However, Shopee says it will only accept applications from real businesses.

# To-do
- Need to figure out the category/attribute options for both Tiki and Lazada. See if they are compatible or how can map between them.

# Tech to-do
- Make the dashboard page
  - Setup a callback url on nextjs under setup-shop/callback for each platform's oauth callback
  - Setup my webserver to be able to receive callback requests from the platforms
  - Create database logic that tracks what platform is connected to what shop
  - Setup frontend to redirect to oauth authorization url of each platform on button connect click
  - Process the return callbacks and verify them
  - Store the access tokens in the database
- Write test for get shop
- Backend API should return error_code, this code will be used for i18n on the frontend. The backend can
  also return a human readable message in the requested locale.
  - https://go.dev/tour/methods/15 (make sure the type casting is error checked)
- Use css variables for setting primary/secondary colors
  - https://tailwindcss.com/docs/customizing-colors#naming-your-colors
- Setup frontend to redirect default language to its /lang endpoint with nextjs middleware
  - https://nextjs.org/docs/advanced-features/i18n-routing#prefixing-the-default-locale
- Debounce validation to prevent stutter?
  - https://codesandbox.io/s/mmywp9jl1y?file=/index.js:101-108

# Long-term to-do
- Figure out how to handle multiple shops / transition between shops
- Interesting layout example for nextjs:
  - https://github.com/vercel/next.js/issues/8193#issuecomment-873281365
- Implement: https://github.com/gorilla/csrf
  - https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html#synchronizer-token-pattern
- Optimize yup schema validation for react final form?
  - https://gist.github.com/nfantone/9ab600760db8774ab4873cb1a3a22f26
- Setup correct meta tags / html header handeling (next/head)
  - Maybe use? https://github.com/garmeeh/next-seo

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

# Frame feature research

## Image manipulation
- https://github.com/gographics/imagick
- https://github.com/omgovich/react-colorful#customization
- https://github.com/jgraph/drawio/blob/dev/src/main/webapp/js/grapheditor/Shapes.js
- https://imagemagick.org/api/drawing-wand.php
- Fast image manipulation:
  - https://github.com/h2non/bimg

## Job queues
- https://cadenceworkflow.io/docs/go-client/
- https://docs.temporal.io/docs/temporal-explained/introduction
  - https://github.com/temporalio/samples-go
- https://github.com/hibiken/asynq

# Oso - authorization
- https://news.ycombinator.com/item?id=30878926
- https://docs.osohq.com/go/getting-started/quickstart.html
- https://www.osohq.com/academy/microservices-authorization
- https://www.osohq.com/post/microservices-authorization-patterns
- https://www.osohq.com/post/why-authorization-is-hard

# Temporal
- https://github.com/DataDog/temporalite
- https://github.com/temporalio/temporal
- https://github.com/temporalio/tctl
- https://github.com/temporalio/docker-builds
- https://docs.temporal.io/docs/server/options

# GRPC
- https://kennethjenkins.net/posts/go-nginx-grpc/
- https://github.com/improbable-eng/grpc-web/go/grpcweb
- https://github.com/percybolmer/grpcexample/blob/master/main.go
- https://github.com/grpc-ecosystem/awesome-grpc

# Links
- https://open.tiki.vn/docs/docs/current/getting-started/
- https://open.lazada.com/
- https://open.shopee.com/

# Competitors
- https://ecomkey.asia/en/price-list
- https://www.anchanto.com
- https://www.crescodata.com/
- https://www.onpoint.vn/

# Random
- https://github.com/charithe/durationcheck

# Third party project links
- https://docs.google.com/spreadsheets/d/1CU_IFc_jVjZ7PLc6UEWmPNszufvYhhrVAGpJUemNU40/edit#gid=0

# Graphic Design Inspiration
- https://workos.com/
- https://vercel.com/login/email?
- https://monday.com/
- https://asana.com/#close
- https://www.airtable.com/
- https://www.float.com/
- https://clickup.com/

# Programming
- For API design: https://stripe.com/docs/api/promotion_codes/list

