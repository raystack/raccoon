export function createJsonSerializer() {
    return function serialize(data) {
        if (!data) {
            throw new Error("Invalid data object for serialization");
        }

        const jsonString = JSON.stringify(data);

        return Array.from(Buffer.from(jsonString));
    }
}
