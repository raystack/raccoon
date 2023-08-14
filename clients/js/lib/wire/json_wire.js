export function createJsonMarshaller() {
    function marshal(data) {
        if (!data) {
            throw new Error("Invalid data object for marshalling");
        }

        if (!data.constructor || !data.constructor.toObject || typeof data.constructor.toObject !== "function") {
            throw new Error("Invalid Protobuf message object for marshalling");
        }

        const dataObject = data.constructor.toObject(data)
        return JSON.stringify(dataObject);
    }

    function unmarshal(data, type) {
        if (!data) {
            throw new Error("Invalid data object for unmarshalling");
        }

        if (!type || typeof type.fromObject !== "function") {
            throw new Error("Invalid Protobuf message type for unmarshalling");
        }

        return type.fromObject(data);
    }

    function getContentType() {
        return "application/json";
    }

    return {
        marshal,
        unmarshal,
        getContentType
    };
}
