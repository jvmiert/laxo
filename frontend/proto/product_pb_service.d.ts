// package: product
// file: product.proto

import * as product_pb from "./product_pb";
import {grpc} from "@improbable-eng/grpc-web";

type ProductServiceCreateFrame = {
  readonly methodName: string;
  readonly service: typeof ProductService;
  readonly requestStream: false;
  readonly responseStream: true;
  readonly requestType: typeof product_pb.CreateFrameRequest;
  readonly responseType: typeof product_pb.CreateFrameReply;
};

export class ProductService {
  static readonly serviceName: string;
  static readonly CreateFrame: ProductServiceCreateFrame;
}

export type ServiceError = { message: string, code: number; metadata: grpc.Metadata }
export type Status = { details: string, code: number; metadata: grpc.Metadata }

interface UnaryResponse {
  cancel(): void;
}
interface ResponseStream<T> {
  cancel(): void;
  on(type: 'data', handler: (message: T) => void): ResponseStream<T>;
  on(type: 'end', handler: (status?: Status) => void): ResponseStream<T>;
  on(type: 'status', handler: (status: Status) => void): ResponseStream<T>;
}
interface RequestStream<T> {
  write(message: T): RequestStream<T>;
  end(): void;
  cancel(): void;
  on(type: 'end', handler: (status?: Status) => void): RequestStream<T>;
  on(type: 'status', handler: (status: Status) => void): RequestStream<T>;
}
interface BidirectionalStream<ReqT, ResT> {
  write(message: ReqT): BidirectionalStream<ReqT, ResT>;
  end(): void;
  cancel(): void;
  on(type: 'data', handler: (message: ResT) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'end', handler: (status?: Status) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'status', handler: (status: Status) => void): BidirectionalStream<ReqT, ResT>;
}

export class ProductServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  createFrame(requestMessage: product_pb.CreateFrameRequest, metadata?: grpc.Metadata): ResponseStream<product_pb.CreateFrameReply>;
}

