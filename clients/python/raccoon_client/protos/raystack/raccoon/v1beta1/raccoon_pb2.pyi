from google.protobuf import timestamp_pb2 as _timestamp_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Status(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = []
    STATUS_UNSPECIFIED: _ClassVar[Status]
    STATUS_SUCCESS: _ClassVar[Status]
    STATUS_ERROR: _ClassVar[Status]

class Code(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = []
    CODE_UNSPECIFIED: _ClassVar[Code]
    CODE_OK: _ClassVar[Code]
    CODE_BAD_REQUEST: _ClassVar[Code]
    CODE_INTERNAL_ERROR: _ClassVar[Code]
    CODE_MAX_CONNECTION_LIMIT_REACHED: _ClassVar[Code]
    CODE_MAX_USER_LIMIT_REACHED: _ClassVar[Code]
STATUS_UNSPECIFIED: Status
STATUS_SUCCESS: Status
STATUS_ERROR: Status
CODE_UNSPECIFIED: Code
CODE_OK: Code
CODE_BAD_REQUEST: Code
CODE_INTERNAL_ERROR: Code
CODE_MAX_CONNECTION_LIMIT_REACHED: Code
CODE_MAX_USER_LIMIT_REACHED: Code

class SendEventRequest(_message.Message):
    __slots__ = ["req_guid", "sent_time", "events"]
    REQ_GUID_FIELD_NUMBER: _ClassVar[int]
    SENT_TIME_FIELD_NUMBER: _ClassVar[int]
    EVENTS_FIELD_NUMBER: _ClassVar[int]
    req_guid: str
    sent_time: _timestamp_pb2.Timestamp
    events: _containers.RepeatedCompositeFieldContainer[Event]
    def __init__(self, req_guid: _Optional[str] = ..., sent_time: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., events: _Optional[_Iterable[_Union[Event, _Mapping]]] = ...) -> None: ...

class Event(_message.Message):
    __slots__ = ["event_bytes", "type"]
    EVENT_BYTES_FIELD_NUMBER: _ClassVar[int]
    TYPE_FIELD_NUMBER: _ClassVar[int]
    event_bytes: bytes
    type: str
    def __init__(self, event_bytes: _Optional[bytes] = ..., type: _Optional[str] = ...) -> None: ...

class SendEventResponse(_message.Message):
    __slots__ = ["status", "code", "sent_time", "reason", "data"]
    class DataEntry(_message.Message):
        __slots__ = ["key", "value"]
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: str
        def __init__(self, key: _Optional[str] = ..., value: _Optional[str] = ...) -> None: ...
    STATUS_FIELD_NUMBER: _ClassVar[int]
    CODE_FIELD_NUMBER: _ClassVar[int]
    SENT_TIME_FIELD_NUMBER: _ClassVar[int]
    REASON_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    status: Status
    code: Code
    sent_time: int
    reason: str
    data: _containers.ScalarMap[str, str]
    def __init__(self, status: _Optional[_Union[Status, str]] = ..., code: _Optional[_Union[Code, str]] = ..., sent_time: _Optional[int] = ..., reason: _Optional[str] = ..., data: _Optional[_Mapping[str, str]] = ...) -> None: ...
