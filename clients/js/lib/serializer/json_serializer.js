export function createJsonSerializer() {
    return function serialize(data) {
        if (!data) {
            throw new Error("Invalid data object for serialization");
        }

        const jsonString = JSON.stringify(data);

        const encoder = new TextEncoder();

        return Array.from(encoder.encode(jsonString));
    }
}
