type ErrorDetails = {
  [key: string]: string;
};

export type ResponseError = {
  success: boolean;
  error: boolean;
  errorDetails: ErrorDetails;
};
