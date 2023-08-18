import unittest

from raccoon_client.serde import util
from raccoon_client.serde.enum import Serialiser, WireType
from raccoon_client.serde.json_serde import JsonSerde


class UtilTest(unittest.TestCase):
    def test_serde_factory(self):
        serde = util.get_serde(Serialiser.JSON)
        self.assertIsInstance(serde, JsonSerde)

    def test_wire_factory(self):
        wire = util.get_wire_type(WireType.JSON)
        self.assertIsInstance(wire, JsonSerde)
