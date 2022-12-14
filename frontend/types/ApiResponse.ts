type ErrorDetailsObject = {
  code: string;
  error: string;
};
type ErrorDetails = {
  [key: string]: ErrorDetailsObject;
};

export type ResponseError = {
  success: boolean;
  error: boolean;
  errorDetails: ErrorDetails;
};

export type GetShopResponsePlatforms = {
  id: string;
  name: string;
  created: number;
};

export type GetShopResponseShops = {
  id: string;
  userID: string;
  name: string;
  assetsToken: string;
  platforms: Array<GetShopResponsePlatforms>;
};

export type GetShopResponse = {
  shops: Array<GetShopResponseShops>;
  total: number;
};

export type GetPlatformsPlatforms = {
  platform: string;
  url: string;
};

export type GetPlatformsResponse = {
  shopID: string;
  platforms: Array<GetPlatformsPlatforms>;
  connectedPlatforms: Array<string>;
};

export type Notification = {
  id: string;
  redisID: string;
  notificationGroupID: string;
  created: Date;
  read?: Date;
  currentMainStep: number;
  currentSubStep?: number;
  mainMessage: string;
  subMessage: string;
  error: boolean;
};

export type NotificationGroup = {
  id: string;
  userID: string;
  workflowID: string;
  entityID: string;
  entityType: string;
  totalMainSteps: number;
  totalSubSteps?: number;
  platformName: string;
};

export type NotificationResponseObject = {
  notification: Notification;
  notificationGroup: NotificationGroup;
};

export type GetNotificationResponse = {
  notifications: Array<NotificationResponseObject>;
  total: number;
};

export type GetLazadaPlatformResponse = {
  id: string;
  shopID: string;
  country: string;
  accountplatform: string;
  account: string;
  userIDVn: string;
  sellerIDVn: string;
  shortCodeVn: string;
  refreshExpiresIn: Date;
  accessExpiresIn: Date;
  created: Date;
};

export type LaxoProductPlatforms = {
  id: string;
  platformName: string;
  platformSKU: string;
  name: string;
  status: string;
  productURL: string;
  syncStatus: boolean;
};

export type LaxoProductAsset = {
  id: string;
  originalFilename: string;
  extension: string;
  status: string;
  fileSize: number;
  order: number;
  width: number;
  height: number;
  created: Date;
};

export type LaxoProduct = {
  product: {
    id: string;
    name: string;
    description: string;
    descriptionSlate: string;
    msku: string;
    shopID: string;
    mediaID?: string;
    created: Date;
    updated: Date;
    sellingPrice: {
      Int: number;
      Exp: number;
      Status: number;
    };
    costPrice: {
      Int: number;
      Exp: number;
      Status: number;
    };
  };
  mediaList: Array<string>;
  platforms: Array<LaxoProductPlatforms>;
};

export type LaxoProductDetails = Omit<LaxoProduct, "mediaList"> & {
  mediaList: Array<LaxoProductAsset>;
};

export type PaginateObject = {
  total: number;
  pages: number;
  limit: number;
  offset: number;
};

export type LaxoProductResponse = {
  products: Array<LaxoProduct>;
  paginate: PaginateObject;
};

export type LaxoProductDetailsResponse = LaxoProductDetails;

export type LaxoAssetResponse = {
  assets: Array<LaxoProductAsset>;
  paginate: PaginateObject;
};
