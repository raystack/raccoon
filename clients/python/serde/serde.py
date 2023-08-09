class Serde:
    def serialise(self, event):
        raise NotImplementedError()

    def deserialise(self, data, template):
        raise NotImplementedError()

    def get_content_type(self):
        raise NotImplementedError()
