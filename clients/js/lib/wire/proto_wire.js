export function createProtoMarshaller() {
    function marshal(data) {
        if (!data) {
            throw new Error("Invalid data object for marshalling");
        }

        if (!data.constructor || !data.constructor.encode || typeof data.constructor.encode !== "function") {
            throw new Error("Invalid Protobuf message object for marshalling");
        }

        return new Uint8Array(data.constructor.encode(data).finish());
    }

    function unmarshal(data, type) {
        if (!data) {
            throw new Error("Invalid data for unmarshalling");
        }

        if (!type || !type.decode || typeof type.decode !== "function") {
            throw new Error("Invalid Protobuf message type for unmarshalling");
        }

        return type.decode(new TextEncoder().encode((data)));
    }

    function getContentType() {
        return "application/proto";
    }

    return {
        marshal,
        unmarshal,
        getContentType
    };
}
