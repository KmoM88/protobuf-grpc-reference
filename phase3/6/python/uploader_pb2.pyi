from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class ImageChunk(_message.Message):
    __slots__ = ("data",)
    DATA_FIELD_NUMBER: _ClassVar[int]
    data: bytes
    def __init__(self, data: _Optional[bytes] = ...) -> None: ...

class UploadSummary(_message.Message):
    __slots__ = ("total_chunks", "total_bytes")
    TOTAL_CHUNKS_FIELD_NUMBER: _ClassVar[int]
    TOTAL_BYTES_FIELD_NUMBER: _ClassVar[int]
    total_chunks: int
    total_bytes: int
    def __init__(self, total_chunks: _Optional[int] = ..., total_bytes: _Optional[int] = ...) -> None: ...
