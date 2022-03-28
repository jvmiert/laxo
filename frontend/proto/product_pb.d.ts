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

