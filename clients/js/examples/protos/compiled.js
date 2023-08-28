/*eslint-disable block-scoped-var, id-length, no-control-regex, no-magic-numbers, no-prototype-builtins, no-redeclare, no-shadow, no-var, sort-vars*/
import $protobuf from "protobufjs";

// Common aliases
const $Reader = $protobuf.Reader, $Writer = $protobuf.Writer, $util = $protobuf.util;

// Exported root namespace
const $root = $protobuf.roots["default"] || ($protobuf.roots["default"] = {});

export const clickevents = $root.clickevents = (() => {

    /**
     * Namespace clickevents.
     * @exports clickevents
     * @namespace
     */
    const clickevents = {};

    clickevents.ClickEvent = (function() {

        /**
         * Properties of a ClickEvent.
         * @memberof clickevents
         * @interface IClickEvent
         * @property {string|null} [eventGuid] ClickEvent eventGuid
         * @property {number|Long|null} [componentIndex] ClickEvent componentIndex
         * @property {string|null} [componentName] ClickEvent componentName
         * @property {google.protobuf.ITimestamp|null} [sentTime] ClickEvent sentTime
         */

        /**
         * Constructs a new ClickEvent.
         * @memberof clickevents
         * @classdesc Represents a ClickEvent.
         * @implements IClickEvent
         * @constructor
         * @param {clickevents.IClickEvent=} [properties] Properties to set
         */
        function ClickEvent(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * ClickEvent eventGuid.
         * @member {string} eventGuid
         * @memberof clickevents.ClickEvent
         * @instance
         */
        ClickEvent.prototype.eventGuid = "";

        /**
         * ClickEvent componentIndex.
         * @member {number|Long} componentIndex
         * @memberof clickevents.ClickEvent
         * @instance
         */
        ClickEvent.prototype.componentIndex = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

        /**
         * ClickEvent componentName.
         * @member {string} componentName
         * @memberof clickevents.ClickEvent
         * @instance
         */
        ClickEvent.prototype.componentName = "";

        /**
         * ClickEvent sentTime.
         * @member {google.protobuf.ITimestamp|null|undefined} sentTime
         * @memberof clickevents.ClickEvent
         * @instance
         */
        ClickEvent.prototype.sentTime = null;

        /**
         * Creates a new ClickEvent instance using the specified properties.
         * @function create
         * @memberof clickevents.ClickEvent
         * @static
         * @param {clickevents.IClickEvent=} [properties] Properties to set
         * @returns {clickevents.ClickEvent} ClickEvent instance
         */
        ClickEvent.create = function create(properties) {
            return new ClickEvent(properties);
        };

        /**
         * Encodes the specified ClickEvent message. Does not implicitly {@link clickevents.ClickEvent.verify|verify} messages.
         * @function encode
         * @memberof clickevents.ClickEvent
         * @static
         * @param {clickevents.IClickEvent} message ClickEvent message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        ClickEvent.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.eventGuid != null && Object.hasOwnProperty.call(message, "eventGuid"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.eventGuid);
            if (message.componentIndex != null && Object.hasOwnProperty.call(message, "componentIndex"))
                writer.uint32(/* id 2, wireType 0 =*/16).int64(message.componentIndex);
            if (message.componentName != null && Object.hasOwnProperty.call(message, "componentName"))
                writer.uint32(/* id 3, wireType 2 =*/26).string(message.componentName);
            if (message.sentTime != null && Object.hasOwnProperty.call(message, "sentTime"))
                $root.google.protobuf.Timestamp.encode(message.sentTime, writer.uint32(/* id 4, wireType 2 =*/34).fork()).ldelim();
            return writer;
        };

        /**
         * Encodes the specified ClickEvent message, length delimited. Does not implicitly {@link clickevents.ClickEvent.verify|verify} messages.
         * @function encodeDelimited
         * @memberof clickevents.ClickEvent
         * @static
         * @param {clickevents.IClickEvent} message ClickEvent message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        ClickEvent.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a ClickEvent message from the specified reader or buffer.
         * @function decode
         * @memberof clickevents.ClickEvent
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {clickevents.ClickEvent} ClickEvent
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        ClickEvent.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.clickevents.ClickEvent();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1: {
                        message.eventGuid = reader.string();
                        break;
                    }
                case 2: {
                        message.componentIndex = reader.int64();
                        break;
                    }
                case 3: {
                        message.componentName = reader.string();
                        break;
                    }
                case 4: {
                        message.sentTime = $root.google.protobuf.Timestamp.decode(reader, reader.uint32());
                        break;
                    }
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a ClickEvent message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof clickevents.ClickEvent
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {clickevents.ClickEvent} ClickEvent
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        ClickEvent.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a ClickEvent message.
         * @function verify
         * @memberof clickevents.ClickEvent
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        ClickEvent.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.eventGuid != null && message.hasOwnProperty("eventGuid"))
                if (!$util.isString(message.eventGuid))
                    return "eventGuid: string expected";
            if (message.componentIndex != null && message.hasOwnProperty("componentIndex"))
                if (!$util.isInteger(message.componentIndex) && !(message.componentIndex && $util.isInteger(message.componentIndex.low) && $util.isInteger(message.componentIndex.high)))
                    return "componentIndex: integer|Long expected";
            if (message.componentName != null && message.hasOwnProperty("componentName"))
                if (!$util.isString(message.componentName))
                    return "componentName: string expected";
            if (message.sentTime != null && message.hasOwnProperty("sentTime")) {
                let error = $root.google.protobuf.Timestamp.verify(message.sentTime);
                if (error)
                    return "sentTime." + error;
            }
            return null;
        };

        /**
         * Creates a ClickEvent message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof clickevents.ClickEvent
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {clickevents.ClickEvent} ClickEvent
         */
        ClickEvent.fromObject = function fromObject(object) {
            if (object instanceof $root.clickevents.ClickEvent)
                return object;
            let message = new $root.clickevents.ClickEvent();
            if (object.eventGuid != null)
                message.eventGuid = String(object.eventGuid);
            if (object.componentIndex != null)
                if ($util.Long)
                    (message.componentIndex = $util.Long.fromValue(object.componentIndex)).unsigned = false;
                else if (typeof object.componentIndex === "string")
                    message.componentIndex = parseInt(object.componentIndex, 10);
                else if (typeof object.componentIndex === "number")
                    message.componentIndex = object.componentIndex;
                else if (typeof object.componentIndex === "object")
                    message.componentIndex = new $util.LongBits(object.componentIndex.low >>> 0, object.componentIndex.high >>> 0).toNumber();
            if (object.componentName != null)
                message.componentName = String(object.componentName);
            if (object.sentTime != null) {
                if (typeof object.sentTime !== "object")
                    throw TypeError(".clickevents.ClickEvent.sentTime: object expected");
                message.sentTime = $root.google.protobuf.Timestamp.fromObject(object.sentTime);
            }
            return message;
        };

        /**
         * Creates a plain object from a ClickEvent message. Also converts values to other types if specified.
         * @function toObject
         * @memberof clickevents.ClickEvent
         * @static
         * @param {clickevents.ClickEvent} message ClickEvent
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        ClickEvent.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.eventGuid = "";
                if ($util.Long) {
                    let long = new $util.Long(0, 0, false);
                    object.componentIndex = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
                } else
                    object.componentIndex = options.longs === String ? "0" : 0;
                object.componentName = "";
                object.sentTime = null;
            }
            if (message.eventGuid != null && message.hasOwnProperty("eventGuid"))
                object.eventGuid = message.eventGuid;
            if (message.componentIndex != null && message.hasOwnProperty("componentIndex"))
                if (typeof message.componentIndex === "number")
                    object.componentIndex = options.longs === String ? String(message.componentIndex) : message.componentIndex;
                else
                    object.componentIndex = options.longs === String ? $util.Long.prototype.toString.call(message.componentIndex) : options.longs === Number ? new $util.LongBits(message.componentIndex.low >>> 0, message.componentIndex.high >>> 0).toNumber() : message.componentIndex;
            if (message.componentName != null && message.hasOwnProperty("componentName"))
                object.componentName = message.componentName;
            if (message.sentTime != null && message.hasOwnProperty("sentTime"))
                object.sentTime = $root.google.protobuf.Timestamp.toObject(message.sentTime, options);
            return object;
        };

        /**
         * Converts this ClickEvent to JSON.
         * @function toJSON
         * @memberof clickevents.ClickEvent
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        ClickEvent.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        /**
         * Gets the default type url for ClickEvent
         * @function getTypeUrl
         * @memberof clickevents.ClickEvent
         * @static
         * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns {string} The default type url
         */
        ClickEvent.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
            if (typeUrlPrefix === undefined) {
                typeUrlPrefix = "type.googleapis.com";
            }
            return typeUrlPrefix + "/clickevents.ClickEvent";
        };

        return ClickEvent;
    })();

    clickevents.PageEvent = (function() {

        /**
         * Properties of a PageEvent.
         * @memberof clickevents
         * @interface IPageEvent
         * @property {string|null} [eventGuid] PageEvent eventGuid
         * @property {string|null} [eventName] PageEvent eventName
         * @property {google.protobuf.ITimestamp|null} [sentTime] PageEvent sentTime
         */

        /**
         * Constructs a new PageEvent.
         * @memberof clickevents
         * @classdesc Represents a PageEvent.
         * @implements IPageEvent
         * @constructor
         * @param {clickevents.IPageEvent=} [properties] Properties to set
         */
        function PageEvent(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * PageEvent eventGuid.
         * @member {string} eventGuid
         * @memberof clickevents.PageEvent
         * @instance
         */
        PageEvent.prototype.eventGuid = "";

        /**
         * PageEvent eventName.
         * @member {string} eventName
         * @memberof clickevents.PageEvent
         * @instance
         */
        PageEvent.prototype.eventName = "";

        /**
         * PageEvent sentTime.
         * @member {google.protobuf.ITimestamp|null|undefined} sentTime
         * @memberof clickevents.PageEvent
         * @instance
         */
        PageEvent.prototype.sentTime = null;

        /**
         * Creates a new PageEvent instance using the specified properties.
         * @function create
         * @memberof clickevents.PageEvent
         * @static
         * @param {clickevents.IPageEvent=} [properties] Properties to set
         * @returns {clickevents.PageEvent} PageEvent instance
         */
        PageEvent.create = function create(properties) {
            return new PageEvent(properties);
        };

        /**
         * Encodes the specified PageEvent message. Does not implicitly {@link clickevents.PageEvent.verify|verify} messages.
         * @function encode
         * @memberof clickevents.PageEvent
         * @static
         * @param {clickevents.IPageEvent} message PageEvent message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        PageEvent.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.eventGuid != null && Object.hasOwnProperty.call(message, "eventGuid"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.eventGuid);
            if (message.eventName != null && Object.hasOwnProperty.call(message, "eventName"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.eventName);
            if (message.sentTime != null && Object.hasOwnProperty.call(message, "sentTime"))
                $root.google.protobuf.Timestamp.encode(message.sentTime, writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
            return writer;
        };

        /**
         * Encodes the specified PageEvent message, length delimited. Does not implicitly {@link clickevents.PageEvent.verify|verify} messages.
         * @function encodeDelimited
         * @memberof clickevents.PageEvent
         * @static
         * @param {clickevents.IPageEvent} message PageEvent message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        PageEvent.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a PageEvent message from the specified reader or buffer.
         * @function decode
         * @memberof clickevents.PageEvent
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {clickevents.PageEvent} PageEvent
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        PageEvent.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.clickevents.PageEvent();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1: {
                        message.eventGuid = reader.string();
                        break;
                    }
                case 2: {
                        message.eventName = reader.string();
                        break;
                    }
                case 3: {
                        message.sentTime = $root.google.protobuf.Timestamp.decode(reader, reader.uint32());
                        break;
                    }
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a PageEvent message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof clickevents.PageEvent
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {clickevents.PageEvent} PageEvent
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        PageEvent.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a PageEvent message.
         * @function verify
         * @memberof clickevents.PageEvent
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        PageEvent.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.eventGuid != null && message.hasOwnProperty("eventGuid"))
                if (!$util.isString(message.eventGuid))
                    return "eventGuid: string expected";
            if (message.eventName != null && message.hasOwnProperty("eventName"))
                if (!$util.isString(message.eventName))
                    return "eventName: string expected";
            if (message.sentTime != null && message.hasOwnProperty("sentTime")) {
                let error = $root.google.protobuf.Timestamp.verify(message.sentTime);
                if (error)
                    return "sentTime." + error;
            }
            return null;
        };

        /**
         * Creates a PageEvent message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof clickevents.PageEvent
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {clickevents.PageEvent} PageEvent
         */
        PageEvent.fromObject = function fromObject(object) {
            if (object instanceof $root.clickevents.PageEvent)
                return object;
            let message = new $root.clickevents.PageEvent();
            if (object.eventGuid != null)
                message.eventGuid = String(object.eventGuid);
            if (object.eventName != null)
                message.eventName = String(object.eventName);
            if (object.sentTime != null) {
                if (typeof object.sentTime !== "object")
                    throw TypeError(".clickevents.PageEvent.sentTime: object expected");
                message.sentTime = $root.google.protobuf.Timestamp.fromObject(object.sentTime);
            }
            return message;
        };

        /**
         * Creates a plain object from a PageEvent message. Also converts values to other types if specified.
         * @function toObject
         * @memberof clickevents.PageEvent
         * @static
         * @param {clickevents.PageEvent} message PageEvent
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        PageEvent.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.eventGuid = "";
                object.eventName = "";
                object.sentTime = null;
            }
            if (message.eventGuid != null && message.hasOwnProperty("eventGuid"))
                object.eventGuid = message.eventGuid;
            if (message.eventName != null && message.hasOwnProperty("eventName"))
                object.eventName = message.eventName;
            if (message.sentTime != null && message.hasOwnProperty("sentTime"))
                object.sentTime = $root.google.protobuf.Timestamp.toObject(message.sentTime, options);
            return object;
        };

        /**
         * Converts this PageEvent to JSON.
         * @function toJSON
         * @memberof clickevents.PageEvent
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        PageEvent.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        /**
         * Gets the default type url for PageEvent
         * @function getTypeUrl
         * @memberof clickevents.PageEvent
         * @static
         * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns {string} The default type url
         */
        PageEvent.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
            if (typeUrlPrefix === undefined) {
                typeUrlPrefix = "type.googleapis.com";
            }
            return typeUrlPrefix + "/clickevents.PageEvent";
        };

        return PageEvent;
    })();

    return clickevents;
})();

export const google = $root.google = (() => {

    /**
     * Namespace google.
     * @exports google
     * @namespace
     */
    const google = {};

    google.protobuf = (function() {

        /**
         * Namespace protobuf.
         * @memberof google
         * @namespace
         */
        const protobuf = {};

        protobuf.Timestamp = (function() {

            /**
             * Properties of a Timestamp.
             * @memberof google.protobuf
             * @interface ITimestamp
             * @property {number|Long|null} [seconds] Timestamp seconds
             * @property {number|null} [nanos] Timestamp nanos
             */

            /**
             * Constructs a new Timestamp.
             * @memberof google.protobuf
             * @classdesc Represents a Timestamp.
             * @implements ITimestamp
             * @constructor
             * @param {google.protobuf.ITimestamp=} [properties] Properties to set
             */
            function Timestamp(properties) {
                if (properties)
                    for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                        if (properties[keys[i]] != null)
                            this[keys[i]] = properties[keys[i]];
            }

            /**
             * Timestamp seconds.
             * @member {number|Long} seconds
             * @memberof google.protobuf.Timestamp
             * @instance
             */
            Timestamp.prototype.seconds = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

            /**
             * Timestamp nanos.
             * @member {number} nanos
             * @memberof google.protobuf.Timestamp
             * @instance
             */
            Timestamp.prototype.nanos = 0;

            /**
             * Creates a new Timestamp instance using the specified properties.
             * @function create
             * @memberof google.protobuf.Timestamp
             * @static
             * @param {google.protobuf.ITimestamp=} [properties] Properties to set
             * @returns {google.protobuf.Timestamp} Timestamp instance
             */
            Timestamp.create = function create(properties) {
                return new Timestamp(properties);
            };

            /**
             * Encodes the specified Timestamp message. Does not implicitly {@link google.protobuf.Timestamp.verify|verify} messages.
             * @function encode
             * @memberof google.protobuf.Timestamp
             * @static
             * @param {google.protobuf.ITimestamp} message Timestamp message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            Timestamp.encode = function encode(message, writer) {
                if (!writer)
                    writer = $Writer.create();
                if (message.seconds != null && Object.hasOwnProperty.call(message, "seconds"))
                    writer.uint32(/* id 1, wireType 0 =*/8).int64(message.seconds);
                if (message.nanos != null && Object.hasOwnProperty.call(message, "nanos"))
                    writer.uint32(/* id 2, wireType 0 =*/16).int32(message.nanos);
                return writer;
            };

            /**
             * Encodes the specified Timestamp message, length delimited. Does not implicitly {@link google.protobuf.Timestamp.verify|verify} messages.
             * @function encodeDelimited
             * @memberof google.protobuf.Timestamp
             * @static
             * @param {google.protobuf.ITimestamp} message Timestamp message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            Timestamp.encodeDelimited = function encodeDelimited(message, writer) {
                return this.encode(message, writer).ldelim();
            };

            /**
             * Decodes a Timestamp message from the specified reader or buffer.
             * @function decode
             * @memberof google.protobuf.Timestamp
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @param {number} [length] Message length if known beforehand
             * @returns {google.protobuf.Timestamp} Timestamp
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            Timestamp.decode = function decode(reader, length) {
                if (!(reader instanceof $Reader))
                    reader = $Reader.create(reader);
                let end = length === undefined ? reader.len : reader.pos + length, message = new $root.google.protobuf.Timestamp();
                while (reader.pos < end) {
                    let tag = reader.uint32();
                    switch (tag >>> 3) {
                    case 1: {
                            message.seconds = reader.int64();
                            break;
                        }
                    case 2: {
                            message.nanos = reader.int32();
                            break;
                        }
                    default:
                        reader.skipType(tag & 7);
                        break;
                    }
                }
                return message;
            };

            /**
             * Decodes a Timestamp message from the specified reader or buffer, length delimited.
             * @function decodeDelimited
             * @memberof google.protobuf.Timestamp
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @returns {google.protobuf.Timestamp} Timestamp
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            Timestamp.decodeDelimited = function decodeDelimited(reader) {
                if (!(reader instanceof $Reader))
                    reader = new $Reader(reader);
                return this.decode(reader, reader.uint32());
            };

            /**
             * Verifies a Timestamp message.
             * @function verify
             * @memberof google.protobuf.Timestamp
             * @static
             * @param {Object.<string,*>} message Plain object to verify
             * @returns {string|null} `null` if valid, otherwise the reason why it is not
             */
            Timestamp.verify = function verify(message) {
                if (typeof message !== "object" || message === null)
                    return "object expected";
                if (message.seconds != null && message.hasOwnProperty("seconds"))
                    if (!$util.isInteger(message.seconds) && !(message.seconds && $util.isInteger(message.seconds.low) && $util.isInteger(message.seconds.high)))
                        return "seconds: integer|Long expected";
                if (message.nanos != null && message.hasOwnProperty("nanos"))
                    if (!$util.isInteger(message.nanos))
                        return "nanos: integer expected";
                return null;
            };

            /**
             * Creates a Timestamp message from a plain object. Also converts values to their respective internal types.
             * @function fromObject
             * @memberof google.protobuf.Timestamp
             * @static
             * @param {Object.<string,*>} object Plain object
             * @returns {google.protobuf.Timestamp} Timestamp
             */
            Timestamp.fromObject = function fromObject(object) {
                if (object instanceof $root.google.protobuf.Timestamp)
                    return object;
                let message = new $root.google.protobuf.Timestamp();
                if (object.seconds != null)
                    if ($util.Long)
                        (message.seconds = $util.Long.fromValue(object.seconds)).unsigned = false;
                    else if (typeof object.seconds === "string")
                        message.seconds = parseInt(object.seconds, 10);
                    else if (typeof object.seconds === "number")
                        message.seconds = object.seconds;
                    else if (typeof object.seconds === "object")
                        message.seconds = new $util.LongBits(object.seconds.low >>> 0, object.seconds.high >>> 0).toNumber();
                if (object.nanos != null)
                    message.nanos = object.nanos | 0;
                return message;
            };

            /**
             * Creates a plain object from a Timestamp message. Also converts values to other types if specified.
             * @function toObject
             * @memberof google.protobuf.Timestamp
             * @static
             * @param {google.protobuf.Timestamp} message Timestamp
             * @param {$protobuf.IConversionOptions} [options] Conversion options
             * @returns {Object.<string,*>} Plain object
             */
            Timestamp.toObject = function toObject(message, options) {
                if (!options)
                    options = {};
                let object = {};
                if (options.defaults) {
                    if ($util.Long) {
                        let long = new $util.Long(0, 0, false);
                        object.seconds = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
                    } else
                        object.seconds = options.longs === String ? "0" : 0;
                    object.nanos = 0;
                }
                if (message.seconds != null && message.hasOwnProperty("seconds"))
                    if (typeof message.seconds === "number")
                        object.seconds = options.longs === String ? String(message.seconds) : message.seconds;
                    else
                        object.seconds = options.longs === String ? $util.Long.prototype.toString.call(message.seconds) : options.longs === Number ? new $util.LongBits(message.seconds.low >>> 0, message.seconds.high >>> 0).toNumber() : message.seconds;
                if (message.nanos != null && message.hasOwnProperty("nanos"))
                    object.nanos = message.nanos;
                return object;
            };

            /**
             * Converts this Timestamp to JSON.
             * @function toJSON
             * @memberof google.protobuf.Timestamp
             * @instance
             * @returns {Object.<string,*>} JSON object
             */
            Timestamp.prototype.toJSON = function toJSON() {
                return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
            };

            /**
             * Gets the default type url for Timestamp
             * @function getTypeUrl
             * @memberof google.protobuf.Timestamp
             * @static
             * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
             * @returns {string} The default type url
             */
            Timestamp.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
                if (typeUrlPrefix === undefined) {
                    typeUrlPrefix = "type.googleapis.com";
                }
                return typeUrlPrefix + "/google.protobuf.Timestamp";
            };

            return Timestamp;
        })();

        return protobuf;
    })();

    return google;
})();

export { $root as default };
