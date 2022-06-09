type ErrorDetails = {
  [key: string]: string;
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
  name: string;
  productURL: string;
};

export type LaxoProduct = {
  product: {
    id: string;
    name: string;
    description: string;
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

export interface LaxoProductDetailsResponse extends LaxoProduct {}
