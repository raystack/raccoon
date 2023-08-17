import unittest

from serde import util
from serde.enum import Serialiser, WireType
from serde.json_serde import JsonSerde


class UtilTest(unittest.TestCase):
    def test_serde_factory(self):
        serde = util.get_serde(Serialiser.JSON)
        self.assertIsInstance(serde, JsonSerde)

    def test_wire_factory(self):
        wire = util.get_wire_type(WireType.JSON)
        self.assertIsInstance(wire, JsonSerde)
