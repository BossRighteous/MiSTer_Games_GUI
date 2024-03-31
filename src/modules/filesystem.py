import struct
from os import scandir, DirEntry
from typing import Union
from .graphics import FontSheet
from .graphics import Bitmap

class FileSystem:
    staticmethod
    def read_dat_bytes(path: str) -> bytes:
        with open(path, mode='rb') as file: # b is important -> binary
            return file.read()
    
    staticmethod
    def font_sheet_from_font_dat(path: str) -> FontSheet:
        buffer = FileSystem.read_dat_bytes(path)
        width, = struct.unpack_from("<H", buffer, 0)
        height, = struct.unpack_from("<H", buffer, 2)
        cell_width, = struct.unpack_from("<B", buffer, 4)
        cell_height, = struct.unpack_from("<B", buffer, 5)
        img_bytes = buffer[6:]
        sprite = FontSheet(width, height, cell_width, cell_height, img_bytes)
        # sprite.blit_rect(0, 0, width, height, img_bytes)
        print(sprite.width, sprite.height, sprite.cell_width, sprite.cell_height, len(img_bytes))
        return sprite
    
    staticmethod
    def bitmap_from_image_dat(path: str) -> Bitmap:        
        buffer = FileSystem.read_dat_bytes(path)
        width, = struct.unpack_from("<H", buffer, 0)
        height, = struct.unpack_from("<H", buffer, 2)
        img_bytes = buffer[4:]
        bitmap = Bitmap(width, height, img_bytes)
        #sprite.blit_rect(0, 0, width, height, img_bytes)
        print(bitmap.width, bitmap.height, len(img_bytes))
        return bitmap
    
    staticmethod
    def get_direntry_iterator(path: str):
        return scandir(path)
    
    def next_direntry(iterator) -> Union[DirEntry, None]:
        return next(iterator, None)

