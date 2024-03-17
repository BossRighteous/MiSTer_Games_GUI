
from random import randint
from time import time
from typing import Callable, Dict, TypeVar, Deque
from modules.filesystem import FileSystem
from modules.graphics import Bitmap, Color, FontSheet
from modules.input import Inputs
from modules.loop_context import LoopContext
from modules.theme import Theme
from collections import deque
from math import floor

Arcade15Self = TypeVar("Arcade15Self", bound="Arcade15")
class Arcade15(Theme):
    ctx: LoopContext
    inputs: Inputs
    idle_queue: Deque[Callable[[Arcade15Self, LoopContext], None]]
    last_time: float
    font: FontSheet
    fps: Bitmap
    bitmaps: Dict[str, Bitmap]
    command_map: Dict[str, str]

    def __init__(self, ctx: LoopContext, inputs: Inputs) -> None:
        # Don't do heavy lifting here as it's part of the GUI constructor
        self.ctx = ctx
        self.inputs = inputs
        self.idle_queue = deque()
        self.last_time = time()
        self.bitmaps = {}
        self.command_map = {}
    
    def pre_loop(self) -> None:
        ctx = self.ctx
        self.font = FileSystem.font_sheet_from_font_dat('./themes/standard_arcade15/font.png.font.dat')
        font = self.font
        bitmaps = self.bitmaps
        cell_w = self.font.cell_width
        cell_h = self.font.cell_height

        self.bgcolor = Color.fromInts(250,128,250)
        ctx.surface.replace_colors(Color.fromInts(0,0,0),self.bgcolor)
        self.fps = font.make_text_bitmap("60", cell_w*2, cell_h, False)
        self.fps.x = ctx.width - cell_w - self.fps.width
        self.fps.y = cell_h

        self.command_map = {
            "0": "jsb0",
            "1": "jsb1",
            "2": "jsb2",
            "3": "jsb3",
            "4": "jsb4",
            "5": "jsb5",
            "6": "jsb6",
            "x": "jsa0",
            "y": "jsa1",
        }
        self.inputs.bind_command_map(self.command_map)
        bx = cell_w
        by = cell_h
        default_color = Color.fromInts(255,255,255)
        for key in self.command_map.keys():
            if key == "4" or key == "x":
                bx = cell_w
                by = by + cell_h
            if key != "x" and key != "y":
                bitmap = self.font.make_text_bitmap(key, cell_w, cell_h, False)
                bitmap.set_position(bx, by)
                bitmap.text_color = default_color
                bitmaps[key] = bitmap
            else:
                bitmap = self.font.make_text_bitmap("+", cell_w, cell_h, False)
                bitmap.set_position(bx, by)
                bitmap.text_color = default_color
                bitmaps[key+"+"] = bitmap
                bx = bx + cell_w
                bitmapm = self.font.make_text_bitmap("-", cell_w, cell_h, False)
                bitmapm.set_position(bx, by)
                bitmapm.text_color = default_color
                bitmaps[key+"-"] = bitmapm
            bx = bx + cell_w
    
    def on_tick(self,) -> None:
        curr_time = time()
        ctx = self.ctx
        is_idle = True
        surface = ctx.surface
        cell_w = self.font.cell_width
        cell_h = self.font.cell_height

        # calculate FPS int
        fps_int = floor(1/(curr_time-self.last_time))
        self.fps.bytes = self.font.make_text_bitmap(str(fps_int), cell_w*2, cell_h, False).bytes

        # Update bitmaps from input states
        for command_composite, bmp in self.bitmaps.items():
            # drop +/- for axis
            command = command_composite[0]
            event = self.inputs.get_command_event(command)

            # Can set sythetic on Windows
            #self.inputs.set_synthetic(ctx.frame_count, command, value)

            if not event.is_axis() or (event.value != 0 and ((command_composite[1] == "+" and event.value == 1) or (command_composite[1] == "-" and event.value == -1))):
                if event.is_just_pressed():
                    new_color = Color.fromInts(255,0,0)
                    bmp.replace_colors(bmp.text_color, new_color)
                    bmp.text_color = new_color
                elif event.is_just_released():
                    new_color = Color.fromInts(255,255,255)
                    bmp.replace_colors(bmp.text_color, new_color)
                    bmp.text_color = new_color
                elif event.is_pressed():
                    new_color = Color.fromInts(0,255,0)
                    bmp.replace_colors(bmp.text_color, new_color)
                    bmp.text_color = new_color
            elif event.is_just_released():
                new_color = Color.fromInts(255,255,255)
                bmp.replace_colors(bmp.text_color, new_color)
                bmp.text_color = new_color

        # composite bitmaps
        surface.blit_bitmap(self.fps)
        for bitmap in self.bitmaps.values():
            surface.blit_bitmap(bitmap)

        if is_idle and self.idle_queue:
            self.idle_queue.popleft()(self, ctx)
 
        self.last_time = curr_time
    
    def post_loop(self) -> None:
        ctx = self.ctx
        pass