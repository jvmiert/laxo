# Tech to-do

- Handle product sync to Lazada after saving product
  - Create temporal workflow/activities for syncing the Lazada product
  - Enable support for updating product images (and image order)
    - Dispatch a temporal task to update
      - Check if a task is currently in sleeping/executing mode through Redis?
      - Start the task with some debounce sleep period
      - Whenever a new image change comes in, reset the sleep period
      - After sleep period is over, execute task
    - Make sure the maximum image size for each platform is respected
  - Add support for notification when we update a product
    - Make a clear distiction between saving the Laxo product and syncing the platform product
  - Add manual sync button for products
- Create new product form
  - Do validation
- Fix ordering issue when adding new image to product
  - We are not populating the order field when linking a new image to a product
- Implement tooltips https://github.com/floating-ui/floating-ui
  - For the drag and drop image ordering
- Use valid status in Axios to stop throwing certain validation errors that return != 200 http status code
  - https://axios-http.com/docs/handling_errors
- Create email collecting backend
  - Create login frontend so production doesn't allow users to sign up yet
    - https://github.com/vercel/next.js/blob/canary/examples/with-env-from-next-config-js/next.config.js

# Long-term to-do

- Create empty states for various screens
  - No shop screen
  - No products screen
  - No product result screen
  - No assets available
- Make product dashboard page
- Make order dashboard page
  - Use ShoppingBagIcon from heroicons
- Make stock dashboard page
- Make a password strength indicator
  - https://www.openwall.com/passwdqc/
  - https://github.com/odin-public/passwdqc-js
- Figure out how to handle multiple shops / transition between shops
- Properly type the axios posts/errors/returns in my post hooks (e.g. the useOAuthApi)
- Interesting layout example for nextjs:
  - https://github.com/vercel/next.js/issues/8193#issuecomment-873281365
- Optimize yup schema validation for react final form?
  - https://gist.github.com/nfantone/9ab600760db8774ab4873cb1a3a22f26
- Setup correct meta tags / html header handeling (next/head)
  - Maybe use? https://github.com/garmeeh/next-seo
- Debounce validation to prevent stutter?
  - https://codesandbox.io/s/mmywp9jl1y?file=/index.js:101-108
- Implement CSRF protection: https://github.com/gorilla/csrf
  - https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html#double-submit-cookie
  - https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html#synchronizer-token-pattern

# Thoughts

It seems for SKU creation, the best approach is to create a logical data model for the service. This data model will then be the source of truth from which we create translation functions for each individual platform. The tricky thing will be that it seems Lazada has predefined SKU attributes. They do not allow arbitrary SKU creation. It seems Tiki also follows this strategy.

# Image notes

## Tiki

For the best user experience, TIKI only display image have size greater than 500×500 pixel in the media gallery and lower than 700 width pixel inside description

- https://open.tiki.vn/docs/docs/current/guides/tiki-theory-v2/product-v2/

## Lazada

"Use this API to upload a single image file to Lazada site. Allowed image formats are JPG and PNG. The maximum size of an image file is 1MB."

- https://open.lazada.com/doc/api.htm?spm=a2o9m.11193531.0.0.58e26bbe1KlJWh#/api?cid=5&path=/image/upload

## Shopee

"Image file. Max 2.0 MB each. Image format accepted: JPG, JPEG, PNG"

- https://open.shopee.com/documents?module=91&type=1&id=660&version=2

# Facebook integration

There is no commerce API available yet for Vietnam. However, for the business manager API, we can create product catalogues it seems.

# Frame feature research

## Image manipulation

- https://github.com/discord/lilliput
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

# Integrations

- https://developers.tiktok-shops.com/

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
- https://ginee.com/
- https://www.octosells.com/
- https://www.bigseller.com/en_US/index.htm
- https://www.linnworks.com/
- https://pipe17.com/
- https://www.bigcommerce.com/essentials/

# Random

- https://github.com/charithe/durationcheck
- https://heroicons.com/?
- https://github.com/mitranim/gow
- https://github.com/airbnb/lottie-web/tree/master/build/player

# Chatting integration / customer support

- https://github.com/chatwoot/chatwoot
- https://github.com/chaskiq/chaskiq
- https://github.com/getfider/fider#fider
- https://github.com/papercups-io/papercups
- https://github.com/mattermost/mattermost-server

# i18n service

- https://github.com/WeblateOrg/weblate

# Third party project links

- https://docs.google.com/spreadsheets/d/1CU_IFc_jVjZ7PLc6UEWmPNszufvYhhrVAGpJUemNU40/edit#gid=0

# Graphic Design Inspiration

- https://www.shopify.com/
- https://workos.com/
- https://vercel.com/login/email?
- https://www.airtable.com/
- https://www.float.com/
- https://clickup.com/
- https://polypane.app/
- https://www.honeycomb.io/#
- https://hireproof.io/
- https://www.convex.dev/
- https://indent.com/
- https://tailscale.com/
- https://sellerportal.deliverr.com/
- https://codesandbox.io/s/dp7to

# Programming

- For API design: https://stripe.com/docs/api/promotion_codes/list
- React inspiration: https://github.com/alan2207/bulletproof-react/tree/master/src
- imagick reference: https://github.com/gographics/imagick/tree/im-7/examples
- https://github.com/dedupeio/dedupe
- Accessing raw JSON: https://github.com/tidwall/gjson
