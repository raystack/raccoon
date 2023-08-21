import unittest

from raccoon_client.serde import util
from raccoon_client.serde.enum import Serialiser, WireType
from raccoon_client.serde.json_serde import JsonSerde
from raccoon_client.serde.protobuf_serde import ProtobufSerde


class UtilTest(unittest.TestCase):
    def test_serde_factory(self):
        serde = util.get_serde(Serialiser.JSON)
        self.assertIsInstance(serde, JsonSerde)

    def test_wire_factory(self):
        wire = util.get_wire_type(WireType.JSON)
        self.assertIsInstance(wire, JsonSerde)

    def test_serde_factory_for_proto(self):
        serde = util.get_serde(Serialiser.PROTOBUF)
        self.assertIsInstance(serde, ProtobufSerde)

    def test_wire_factory_for_proto(self):
        wire = util.get_wire_type(WireType.PROTOBUF)
        self.assertIsInstance(wire, ProtobufSerde)

    def test_invalid_wire_factory_input(self):
        self.assertRaises(ValueError, util.get_wire_type, "invalid_type")

    def test_invalid_serde_factory_input(self):
        self.assertRaises(ValueError, util.get_serde, "invalid_type")
