const createJsonSerializer = require('../../lib/serializer/json_serializer');

describe('JsonSerializer', () => {
    test('should convert data to JSON', () => {
        const jsonSerializer = createJsonSerializer();
        const data = { key: 'value' };
        const serializedData = jsonSerializer.serialize(data);

        expect(JSON.parse(Buffer.from(serializedData).toString())).toEqual(data);
    });

    test('should handle empty object', () => {
        const jsonSerializer = createJsonSerializer();
        const data = {};
        const serializedData = jsonSerializer.serialize(data);

        expect(JSON.parse(Buffer.from(serializedData).toString())).toEqual(data);
    });

    test('should handle empty array', () => {
        const jsonSerializer = createJsonSerializer();
        const data = [];
        const serializedData = jsonSerializer.serialize(data);

        expect(JSON.parse(Buffer.from(serializedData).toString())).toEqual(data);
    });

    test('should handle null', () => {
        const jsonSerializer = createJsonSerializer();
        const data = null;
        const serializedData = jsonSerializer.serialize(data);

        expect(JSON.parse(Buffer.from(serializedData).toString())).toEqual(data);
    });

    test('should handle undefined with proper error', () => {
        const jsonSerializer = createJsonSerializer();
        const data = undefined;

        expect(() => {
            jsonSerializer.serialize(data);
        }).toThrowError("The first argument must be of type string or an instance of Buffer, ArrayBuffer, or Array or an Array-like Object. Received undefined");
    });


    test('should handle numbers', () => {
        const jsonSerializer = createJsonSerializer();
        const data = 12345;
        const serializedData = jsonSerializer.serialize(data);

        expect(JSON.parse(Buffer.from(serializedData).toString())).toEqual(data);
    });

    test('should handle strings', () => {
        const jsonSerializer = createJsonSerializer();
        const data = 'hello';
        const serializedData = jsonSerializer.serialize(data);

        expect(JSON.parse(Buffer.from(serializedData).toString())).toEqual(data);
    });
});
