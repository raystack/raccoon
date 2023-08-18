class Event:
    type: str
    event: object


class Client:

    def send(self, events: [Event]):
        raise NotImplementedError()


