# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: raystack/raccoon/v1beta1/raccoon.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import timestamp_pb2 as google_dot_protobuf_dot_timestamp__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n&raystack/raccoon/v1beta1/raccoon.proto\x12\x18raystack.raccoon.v1beta1\x1a\x1fgoogle/protobuf/timestamp.proto\"\x9f\x01\n\x10SendEventRequest\x12\x19\n\x08req_guid\x18\x01 \x01(\tR\x07reqGuid\x12\x37\n\tsent_time\x18\x02 \x01(\x0b\x32\x1a.google.protobuf.TimestampR\x08sentTime\x12\x37\n\x06\x65vents\x18\x03 \x03(\x0b\x32\x1f.raystack.raccoon.v1beta1.EventR\x06\x65vents\"<\n\x05\x45vent\x12\x1f\n\x0b\x65vent_bytes\x18\x01 \x01(\x0cR\neventBytes\x12\x12\n\x04type\x18\x02 \x01(\tR\x04type\"\xba\x02\n\x11SendEventResponse\x12\x38\n\x06status\x18\x01 \x01(\x0e\x32 .raystack.raccoon.v1beta1.StatusR\x06status\x12\x32\n\x04\x63ode\x18\x02 \x01(\x0e\x32\x1e.raystack.raccoon.v1beta1.CodeR\x04\x63ode\x12\x1b\n\tsent_time\x18\x03 \x01(\x03R\x08sentTime\x12\x16\n\x06reason\x18\x04 \x01(\tR\x06reason\x12I\n\x04\x64\x61ta\x18\x05 \x03(\x0b\x32\x35.raystack.raccoon.v1beta1.SendEventResponse.DataEntryR\x04\x64\x61ta\x1a\x37\n\tDataEntry\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n\x05value\x18\x02 \x01(\tR\x05value:\x02\x38\x01*F\n\x06Status\x12\x16\n\x12STATUS_UNSPECIFIED\x10\x00\x12\x12\n\x0eSTATUS_SUCCESS\x10\x01\x12\x10\n\x0cSTATUS_ERROR\x10\x02*\xa0\x01\n\x04\x43ode\x12\x14\n\x10\x43ODE_UNSPECIFIED\x10\x00\x12\x0b\n\x07\x43ODE_OK\x10\x01\x12\x14\n\x10\x43ODE_BAD_REQUEST\x10\x02\x12\x17\n\x13\x43ODE_INTERNAL_ERROR\x10\x03\x12%\n!CODE_MAX_CONNECTION_LIMIT_REACHED\x10\x04\x12\x1f\n\x1b\x43ODE_MAX_USER_LIMIT_REACHED\x10\x05\x32t\n\x0c\x45ventService\x12\x64\n\tSendEvent\x12*.raystack.raccoon.v1beta1.SendEventRequest\x1a+.raystack.raccoon.v1beta1.SendEventResponseB[\n\x1aio.raystack.proton.raccoonB\nEventProtoP\x01Z/github.com/raystack/proton/raccoon/v1;raccoonv1b\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'raystack.raccoon.v1beta1.raccoon_pb2', _globals)
if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'\n\032io.raystack.proton.raccoonB\nEventProtoP\001Z/github.com/raystack/proton/raccoon/v1;raccoonv1'
  _SENDEVENTRESPONSE_DATAENTRY._options = None
  _SENDEVENTRESPONSE_DATAENTRY._serialized_options = b'8\001'
  _globals['_STATUS']._serialized_start=642
  _globals['_STATUS']._serialized_end=712
  _globals['_CODE']._serialized_start=715
  _globals['_CODE']._serialized_end=875
  _globals['_SENDEVENTREQUEST']._serialized_start=102
  _globals['_SENDEVENTREQUEST']._serialized_end=261
  _globals['_EVENT']._serialized_start=263
  _globals['_EVENT']._serialized_end=323
  _globals['_SENDEVENTRESPONSE']._serialized_start=326
  _globals['_SENDEVENTRESPONSE']._serialized_end=640
  _globals['_SENDEVENTRESPONSE_DATAENTRY']._serialized_start=585
  _globals['_SENDEVENTRESPONSE_DATAENTRY']._serialized_end=640
  _globals['_EVENTSERVICE']._serialized_start=877
  _globals['_EVENTSERVICE']._serialized_end=993
# @@protoc_insertion_point(module_scope)
