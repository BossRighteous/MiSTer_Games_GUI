
import os
import struct
from typing import Dict, Union

"""
https://www.kernel.org/doc/Documentation/input/joystick-api.txt
"""

class JSEvent():
    EVENT_SIZE: int = 8
    FORMAT: str = '<LhBB'
    AXIS_MAX = 32767
    AXIS_DEADZONE = 8191 # 1/4 MAX

    timestamp: int
    value: int
    type: int
    number: int
    changed: bool

    def __init__(
        self,
        timestamp: int,
        value: int,
        type: int, # 1 or 2 if not using
        number: int,
    ) -> None:
        self.timestamp = timestamp
        self.value = value
        self.type = type
        self.number = number
        self.changed = False
    
    def normalize_value(self) -> int:
        if self.type == 2 and self.value != 0:
            absval = abs(self.value)
            if absval < self.AXIS_DEADZONE:
                self.value = 0
            else:
                self.value = self.value / absval
    
    def to_bytes(self) -> bytes:
        return struct.pack(self.FORMAT, self.timestamp, self.value, self.type, self.number)
    
    def is_axis(self) -> bool:
        return self.type == 2
    
    def map_index(self) -> bytes:            
        return ("jsb" if self.type == 1 else "jsa") + str(self.number)
    
    def is_just_pressed(self) -> bool:
        return self.changed == True and self.value == 1
    
    def is_just_released(self) -> bool:
        return self.changed == True and self.value == 0
    
    def is_pressed(self) -> bool:
        return self.value != 0


def JSEvent_from_bytes(buffer: bytes) -> Union[JSEvent, None]:
        (timestamp, value, type, number) = struct.unpack(JSEvent.FORMAT, buffer)
        type = type &~ 128
        if type != 1 and type != 2:
            return None
        return JSEvent(timestamp, value, type, number)


class Inputs():
    js_stream: str
    command_map: Dict[str, str] # command_name: #jsb2
    js_map: Dict[str, JSEvent] # jsb2: JSEvent
    last_js_timestamp: int

    def __init__(self, js_stream: str = "/dev/input/js0") -> None:
        self.js_stream = js_stream
        self.last_js_timestamp = 0
        self.command_map = {}
        self.js_map = {}
    
    def bind_command_map(self, command_map: Dict[str, str]) -> None:
        self.command_map = command_map
        self.js_map = {}
        for key in self.command_map.values():
            type = 1 if key[0:3] == 'jsb' else 2
            number = int(key[3:])
            self.js_map[key] = JSEvent(0,0,type,number)
    
    def poll(self, ) -> None:
        self._poll_js_inputs()
    
    def _poll_js_inputs(self,) -> None:
        fd = os.open(self.js_stream, os.O_RDONLY|os.O_NONBLOCK, mode=0x666)
        in_file = os.fdopen(fd, "rb", 256)
        buffer = in_file.read(256)
        in_file.close()
        if buffer and len(buffer):
            self._parse_js_buffer(buffer)
    
    def _parse_js_buffer(self, buffer: bytes) -> None:
        for i in range(0, len(buffer), JSEvent.EVENT_SIZE):
            event_buffer = buffer[i:i+JSEvent.EVENT_SIZE]
            event = JSEvent_from_bytes(event_buffer)
            if event is not None and event.map_index() in self.js_map:
                self._update_js_map_event(event)

    def _update_js_map_event(self, event: JSEvent) -> None:
        if event.timestamp < self.last_js_timestamp:
            return
        
        if self.last_js_timestamp < event.timestamp:
            self.last_js_timestamp = event.timestamp
        
        map_index = event.map_index()
        last_value = self.js_map[map_index].value
        event.normalize_value()
        event.changed = last_value != event.value
        self.js_map[map_index] = event
    
    def get_command_event(self, command: str) -> JSEvent:
        return self.js_map[self.command_map[command]]
    
    def set_synthetic(self, frame_count: int, command: str, value: int) -> JSEvent:
        last_event = self.get_command_event(command)
        event = JSEvent(frame_count, value, last_event.type, last_event.number)
        self._update_js_map_event(event)


