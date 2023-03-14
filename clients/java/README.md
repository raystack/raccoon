# Java Client for Raccoon


## Requirements

- [Gradle v7+](https://gradle.org/)
- [JDK 8+](https://openjdk.java.net/projects/jdk8/)


### Add raccoon as dependency

#### Gradle

```groovy
  implementation group: 'com.gotocompany', name: 'raccoon', version: '0.1.5'
```

#### Maven

```xml
<dependency>
  <groupId>com.gotocompany</groupId>
  <artifactId>racoon</artifactId>
  <version>0.1.5</version>
</dependency>
```

## Usage

#### Construct a new REST client, then use the various options to build the client.
For example:

```java
  // build the configuration.
  RestConfig config = RestConfig.builder()
                  .url("http://localhost:8080/api/v1/events")
                  .header("x-connection-id", "123")
                  .serializer(new ProtoSerializer()) // default is Json
                  .marshaler(new ProtoWire()) // default is Json
                  .retryMax(5) // default is 3
                  .retryWait(2000) // default is one second
                  .build();

  // get the rest client instance.
  RaccoonClient restClient = RaccoonClientFactory.getRestClient(config);

  // prepare the event to be send.
  PageEventProto.PageEvent pageEvent = PageEventProto.PageEvent.newBuilder()
                  .setEventGuid(UUID.randomUUID().toString())
                  .setEventName("clicked")
                  .setSentTime(Timestamp.newBuilder().build())
                  .build();

  // send the event batch to the raccoon.
  Response response = restClient.send(
                  new Event[] {
                                  new Event("page", pageEvent)
                  });

  // check the status of the request.
  if (response.isSuccess() && response.getStatus() == ResponseStatus.STATUS_SUCCESS) {
          System.out.println("The event sent successfully");
  }
```