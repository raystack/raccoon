class Wire:
    def marshal(self, event):
        raise NotImplementedError("not implemented")

    def unmarshal(self, data, template):
        raise NotImplementedError("not implemented")

    def get_content_type(self):
        raise NotImplementedError("not implemented")
