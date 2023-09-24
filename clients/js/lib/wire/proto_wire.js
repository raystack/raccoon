import WireType from '../types/wire_type.js';

function createProtoMarshaller() {
    function marshal(data) {
        if (!data) {
            throw new Error('Invalid data object for marshalling');
        }

        if (
            !data.constructor ||
            !data.constructor.encode ||
            typeof data.constructor.encode !== 'function'
        ) {
            throw new Error('Invalid Protobuf message object for marshalling');
        }

        return new Uint8Array(data.constructor.encode(data).finish());
    }

    function unmarshal(data, type) {
        if (!data) {
            throw new Error('Invalid data for unmarshalling');
        }

        if (!type || !type.decode || typeof type.decode !== 'function') {
            throw new Error('Invalid Protobuf message type for unmarshalling');
        }

        return type.decode(new Uint8Array(data));
    }

    function getContentType() {
        return WireType.PROTOBUF;
    }

    function getResponseType() {
        return 'arraybuffer';
    }

    return {
        marshal,
        unmarshal,
        getContentType,
        getResponseType
    };
}

export default createProtoMarshaller;
