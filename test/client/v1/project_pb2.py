# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: v1/project.proto
"""Generated protocol buffer code."""
from google.protobuf.internal import builder as _builder
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x10v1/project.proto\x12\x06\x61pi.v1\"\x1b\n\x0bPingRequest\x12\x0c\n\x04name\x18\x01 \x01(\t\"\x1c\n\tPingReply\x12\x0f\n\x07message\x18\x01 \x01(\t\",\n\x0eProjectRequest\x12\x0c\n\x04name\x18\x01 \x01(\t\x12\x0c\n\x04\x63ode\x18\x02 \x01(\t\"\x1e\n\rProjectFilter\x12\r\n\x05regex\x18\x01 \x01(\t\"1\n\x07Project\x12\n\n\x02id\x18\x01 \x01(\t\x12\x0c\n\x04name\x18\x02 \x01(\t\x12\x0c\n\x04\x63ode\x18\x03 \x01(\t2\xae\x01\n\x06Studio\x12\x30\n\x04Ping\x12\x13.api.v1.PingRequest\x1a\x11.api.v1.PingReply\"\x00\x12:\n\rCreateProject\x12\x16.api.v1.ProjectRequest\x1a\x0f.api.v1.Project\"\x00\x12\x36\n\x08Projects\x12\x15.api.v1.ProjectFilter\x1a\x0f.api.v1.Project\"\x00\x30\x01\x42)Z\'github.com/studio1767/studio-api/api_v1b\x06proto3')

_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, globals())
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'v1.project_pb2', globals())
if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'Z\'github.com/studio1767/studio-api/api_v1'
  _PINGREQUEST._serialized_start=28
  _PINGREQUEST._serialized_end=55
  _PINGREPLY._serialized_start=57
  _PINGREPLY._serialized_end=85
  _PROJECTREQUEST._serialized_start=87
  _PROJECTREQUEST._serialized_end=131
  _PROJECTFILTER._serialized_start=133
  _PROJECTFILTER._serialized_end=163
  _PROJECT._serialized_start=165
  _PROJECT._serialized_end=214
  _STUDIO._serialized_start=217
  _STUDIO._serialized_end=391
# @@protoc_insertion_point(module_scope)
