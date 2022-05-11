// package: user
// file: user.proto

import * as jspb from "google-protobuf";

export class Notification extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getRedisid(): string;
  setRedisid(value: string): void;

  getGroupid(): string;
  setGroupid(value: string): void;

  getCreated(): number;
  setCreated(value: number): void;

  hasRead(): boolean;
  clearRead(): void;
  getRead(): number;
  setRead(value: number): void;

  getCurrentmainstep(): number;
  setCurrentmainstep(value: number): void;

  hasCurrentsubstep(): boolean;
  clearCurrentsubstep(): void;
  getCurrentsubstep(): number;
  setCurrentsubstep(value: number): void;

  getMainmessage(): string;
  setMainmessage(value: string): void;

  getSubmessage(): string;
  setSubmessage(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Notification.AsObject;
  static toObject(includeInstance: boolean, msg: Notification): Notification.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Notification, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Notification;
  static deserializeBinaryFromReader(message: Notification, reader: jspb.BinaryReader): Notification;
}

export namespace Notification {
  export type AsObject = {
    id: string,
    redisid: string,
    groupid: string,
    created: number,
    read: number,
    currentmainstep: number,
    currentsubstep: number,
    mainmessage: string,
    submessage: string,
  }
}

export class NotificationGroup extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getUserid(): string;
  setUserid(value: string): void;

  getWorkflowid(): string;
  setWorkflowid(value: string): void;

  getEntityid(): string;
  setEntityid(value: string): void;

  getEntitytype(): string;
  setEntitytype(value: string): void;

  getTotalmainsteps(): number;
  setTotalmainsteps(value: number): void;

  hasTotalsubsteps(): boolean;
  clearTotalsubsteps(): void;
  getTotalsubsteps(): number;
  setTotalsubsteps(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NotificationGroup.AsObject;
  static toObject(includeInstance: boolean, msg: NotificationGroup): NotificationGroup.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: NotificationGroup, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NotificationGroup;
  static deserializeBinaryFromReader(message: NotificationGroup, reader: jspb.BinaryReader): NotificationGroup;
}

export namespace NotificationGroup {
  export type AsObject = {
    id: string,
    userid: string,
    workflowid: string,
    entityid: string,
    entitytype: string,
    totalmainsteps: number,
    totalsubsteps: number,
  }
}

export class NotificationUpdateRequest extends jspb.Message {
  getNotificationredisid(): string;
  setNotificationredisid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NotificationUpdateRequest.AsObject;
  static toObject(includeInstance: boolean, msg: NotificationUpdateRequest): NotificationUpdateRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: NotificationUpdateRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NotificationUpdateRequest;
  static deserializeBinaryFromReader(message: NotificationUpdateRequest, reader: jspb.BinaryReader): NotificationUpdateRequest;
}

export namespace NotificationUpdateRequest {
  export type AsObject = {
    notificationredisid: string,
  }
}

export class NotificationUpdateReply extends jspb.Message {
  getKeepalive(): boolean;
  setKeepalive(value: boolean): void;

  hasNotification(): boolean;
  clearNotification(): void;
  getNotification(): Notification | undefined;
  setNotification(value?: Notification): void;

  hasNotificationgroup(): boolean;
  clearNotificationgroup(): void;
  getNotificationgroup(): NotificationGroup | undefined;
  setNotificationgroup(value?: NotificationGroup): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NotificationUpdateReply.AsObject;
  static toObject(includeInstance: boolean, msg: NotificationUpdateReply): NotificationUpdateReply.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: NotificationUpdateReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NotificationUpdateReply;
  static deserializeBinaryFromReader(message: NotificationUpdateReply, reader: jspb.BinaryReader): NotificationUpdateReply;
}

export namespace NotificationUpdateReply {
  export type AsObject = {
    keepalive: boolean,
    notification?: Notification.AsObject,
    notificationgroup?: NotificationGroup.AsObject,
  }
}

