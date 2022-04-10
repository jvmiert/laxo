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
};
