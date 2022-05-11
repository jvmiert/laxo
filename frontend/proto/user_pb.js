// source: user.proto
/**
 * @fileoverview
 * @enhanceable
 * @suppress {missingRequire} reports error on implicit type usages.
 * @suppress {messageConventions} JS Compiler reports an error if a variable or
 *     field starts with 'MSG_' and isn't a translatable message.
 * @public
 */
// GENERATED CODE -- DO NOT EDIT!
/* eslint-disable */
// @ts-nocheck

var jspb = require('google-protobuf');
var goog = jspb;
var global = (function() {
  if (this) { return this; }
  if (typeof window !== 'undefined') { return window; }
  if (typeof global !== 'undefined') { return global; }
  if (typeof self !== 'undefined') { return self; }
  return Function('return this')();
}.call(null));

goog.exportSymbol('proto.user.Notification', null, global);
goog.exportSymbol('proto.user.NotificationGroup', null, global);
goog.exportSymbol('proto.user.NotificationUpdateReply', null, global);
goog.exportSymbol('proto.user.NotificationUpdateRequest', null, global);
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.user.Notification = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.user.Notification, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.user.Notification.displayName = 'proto.user.Notification';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.user.NotificationGroup = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.user.NotificationGroup, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.user.NotificationGroup.displayName = 'proto.user.NotificationGroup';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.user.NotificationUpdateRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.user.NotificationUpdateRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.user.NotificationUpdateRequest.displayName = 'proto.user.NotificationUpdateRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.user.NotificationUpdateReply = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.user.NotificationUpdateReply, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.user.NotificationUpdateReply.displayName = 'proto.user.NotificationUpdateReply';
}



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.user.Notification.prototype.toObject = function(opt_includeInstance) {
  return proto.user.Notification.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.user.Notification} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.user.Notification.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    redisid: jspb.Message.getFieldWithDefault(msg, 2, ""),
    groupid: jspb.Message.getFieldWithDefault(msg, 3, ""),
    created: jspb.Message.getFieldWithDefault(msg, 4, 0),
    read: jspb.Message.getFieldWithDefault(msg, 5, 0),
    currentmainstep: jspb.Message.getFieldWithDefault(msg, 6, 0),
    currentsubstep: jspb.Message.getFieldWithDefault(msg, 7, 0),
    mainmessage: jspb.Message.getFieldWithDefault(msg, 8, ""),
    submessage: jspb.Message.getFieldWithDefault(msg, 9, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.user.Notification}
 */
proto.user.Notification.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.user.Notification;
  return proto.user.Notification.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.user.Notification} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.user.Notification}
 */
proto.user.Notification.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setRedisid(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setGroupid(value);
      break;
    case 4:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setCreated(value);
      break;
    case 5:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setRead(value);
      break;
    case 6:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setCurrentmainstep(value);
      break;
    case 7:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setCurrentsubstep(value);
      break;
    case 8:
      var value = /** @type {string} */ (reader.readString());
      msg.setMainmessage(value);
      break;
    case 9:
      var value = /** @type {string} */ (reader.readString());
      msg.setSubmessage(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.user.Notification.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.user.Notification.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.user.Notification} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.user.Notification.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getRedisid();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getGroupid();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getCreated();
  if (f !== 0) {
    writer.writeInt64(
      4,
      f
    );
  }
  f = /** @type {number} */ (jspb.Message.getField(message, 5));
  if (f != null) {
    writer.writeInt64(
      5,
      f
    );
  }
  f = message.getCurrentmainstep();
  if (f !== 0) {
    writer.writeInt64(
      6,
      f
    );
  }
  f = /** @type {number} */ (jspb.Message.getField(message, 7));
  if (f != null) {
    writer.writeInt64(
      7,
      f
    );
  }
  f = message.getMainmessage();
  if (f.length > 0) {
    writer.writeString(
      8,
      f
    );
  }
  f = message.getSubmessage();
  if (f.length > 0) {
    writer.writeString(
      9,
      f
    );
  }
};


/**
 * optional string ID = 1;
 * @return {string}
 */
proto.user.Notification.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.user.Notification} returns this
 */
proto.user.Notification.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string redisID = 2;
 * @return {string}
 */
proto.user.Notification.prototype.getRedisid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.user.Notification} returns this
 */
proto.user.Notification.prototype.setRedisid = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string groupID = 3;
 * @return {string}
 */
proto.user.Notification.prototype.getGroupid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.user.Notification} returns this
 */
proto.user.Notification.prototype.setGroupid = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional int64 created = 4;
 * @return {number}
 */
proto.user.Notification.prototype.getCreated = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {number} value
 * @return {!proto.user.Notification} returns this
 */
proto.user.Notification.prototype.setCreated = function(value) {
  return jspb.Message.setProto3IntField(this, 4, value);
};


/**
 * optional int64 read = 5;
 * @return {number}
 */
proto.user.Notification.prototype.getRead = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 5, 0));
};


/**
 * @param {number} value
 * @return {!proto.user.Notification} returns this
 */
proto.user.Notification.prototype.setRead = function(value) {
  return jspb.Message.setField(this, 5, value);
};


/**
 * Clears the field making it undefined.
 * @return {!proto.user.Notification} returns this
 */
proto.user.Notification.prototype.clearRead = function() {
  return jspb.Message.setField(this, 5, undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.user.Notification.prototype.hasRead = function() {
  return jspb.Message.getField(this, 5) != null;
};


/**
 * optional int64 currentMainStep = 6;
 * @return {number}
 */
proto.user.Notification.prototype.getCurrentmainstep = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 6, 0));
};


/**
 * @param {number} value
 * @return {!proto.user.Notification} returns this
 */
proto.user.Notification.prototype.setCurrentmainstep = function(value) {
  return jspb.Message.setProto3IntField(this, 6, value);
};


/**
 * optional int64 currentSubStep = 7;
 * @return {number}
 */
proto.user.Notification.prototype.getCurrentsubstep = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 7, 0));
};


/**
 * @param {number} value
 * @return {!proto.user.Notification} returns this
 */
proto.user.Notification.prototype.setCurrentsubstep = function(value) {
  return jspb.Message.setField(this, 7, value);
};


/**
 * Clears the field making it undefined.
 * @return {!proto.user.Notification} returns this
 */
proto.user.Notification.prototype.clearCurrentsubstep = function() {
  return jspb.Message.setField(this, 7, undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.user.Notification.prototype.hasCurrentsubstep = function() {
  return jspb.Message.getField(this, 7) != null;
};


/**
 * optional string mainMessage = 8;
 * @return {string}
 */
proto.user.Notification.prototype.getMainmessage = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 8, ""));
};


/**
 * @param {string} value
 * @return {!proto.user.Notification} returns this
 */
proto.user.Notification.prototype.setMainmessage = function(value) {
  return jspb.Message.setProto3StringField(this, 8, value);
};


/**
 * optional string subMessage = 9;
 * @return {string}
 */
proto.user.Notification.prototype.getSubmessage = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 9, ""));
};


/**
 * @param {string} value
 * @return {!proto.user.Notification} returns this
 */
proto.user.Notification.prototype.setSubmessage = function(value) {
  return jspb.Message.setProto3StringField(this, 9, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.user.NotificationGroup.prototype.toObject = function(opt_includeInstance) {
  return proto.user.NotificationGroup.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.user.NotificationGroup} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.user.NotificationGroup.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    userid: jspb.Message.getFieldWithDefault(msg, 2, ""),
    workflowid: jspb.Message.getFieldWithDefault(msg, 3, ""),
    entityid: jspb.Message.getFieldWithDefault(msg, 4, ""),
    entitytype: jspb.Message.getFieldWithDefault(msg, 5, ""),
    totalmainsteps: jspb.Message.getFieldWithDefault(msg, 6, 0),
    totalsubsteps: jspb.Message.getFieldWithDefault(msg, 7, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.user.NotificationGroup}
 */
proto.user.NotificationGroup.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.user.NotificationGroup;
  return proto.user.NotificationGroup.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.user.NotificationGroup} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.user.NotificationGroup}
 */
proto.user.NotificationGroup.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setUserid(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setWorkflowid(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setEntityid(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setEntitytype(value);
      break;
    case 6:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setTotalmainsteps(value);
      break;
    case 7:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setTotalsubsteps(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.user.NotificationGroup.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.user.NotificationGroup.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.user.NotificationGroup} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.user.NotificationGroup.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getUserid();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getWorkflowid();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getEntityid();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
  f = message.getEntitytype();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getTotalmainsteps();
  if (f !== 0) {
    writer.writeInt64(
      6,
      f
    );
  }
  f = /** @type {number} */ (jspb.Message.getField(message, 7));
  if (f != null) {
    writer.writeInt64(
      7,
      f
    );
  }
};


/**
 * optional string ID = 1;
 * @return {string}
 */
proto.user.NotificationGroup.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.user.NotificationGroup} returns this
 */
proto.user.NotificationGroup.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string userID = 2;
 * @return {string}
 */
proto.user.NotificationGroup.prototype.getUserid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.user.NotificationGroup} returns this
 */
proto.user.NotificationGroup.prototype.setUserid = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string workflowID = 3;
 * @return {string}
 */
proto.user.NotificationGroup.prototype.getWorkflowid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.user.NotificationGroup} returns this
 */
proto.user.NotificationGroup.prototype.setWorkflowid = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string entityID = 4;
 * @return {string}
 */
proto.user.NotificationGroup.prototype.getEntityid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * @param {string} value
 * @return {!proto.user.NotificationGroup} returns this
 */
proto.user.NotificationGroup.prototype.setEntityid = function(value) {
  return jspb.Message.setProto3StringField(this, 4, value);
};


/**
 * optional string entityType = 5;
 * @return {string}
 */
proto.user.NotificationGroup.prototype.getEntitytype = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/**
 * @param {string} value
 * @return {!proto.user.NotificationGroup} returns this
 */
proto.user.NotificationGroup.prototype.setEntitytype = function(value) {
  return jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * optional int64 totalMainSteps = 6;
 * @return {number}
 */
proto.user.NotificationGroup.prototype.getTotalmainsteps = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 6, 0));
};


/**
 * @param {number} value
 * @return {!proto.user.NotificationGroup} returns this
 */
proto.user.NotificationGroup.prototype.setTotalmainsteps = function(value) {
  return jspb.Message.setProto3IntField(this, 6, value);
};


/**
 * optional int64 totalSubSteps = 7;
 * @return {number}
 */
proto.user.NotificationGroup.prototype.getTotalsubsteps = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 7, 0));
};


/**
 * @param {number} value
 * @return {!proto.user.NotificationGroup} returns this
 */
proto.user.NotificationGroup.prototype.setTotalsubsteps = function(value) {
  return jspb.Message.setField(this, 7, value);
};


/**
 * Clears the field making it undefined.
 * @return {!proto.user.NotificationGroup} returns this
 */
proto.user.NotificationGroup.prototype.clearTotalsubsteps = function() {
  return jspb.Message.setField(this, 7, undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.user.NotificationGroup.prototype.hasTotalsubsteps = function() {
  return jspb.Message.getField(this, 7) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.user.NotificationUpdateRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.user.NotificationUpdateRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.user.NotificationUpdateRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.user.NotificationUpdateRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    notificationredisid: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.user.NotificationUpdateRequest}
 */
proto.user.NotificationUpdateRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.user.NotificationUpdateRequest;
  return proto.user.NotificationUpdateRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.user.NotificationUpdateRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.user.NotificationUpdateRequest}
 */
proto.user.NotificationUpdateRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setNotificationredisid(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.user.NotificationUpdateRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.user.NotificationUpdateRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.user.NotificationUpdateRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.user.NotificationUpdateRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getNotificationredisid();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string notificationRedisID = 1;
 * @return {string}
 */
proto.user.NotificationUpdateRequest.prototype.getNotificationredisid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.user.NotificationUpdateRequest} returns this
 */
proto.user.NotificationUpdateRequest.prototype.setNotificationredisid = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.user.NotificationUpdateReply.prototype.toObject = function(opt_includeInstance) {
  return proto.user.NotificationUpdateReply.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.user.NotificationUpdateReply} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.user.NotificationUpdateReply.toObject = function(includeInstance, msg) {
  var f, obj = {
    keepalive: jspb.Message.getBooleanFieldWithDefault(msg, 1, false),
    notification: (f = msg.getNotification()) && proto.user.Notification.toObject(includeInstance, f),
    notificationgroup: (f = msg.getNotificationgroup()) && proto.user.NotificationGroup.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.user.NotificationUpdateReply}
 */
proto.user.NotificationUpdateReply.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.user.NotificationUpdateReply;
  return proto.user.NotificationUpdateReply.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.user.NotificationUpdateReply} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.user.NotificationUpdateReply}
 */
proto.user.NotificationUpdateReply.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setKeepalive(value);
      break;
    case 2:
      var value = new proto.user.Notification;
      reader.readMessage(value,proto.user.Notification.deserializeBinaryFromReader);
      msg.setNotification(value);
      break;
    case 3:
      var value = new proto.user.NotificationGroup;
      reader.readMessage(value,proto.user.NotificationGroup.deserializeBinaryFromReader);
      msg.setNotificationgroup(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.user.NotificationUpdateReply.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.user.NotificationUpdateReply.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.user.NotificationUpdateReply} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.user.NotificationUpdateReply.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getKeepalive();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
  f = message.getNotification();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.user.Notification.serializeBinaryToWriter
    );
  }
  f = message.getNotificationgroup();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.user.NotificationGroup.serializeBinaryToWriter
    );
  }
};


/**
 * optional bool keepAlive = 1;
 * @return {boolean}
 */
proto.user.NotificationUpdateReply.prototype.getKeepalive = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/**
 * @param {boolean} value
 * @return {!proto.user.NotificationUpdateReply} returns this
 */
proto.user.NotificationUpdateReply.prototype.setKeepalive = function(value) {
  return jspb.Message.setProto3BooleanField(this, 1, value);
};


/**
 * optional Notification notification = 2;
 * @return {?proto.user.Notification}
 */
proto.user.NotificationUpdateReply.prototype.getNotification = function() {
  return /** @type{?proto.user.Notification} */ (
    jspb.Message.getWrapperField(this, proto.user.Notification, 2));
};


/**
 * @param {?proto.user.Notification|undefined} value
 * @return {!proto.user.NotificationUpdateReply} returns this
*/
proto.user.NotificationUpdateReply.prototype.setNotification = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.user.NotificationUpdateReply} returns this
 */
proto.user.NotificationUpdateReply.prototype.clearNotification = function() {
  return this.setNotification(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.user.NotificationUpdateReply.prototype.hasNotification = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional NotificationGroup notificationGroup = 3;
 * @return {?proto.user.NotificationGroup}
 */
proto.user.NotificationUpdateReply.prototype.getNotificationgroup = function() {
  return /** @type{?proto.user.NotificationGroup} */ (
    jspb.Message.getWrapperField(this, proto.user.NotificationGroup, 3));
};


/**
 * @param {?proto.user.NotificationGroup|undefined} value
 * @return {!proto.user.NotificationUpdateReply} returns this
*/
proto.user.NotificationUpdateReply.prototype.setNotificationgroup = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.user.NotificationUpdateReply} returns this
 */
proto.user.NotificationUpdateReply.prototype.clearNotificationgroup = function() {
  return this.setNotificationgroup(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.user.NotificationUpdateReply.prototype.hasNotificationgroup = function() {
  return jspb.Message.getField(this, 3) != null;
};


goog.object.extend(exports, proto.user);
