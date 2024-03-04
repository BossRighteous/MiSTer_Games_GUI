import struct

class Color():
    staticmethod
    def fromInts(r: int, g: int, b: int) -> bytearray:
        bytes = bytearray(3)
        struct.pack_into('<B', bytes, 0, b)
        struct.pack_into('<B', bytes, 1, g)
        struct.pack_into('<B', bytes, 2, r)
        return bytes
    
    staticmethod
    def toInts(bytes: bytearray) -> list[int]:
        return [
            struct.unpack("<B",bytes[0]),
            struct.unpack("<B",bytes[1]),
            struct.unpack("<B",bytes[2]),
        ]


class Bitmap():
    PX_SIZE = 3
    # Byte Order: BGR
    width: int
    height: int
    bytes: bytearray
    transparent: int

    # Allow position for value storage, don't pre-calculate
    x: int
    y: int

    def __init__(self, width: int, height: int) -> None:
        self.transparent = -1
        self.width = width
        self.height = height
        self.bytes = bytearray(width * height * self.PX_SIZE)
        self.x = 0
        self.y = 0
    
    def set_pixel(self, x: int, y: int, color_bytes:bytearray) -> None:
        offset: int = self._get_pixel_byte_offset(x, y)
        self.bytes[offset] = color_bytes[0]
        self.bytes[offset+1] = color_bytes[1]
        self.bytes[offset+2] = color_bytes[2]
    
    def get_pixel(self, x: int, y: int):
        offset: int = self._get_pixel_byte_offset(x, y)
        return self.bytes[offset:offset+3]

    def _get_pixel_byte_offset(self, x: int, y: int) -> int:
        return (y * self.width * self.PX_SIZE) + (x * self.PX_SIZE)

    def fill_rect(self, x: int, y:int, width: int, height: int, color_bytes:bytearray) -> None:
        row_byte_width = self.width * self.PX_SIZE
        fill_bytes_row: bytearray = color_bytes * width
        fill_bytes_width = width * self.PX_SIZE
        offset: int = self._get_pixel_byte_offset(x, y)
        new_bytes = self.bytes[0:offset]
        for _row in range(height):
            new_bytes.extend(fill_bytes_row)
            new_bytes.extend(self.bytes[offset + fill_bytes_width : offset + row_byte_width])
            offset = offset + row_byte_width
        new_bytes.extend(self.bytes[offset + fill_bytes_width:])
        self.bytes = new_bytes
    
    def get_rect(self, x: int, y:int, width: int, height: int) -> bytearray:
        row_byte_width = self.width * self.PX_SIZE
        fill_bytes_width = width * self.PX_SIZE
        bytes = bytearray()
        offset: int = self._get_pixel_byte_offset(x, y)
        for _row in range(height):
            bytes.extend(self.bytes[offset : offset + fill_bytes_width])
            offset = offset + row_byte_width
        return bytes

    def blit_rect(self, x: int, y:int, width: int, height: int, source_rect_bytes: bytearray) -> None:
        row_byte_width = self.width * self.PX_SIZE
        fill_byte_width = width * self.PX_SIZE
        offset: int = self._get_pixel_byte_offset(x, y)
        source_offset: int = 0
        new_bytes = self.bytes[0:offset]
        for _row in range(height):
            new_bytes.extend(source_rect_bytes[source_offset : source_offset + fill_byte_width])
            new_bytes.extend(self.bytes[offset + fill_byte_width : offset + row_byte_width])
            offset = offset + row_byte_width
            source_offset = source_offset + fill_byte_width
        new_bytes.extend(self.bytes[offset:])
        self.bytes = new_bytes
    
    # copy on get to freeze reference
    def get_bytes(self) -> bytearray:
        return self.bytes[0:]