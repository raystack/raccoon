function createProtobufSerializer() {
    return function serialize(data) {
        if (!data) {
            throw new Error('Invalid data object for serialization');
        }

        if (
            !data.constructor ||
            !data.constructor.encode ||
            typeof data.constructor.encode !== 'function'
        ) {
            throw new Error('Invalid Protobuf message object for serialization');
        }

        return Array.from(data.constructor.encode(data).finish());
    };
}

export default createProtobufSerializer;
