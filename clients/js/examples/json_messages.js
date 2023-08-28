// eslint-disable-next-line import/no-unresolved
import { RaccoonClient, SerializationType, WireType } from '@raystack/raccoon';

const logger = console;

//  create json messages
const jsonEvents = [
    {
        type: 'test-topic1',
        data: { key1: 'value1', key2: ['a', 'b'] }
    },
    {
        type: 'test-topic2',
        data: { key1: 'value2', key2: { key3: 'value3', key4: 'value4' } }
    }
];

//  initialise the raccoon client with required configs
const raccoonClient = new RaccoonClient({
    serializationType: SerializationType.JSON,
    wireType: WireType.JSON,
    timeout: 5000,
    url: 'http://localhost:8080/api/v1/events',
    headers: {
        'X-User-ID': 'user-1'
    }
});

//  send the request
raccoonClient
    .send(jsonEvents)
    .then((result) => {
        logger.log('Result:', result);
    })
    .catch((error) => {
        logger.error('Error:', error);
    });
