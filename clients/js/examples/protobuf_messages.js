import { RaccoonClient, SerializationType, WireType } from '@raystack/raccoon';
//import the compiled js file generated via protobufjs
import { google, clickevents } from './protos/compiled.js';

const currentTime = new Date();
const timestamp = google.protobuf.Timestamp.create({
    seconds: Math.floor(currentTime / 1000),
    nanos: (currentTime % 1000) * 1e6
});

//create the protobufjs messages and set the field values
const pageEvent = new clickevents.PageEvent();
pageEvent.eventGuid = "123";
pageEvent.eventName = "page open";
pageEvent.sentTime = timestamp;

const clickEvent = new clickevents.ClickEvent();
clickEvent.eventGuid = "123";
clickEvent.componentIndex = 12;
clickEvent.componentName = "images";
clickEvent.sentTime = timestamp;

const protobufEvents = [
    {
        type: 'test-topic1',
        data: clickEvent
    },
    {
        type: 'test-topic2',
        data: pageEvent
    }
];

//initialise the raccoon client with required configs
const raccoonClient = new RaccoonClient({
    serializationType: SerializationType.PROTOBUF,
    wireType: WireType.JSON,
    timeout: 5000,
    url: 'http://localhost:8080/api/v1/events',
    headers: {
        'X-User-ID': 'user-1'
    }
});

//send the request
raccoonClient.send(protobufEvents)
    .then(result => {
        console.log('Result:', result);
    })
    .catch(error => {
        console.error('Error:', error);
    });
