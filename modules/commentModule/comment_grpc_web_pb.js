/**
 * @fileoverview gRPC-Web generated client stub for commentModule
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.commentModule = require('./comment_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.commentModule.CommentServiceClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

  /**
   * @private @const {?Object} The credentials to be used to connect
   *    to the server
   */
  this.credentials_ = credentials;

  /**
   * @private @const {?Object} Options for the client
   */
  this.options_ = options;
};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.commentModule.CommentServicePromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

  /**
   * @private @const {?Object} The credentials to be used to connect
   *    to the server
   */
  this.credentials_ = credentials;

  /**
   * @private @const {?Object} Options for the client
   */
  this.options_ = options;
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.commentModule.RequestGetComments,
 *   !proto.commentModule.ResponseGetComments>}
 */
const methodDescriptor_CommentService_getComments = new grpc.web.MethodDescriptor(
  '/commentModule.CommentService/getComments',
  grpc.web.MethodType.UNARY,
  proto.commentModule.RequestGetComments,
  proto.commentModule.ResponseGetComments,
  /** @param {!proto.commentModule.RequestGetComments} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.commentModule.ResponseGetComments.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.commentModule.RequestGetComments,
 *   !proto.commentModule.ResponseGetComments>}
 */
const methodInfo_CommentService_getComments = new grpc.web.AbstractClientBase.MethodInfo(
  proto.commentModule.ResponseGetComments,
  /** @param {!proto.commentModule.RequestGetComments} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.commentModule.ResponseGetComments.deserializeBinary
);


/**
 * @param {!proto.commentModule.RequestGetComments} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.commentModule.ResponseGetComments)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.commentModule.ResponseGetComments>|undefined}
 *     The XHR Node Readable Stream
 */
proto.commentModule.CommentServiceClient.prototype.getComments =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/commentModule.CommentService/getComments',
      request,
      metadata || {},
      methodDescriptor_CommentService_getComments,
      callback);
};


/**
 * @param {!proto.commentModule.RequestGetComments} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.commentModule.ResponseGetComments>}
 *     A native promise that resolves to the response
 */
proto.commentModule.CommentServicePromiseClient.prototype.getComments =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/commentModule.CommentService/getComments',
      request,
      metadata || {},
      methodDescriptor_CommentService_getComments);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.commentModule.RequestPutComment,
 *   !proto.commentModule.ResponsePutComment>}
 */
const methodDescriptor_CommentService_putComment = new grpc.web.MethodDescriptor(
  '/commentModule.CommentService/putComment',
  grpc.web.MethodType.UNARY,
  proto.commentModule.RequestPutComment,
  proto.commentModule.ResponsePutComment,
  /** @param {!proto.commentModule.RequestPutComment} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.commentModule.ResponsePutComment.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.commentModule.RequestPutComment,
 *   !proto.commentModule.ResponsePutComment>}
 */
const methodInfo_CommentService_putComment = new grpc.web.AbstractClientBase.MethodInfo(
  proto.commentModule.ResponsePutComment,
  /** @param {!proto.commentModule.RequestPutComment} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.commentModule.ResponsePutComment.deserializeBinary
);


/**
 * @param {!proto.commentModule.RequestPutComment} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.commentModule.ResponsePutComment)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.commentModule.ResponsePutComment>|undefined}
 *     The XHR Node Readable Stream
 */
proto.commentModule.CommentServiceClient.prototype.putComment =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/commentModule.CommentService/putComment',
      request,
      metadata || {},
      methodDescriptor_CommentService_putComment,
      callback);
};


/**
 * @param {!proto.commentModule.RequestPutComment} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.commentModule.ResponsePutComment>}
 *     A native promise that resolves to the response
 */
proto.commentModule.CommentServicePromiseClient.prototype.putComment =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/commentModule.CommentService/putComment',
      request,
      metadata || {},
      methodDescriptor_CommentService_putComment);
};


module.exports = proto.commentModule;

