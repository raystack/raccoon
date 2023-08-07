function createJsonSerializer() {
    function serialize(data) {

        const jsonString = JSON.stringify(data);

        return Array.from(Buffer.from(jsonString));
    }

    return { serialize };
}

module.exports = createJsonSerializer;
