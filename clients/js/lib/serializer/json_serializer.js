function createJsonSerializer() {
    return function serialize(data) {

        const jsonString = JSON.stringify(data);

        return Array.from(Buffer.from(jsonString));
    }
}

module.exports = createJsonSerializer;
