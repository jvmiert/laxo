// package: product
// file: product.proto

import * as jspb from "google-protobuf";

export class CreateFrameRequest extends jspb.Message {
  getImgid(): string;
  setImgid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateFrameRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateFrameRequest): CreateFrameRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateFrameRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateFrameRequest;
  static deserializeBinaryFromReader(message: CreateFrameRequest, reader: jspb.BinaryReader): CreateFrameRequest;
}

export namespace CreateFrameRequest {
  export type AsObject = {
    imgid: string,
  }
}

export class CreateFrameReply extends jspb.Message {
  getImgid(): string;
  setImgid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateFrameReply.AsObject;
  static toObject(includeInstance: boolean, msg: CreateFrameReply): CreateFrameReply.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateFrameReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateFrameReply;
  static deserializeBinaryFromReader(message: CreateFrameReply, reader: jspb.BinaryReader): CreateFrameReply;
}

export namespace CreateFrameReply {
  export type AsObject = {
    imgid: string,
  }
}

export class ProductRetrieveUpdateRequest extends jspb.Message {
  getShopid(): string;
  setShopid(value: string): void;

  getRetrieveid(): string;
  setRetrieveid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProductRetrieveUpdateRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ProductRetrieveUpdateRequest): ProductRetrieveUpdateRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ProductRetrieveUpdateRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProductRetrieveUpdateRequest;
  static deserializeBinaryFromReader(message: ProductRetrieveUpdateRequest, reader: jspb.BinaryReader): ProductRetrieveUpdateRequest;
}

export namespace ProductRetrieveUpdateRequest {
  export type AsObject = {
    shopid: string,
    retrieveid: string,
  }
}

export class ProductRetrieveUpdateReply extends jspb.Message {
  getCurrentstatus(): string;
  setCurrentstatus(value: string): void;

  getTotalproducts(): number;
  setTotalproducts(value: number): void;

  getCurrentproducts(): number;
  setCurrentproducts(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProductRetrieveUpdateReply.AsObject;
  static toObject(includeInstance: boolean, msg: ProductRetrieveUpdateReply): ProductRetrieveUpdateReply.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ProductRetrieveUpdateReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProductRetrieveUpdateReply;
  static deserializeBinaryFromReader(message: ProductRetrieveUpdateReply, reader: jspb.BinaryReader): ProductRetrieveUpdateReply;
}

export namespace ProductRetrieveUpdateReply {
  export type AsObject = {
    currentstatus: string,
    totalproducts: number,
    currentproducts: number,
  }
}

