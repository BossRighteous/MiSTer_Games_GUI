
from os import DirEntry
from random import randint
from time import time
from typing import Callable, Dict, List, TypeVar, Deque
from modules.filesystem import FileSystem
from modules.graphics import Bitmap, Color, FontSheet
from modules.input import Inputs
from modules.loop_context import LoopContext
from modules.theme import Theme
from collections import deque
from math import floor
from enum import IntEnum

class Screen(IntEnum):
    SCANNING = 0
    LISTING = 1
    IMAGE_FIT = 2
    DESCRIPTION = 3
    META = 4
    # LOADING_GAME = 5
    # IMAGE_FULL = 6


class Arcade15(Theme):
    ctx: LoopContext
    inputs: Inputs
    idle_queue: Deque[Callable[[], None]]
    screen: Screen
    bitmaps: Dict[str, Bitmap]
    command_map: Dict[str, str]
    is_screen_changing: bool
    is_loading: bool
    is_input_polling: bool

    font: FontSheet
    fps: Bitmap
    last_time: float

    dir_path: str = './.rsync'
    dir_list: List[DirEntry, str]

    show_fps = True


    def __init__(self, ctx: LoopContext, inputs: Inputs) -> None:
        # Don't do heavy lifting here as it's part of the GUI constructor
        self.ctx = ctx
        self.inputs = inputs
        self.reset()
    
    def reset(self):
        self.screen = Screen.SCANNING
        self.is_screen_changing = False
        self.is_loading = False
        self.is_input_polling = False
        self.idle_queue = deque()
        self.last_time = time()
        self.bitmaps = {}
        self.command_map = {}
    
    def pre_loop(self) -> None:
        ctx = self.ctx
        self.font: FontSheet = FileSystem.font_sheet_from_font_dat('./themes/standard_arcade15/font.png.font.dat')
        cell_w = self.font.cell_width
        cell_h = self.font.cell_height

        if self.show_fps:
            self.fps = self.font.make_text_bitmap("60", cell_w*2, cell_h, False)
            self.fps.x = ctx.width - cell_w - self.fps.width
            self.fps.y = cell_h

        self.bgcolor = Color.fromInts(250,128,250)
        self.is_screen_changing = True
        self.clear_surface()

        self.command_map = {
            "0": "jsb0",
            "1": "jsb1",
            "2": "jsb2",
            "3": "jsb3",
            "x": "jsa0",
            "y": "jsa1",
        }
        self.inputs.bind_command_map(self.command_map)
    
    def on_tick(self,) -> None:
        curr_time = time()
        self.inputs.poll()
        ctx = self.ctx
        surface = ctx.surface

        if self.is_input_polling:
            self.inputs.poll()
        
        if self.screen == Screen.SCANNING:
            if self.is_screen_changing:
                self.is_screen_changing = False
                self.change_screen_to_scanning()
            else:
                self.on_tick_scanning()
        elif self.screen == Screen.LISTING:
            if self.is_screen_changing:
                self.is_screen_changing = False
                self.change_screen_to_listing()
            else:
                self.on_tick_listing()
        elif self.screen == Screen.IMAGE_FIT:
            if self.is_screen_changing:
                self.is_screen_changing = False
                self.change_screen_to_image_fit()
            else:
                self.on_tick_image_fit()
        elif self.screen == Screen.DESCRIPTION:
            if self.is_screen_changing:
                self.is_screen_changing = False
                self.change_screen_to_description()
            else:
                self.on_tick_description()
        elif self.screen == Screen.META:
            if self.is_screen_changing:
                self.is_screen_changing = False
                self.change_screen_to_meta()
            else:
                self.on_tick_meta()  

        # composite bitmaps
        for bitmap in self.bitmaps.values():
            surface.blit_bitmap(bitmap)

        if self.is_loading:
            self.render_loading()

        if self.show_fps:
            self.render_fps(curr_time)

        #if is_idle and self.idle_queue:
        #    self.idle_queue.popleft()()
    
    def post_loop(self) -> None:
        pass

    def clear_surface(self) -> None:        
        self.ctx.surface.fill_rect(0, 0, self.ctx.surface.width, self.ctx.surface.height, self.bgcolor)

    def render_fps(self, curr_time: float) -> None:
        fps_int = floor(1/(curr_time-self.last_time))
        self.fps.bytes = self.font.make_text_bitmap(str(fps_int), self.font.cell_width*2, self.font.cell_height, False).bytes
        self.ctx.surface.blit_bitmap(self.fps)
        self.last_time = curr_time

    def render_loading(self) -> None:
        pass

    def change_screen_to_scanning(self):
        self.is_loading = True
        self.is_input_polling = False
        self.clear_surface()
        # Queue directory iteration
        # closure scope for self
        self.dir_list = ['../']
        scandir_iter = FileSystem.get_direntry_iterator(self.dir_path)
        def iterate_dir():
            requeue = True
            for i in range(10):
                dir_entry = FileSystem.next_direntry(scandir_iter)
                if dir_entry is None:
                    requeue = False
                    break
                self.dir_list.append(dir_entry)
            if requeue:
                self.idle_queue.append(iterate_dir)
        self.idle_queue.append(iterate_dir)

    def on_tick_scanning(self):
        # Burn queue and transition
        if self.idle_queue:
            self.idle_queue.popleft()()
        else:
            self.is_screen_changing = True
            self.screen = Screen.LISTING

    def change_screen_to_listing(self):
        # prep blit for listing
        pass

    def on_tick_listing(self):
        pass

    def change_screen_to_image_fit(self):
        pass

    def on_tick_image_fit(self):
        pass

    def change_screen_to_description(self):
        pass

    def on_tick_description(self):
        pass

    def change_screen_to_meta(self):
        pass

    def on_tick_meta(self):
        pass
