class Wire:
    def marshal(self, obj):
        raise NotImplementedError("not implemented")

    def unmarshal(self, obj, template):
        raise NotImplementedError("not implemented")

    def get_content_type(self):
        raise NotImplementedError("not implemented")
