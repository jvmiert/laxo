// package: product
// file: product.proto

var product_pb = require("./product_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var ProductService = (function () {
  function ProductService() {}
  ProductService.serviceName = "product.ProductService";
  return ProductService;
}());

ProductService.CreateFrame = {
  methodName: "CreateFrame",
  service: ProductService,
  requestStream: false,
  responseStream: true,
  requestType: product_pb.CreateFrameRequest,
  responseType: product_pb.CreateFrameReply
};

exports.ProductService = ProductService;

function ProductServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

ProductServiceClient.prototype.createFrame = function createFrame(requestMessage, metadata) {
  var listeners = {
    data: [],
    end: [],
    status: []
  };
  var client = grpc.invoke(ProductService.CreateFrame, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onMessage: function (responseMessage) {
      listeners.data.forEach(function (handler) {
        handler(responseMessage);
      });
    },
    onEnd: function (status, statusMessage, trailers) {
      listeners.status.forEach(function (handler) {
        handler({ code: status, details: statusMessage, metadata: trailers });
      });
      listeners.end.forEach(function (handler) {
        handler({ code: status, details: statusMessage, metadata: trailers });
      });
      listeners = null;
    }
  });
  return {
    on: function (type, handler) {
      listeners[type].push(handler);
      return this;
    },
    cancel: function () {
      listeners = null;
      client.close();
    }
  };
};

exports.ProductServiceClient = ProductServiceClient;

