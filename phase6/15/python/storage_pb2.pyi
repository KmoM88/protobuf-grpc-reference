from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class FileMetadata(_message.Message):
    __slots__ = ("file_id", "filename", "size_bytes", "storage_node_address", "auth_token")
    FILE_ID_FIELD_NUMBER: _ClassVar[int]
    FILENAME_FIELD_NUMBER: _ClassVar[int]
    SIZE_BYTES_FIELD_NUMBER: _ClassVar[int]
    STORAGE_NODE_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    AUTH_TOKEN_FIELD_NUMBER: _ClassVar[int]
    file_id: str
    filename: str
    size_bytes: int
    storage_node_address: str
    auth_token: str
    def __init__(self, file_id: _Optional[str] = ..., filename: _Optional[str] = ..., size_bytes: _Optional[int] = ..., storage_node_address: _Optional[str] = ..., auth_token: _Optional[str] = ...) -> None: ...

class FileChunk(_message.Message):
    __slots__ = ("file_id", "chunk_index", "data", "offset")
    FILE_ID_FIELD_NUMBER: _ClassVar[int]
    CHUNK_INDEX_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    OFFSET_FIELD_NUMBER: _ClassVar[int]
    file_id: str
    chunk_index: int
    data: bytes
    offset: int
    def __init__(self, file_id: _Optional[str] = ..., chunk_index: _Optional[int] = ..., data: _Optional[bytes] = ..., offset: _Optional[int] = ...) -> None: ...

class UploadStatus(_message.Message):
    __slots__ = ("file_id", "bytes_received", "success", "message")
    FILE_ID_FIELD_NUMBER: _ClassVar[int]
    BYTES_RECEIVED_FIELD_NUMBER: _ClassVar[int]
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    file_id: str
    bytes_received: int
    success: bool
    message: str
    def __init__(self, file_id: _Optional[str] = ..., bytes_received: _Optional[int] = ..., success: bool = ..., message: _Optional[str] = ...) -> None: ...

class Empty(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...
