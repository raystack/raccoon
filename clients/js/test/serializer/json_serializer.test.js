const createJsonSerializer = require('../../lib/serializer/json_serializer');

describe('JsonSerializer', () => {
    test('should convert data to JSON', () => {
        const serialize = createJsonSerializer();
        const data = { key: 'value' };
        const serializedData = serialize(data);

        expect(JSON.parse(Buffer.from(serializedData).toString())).toEqual(data);
    });

    test('should handle empty object', () => {
        const serialize = createJsonSerializer();
        const data = {};
        const serializedData = serialize(data);

        expect(JSON.parse(Buffer.from(serializedData).toString())).toEqual(data);
    });

    test('should handle empty array', () => {
        const serialize = createJsonSerializer();
        const data = [];
        const serializedData = serialize(data);

        expect(JSON.parse(Buffer.from(serializedData).toString())).toEqual(data);
    });

    test('should handle null', () => {
        const serialize = createJsonSerializer();
        const data = null;
        const serializedData = serialize(data);

        expect(JSON.parse(Buffer.from(serializedData).toString())).toEqual(data);
    });

    test('should handle undefined with proper error', () => {
        const serialize = createJsonSerializer();
        const data = undefined;

        expect(() => {
            serialize(data);
        }).toThrowError("The first argument must be of type string or an instance of Buffer, ArrayBuffer, or Array or an Array-like Object. Received undefined");
    });


    test('should handle numbers', () => {
        const serialize = createJsonSerializer();
        const data = 12345;
        const serializedData = serialize(data);

        expect(JSON.parse(Buffer.from(serializedData).toString())).toEqual(data);
    });

    test('should handle strings', () => {
        const serialize = createJsonSerializer();
        const data = 'hello';
        const serializedData = serialize(data);

        expect(JSON.parse(Buffer.from(serializedData).toString())).toEqual(data);
    });
});
