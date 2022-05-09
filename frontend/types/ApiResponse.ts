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

export type Notification = {
  id: string;
  redisID: string;
  notificationGroupID: string;
  created: Date;
  read?: Date;
  currentMainStep: number;
  currentSubStep: number;
  mainMessage: string;
  subMessage: string;
};

export type NotificationGroup = {
  id: string;
  userID: string;
  workflowID: string;
  entityID: string;
  entityType: string;
  totalMainSteps: number;
  totalSubSteps: number;
};

export type NotificationResponseObject = {
  notification: Notification;
  notificationGroup: NotificationGroup;
};

export type GetNotificationResponse = {
  notifications: Array<NotificationResponseObject>;
  total: number;
};
