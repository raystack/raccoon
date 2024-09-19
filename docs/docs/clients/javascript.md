# Javascript

## Requirements
Make sure that Nodejs >= `20.0` is installed on your system. See [installation instructions](https://nodejs.org/en/download/package-manager) on Nodejs's website for more info.

## Installation
Install Raccoon's Javascript client using npm
```javascript
$ npm install --save @raystack/raccoon
```
## Usage

### Quickstart

Below is a self contained example of Raccoon's Javascript client that uses Raccoon's REST API to publish events.

```javascript title="quickstart.js" showLineNumbers
import { RaccoonClient, SerializationType, WireType } from '@raystack/raccoon';

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

const raccoonClient = new RaccoonClient({
    serializationType: SerializationType.JSON,
    wireType: WireType.JSON,
    timeout: 5000,
    url: 'http://localhost:8080/api/v1/events',
    headers: {
        'X-User-ID': 'user-1'
    }
});

raccoonClient
    .send(jsonEvents)
    .then((result) => {
        console.log('Result:', result);
    })
    .catch((error) => {
        console.error('Error:', error);
    });
```

### Guide

#### Creating a client
Raccoon's Javascript client only supports sending event's over Raccoon's HTTP/JSON (REST) API.

To create the client, use `new RaccoonClient(options)`. `options` is javascript object that contains the following properites:

| Property | Description |
| --- | --- |
| url | (required) The base URL for the API requests |
| serializationType |  (required) The serialization type to use, either 'protobuf' or 'json' |
| wireType | The wire configuration, containing ContentType (default: `wireType.JSON`)|
| headers | Custom headers to be included in the HTTP requests (default: `{}`)|
| retryMax | The maximum number of retry attempts for failed requests (default: `3`) |
| retryWait | The time in milliseconds to wait between retry attempts (default: `1000`)| 
| timeout | The timeout in milliseconds (default: `1000`)|
| logger | Logger object for logging (default: `global.console`)

#### Publishing events
To publish events, create an array of objects and pass it to `RaccoonClient#send()`. The return value is a [Promise](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise).

```js
const events = [
    {
        type: "event_type",
        data: {},
    }
];

client.send(events)
    .then(result => console.log(result))
    .catch(err => console.error(err))
```

`type` denotes the event type. This is used by raccoon to route the event to a specific topic downstream. `data` field contains the payload. This data is serialised by the `serializerType` that's configured on the client. 

The following table lists which serializer to use for a given payload type.

| Message Type | Serializer |
| --- | --- |
| JSON | `SerializationType.JSON` |
| Protobuf | `SerializationType.PROTOBUF`|

Once a client is constructed with a specific kind of serializer, you may only pass it events of that specific type. In particular, for `JSON` serialiser the event data must be a javascript object. While for `PROTOBUF` serialiser the event data must be a protobuf message.

## Examples
You can find examples of client usage [here](https://github.com/raystack/raccoon/tree/main/clients/js/examples)