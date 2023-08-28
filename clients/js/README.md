# Raccoon Client for Javascript

A JS client library for [Raccoon](https://github.com/raystack/raccoon), compatible both in browsers as well as nodejs environment.

Features:
* Send JSON/Protobuf encoded messages to raccoon server
* Configurable request wiretype i.e.,JSON/PROTOBUF
* Retry mechanism

**Note:** For encoding protobuf messages, this client utilises [protobufjs](https://github.com/protobufjs/protobuf.js).

## Requirements

- **Node.js**: Version 20.x or above.
- **Browser**: Modern web browser with ES6+ compatibility.

## Install

```bash
npm i @raystack/raccoon --save
```

## Usage

#### Construct a new REST client, then pass the available options on the client.

For example:

```javascript
import { RaccoonClient, SerializationType, WireType } from '@raystack/raccoon';
```

```javascript
const raccoonClient = new RaccoonClient({
    serializationType: SerializationType.PROTOBUF,
    wireType: WireType.JSON,
    timeout: 5000,
    url: 'http://localhost:8080/api/v1/events',
    headers: {
        'X-User-ID': 'user-1'
    }
});
```

#### Send messages

In the below example, we are sending protobuf messages. 

**Note:** In this case, we only need to send the protobufjs object along with values, the client will take care of encoding the object.

```javascript
const pageEvent = new clickevents.PageEvent();
pageEvent.eventGuid = "123";
pageEvent.eventName = "page open";

const clickEvent = new clickevents.ClickEvent();
clickEvent.eventGuid = "123";
clickEvent.componentIndex = 1;
clickEvent.componentName = "images";

const events = [
    {
        type: 'test-topic1',
        data: clickEvent
    },
    {
        type: 'test-topic2',
        data: pageEvent
    }
];

raccoonClient.send(events)
    .then(result => {
        console.log('Result:', result);
    })
    .catch(error => {
        console.error('Error:', error);
    });
```

**Note:** Here, the PageEvent and ClickEvent are imported from generated protobufjs code.

See complete examples for sending JSON and protobuf messages: [examples](examples)

## Configurations

#### `SerializationType`

The serialization type to use for the events being sent to Kafka, either 'PROTOBUF' or 'JSON'.

- Type: `Required`
- Example value: `SerializationType.JSON`

#### `WireType`

The wire configuration using which the request payload should be sent to raccoon server

- Type: `Required`
- Example value: `WireType.JSON`

#### `url`

The URL for the raccoon server.

- Type: `Required`
- Example value: `http://localhost:8080/api/v1/events`

#### `headers`

Custom headers to be included in the HTTP requests.

- Type: `Optional`
- Default value: `{}`

#### `timeout`

The request timeout in milliseconds.

- Type: `Optional`
- Default value: `1000`

#### `retryMax`

The maximum number of retry attempts for failed requests.

- Type: `Optional`
- Default value: `3`

#### `retryWait`

The time in milliseconds to wait between retry attempts.

- Type: `Optional`
- Default value: `5000`

#### `logger`

Logger object for logging.

- Type: `Optional`
- Default value: `console`

## Setting up development environment

### Prerequisite Tools

- [Node.js](https://nodejs.org/) (version >= 20.x)
- [Git](https://git-scm.com/)

1. Clone the repo

   ```sh
   $ git clone https://github.com/raystack/raccoon
   $ cd raccoon/clients/js
   ```

2. Install dependencies

   ```sh
   $ npm install
   ```

3. Run the tests. All of the tests are written with [jest](https://jestjs.io/).

   ```sh
   $ npm test
   ```
4. Run format.

   ```sh
   $ npm run format
   ```
4. Run lint.

   ```sh
   $ npm run lint
   ```

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags](https://www.npmjs.com/package/@raystack/raccoon?activeTab=versions).
