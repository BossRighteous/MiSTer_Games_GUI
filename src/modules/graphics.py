import struct
from math import floor
from typing import Tuple
from typing import List
from typing import Optional
from typing import Union

class Color():
    staticmethod
    def fromInts(r: int, g: int, b: int) -> bytearray:
        bytes = bytearray(3)
        struct.pack_into('<B', bytes, 0, b)
        struct.pack_into('<B', bytes, 1, g)
        struct.pack_into('<B', bytes, 2, r)
        return bytes
    
    staticmethod
    def toInts(bytes: bytearray) -> List[int]:
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

    def __init__(self, width: int, height: int, initial_bytes: Optional[bytearray] = None) -> None:
        self.transparent = -1
        self.width = width
        self.height = height
        self.bytes = initial_bytes if initial_bytes is not None else bytearray(width * height * self.PX_SIZE)
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
    
    def replace_colors(self, existing_color_bytes: bytearray, new_color_bytes: bytearray):
        # assumes triplet bytes
        self.bytes = self.bytes.replace(existing_color_bytes, new_color_bytes)
    
    # copy on get to freeze reference
    def get_bytes(self) -> bytearray:
        return self.bytes[0:]


class SpriteSheet(Bitmap):
    cell_width: int
    cell_height: int
    cell_cols: int
    cell_rows: int
    cell_number: int
    max_cells: int

    def __init__(self, width: int, height: int, cell_width: int, cell_height: int, initial_bytes: Optional[bytearray] = None) -> None:
        super().__init__(width, height, initial_bytes)
        self.cell_width = cell_width
        self.cell_height = cell_height
        self.cell_cols = floor(width / cell_width)
        self.cell_rows = floor(height/cell_height)
        self.max_cells = self.cell_cols * self.cell_rows
        self.cell_number = 0
    
    def _get_position_for_cell(self, cell_number: int) -> Tuple[int, int]:
        # TODO: math wrong here
        x = (cell_number % self.cell_cols) * self.cell_width
        y = (floor(cell_number/ self.cell_cols)) * self.cell_height
        return x, y
    
    def get_current_cell_rect(self) -> bytearray:
        return self.get_cell_rect(self.cell_number)

    def get_cell_rect(self, cell_number:int):
        x, y = self._get_position_for_cell(cell_number)
        return self.get_rect(x, y, self.cell_width, self.cell_height)


class FontSheet(SpriteSheet):
    UTF_BYTE_OFFSET = 32

    def __init__(self, width: int, height: int, cell_width: int, cell_height: int, initial_bytes: Optional[bytearray] = None) -> None:
        super().__init__(width, height, cell_width, cell_height, initial_bytes)
    
    def get_cell_from_char(self, char: Union[str, bytearray]) -> int:
        return ord(char) - self.UTF_BYTE_OFFSET
    
    def get_char_at_cell(self, cell_number: Optional[int]) -> str:
        cell_number = cell_number if cell_number is not None else self.cell_number
        chr(cell_number+self.UTF_BYTE_OFFSET)

    
    def make_text_bitmap(self, chars: Union[str, bytearray], max_width: int, max_height: int, trim_size: bool) -> Bitmap:
        tmp_bmp = Bitmap(max_width, max_height)
        max_cell_count = floor((max_width/self.cell_width) * (max_height/self.cell_height))
        blit_x = 0
        blit_y = 0
        for index, char in enumerate(chars):
            if index > max_cell_count:
                break
            if blit_x + self.cell_width > max_width:
                blit_x = 0
                blit_y = blit_y + self.cell_height
            self.cell_number = self.get_cell_from_char(char)
            tmp_bmp.blit_rect(blit_x, blit_y, self.cell_width, self.cell_height, self.get_current_cell_rect())
            blit_x = blit_x + self.cell_width
        if trim_size:
            final_width = max_width if blit_y > 0 else blit_x
            final_height = blit_y + self.cell_height
            tmp_bmp.bytes = tmp_bmp.get_rect(0, 0, final_width, final_height)
            tmp_bmp.width = final_width
            tmp_bmp.height = final_height
        return tmp_bmp
