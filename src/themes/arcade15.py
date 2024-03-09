
from random import randint
from typing import Callable, TypeVar, Deque
from modules.filesystem import FileSystem
from modules.graphics import Bitmap, Color, FontSheet
from modules.loop_context import LoopContext
from modules.theme import Theme
from collections import deque


Arcade15Self = TypeVar("Arcade15Self", bound="Arcade15")
class Arcade15(Theme):
    idle_queue: Deque[Callable[[Arcade15Self, LoopContext], None]]
    
    def pre_loop(self, ctx: LoopContext) -> None:
        self.bgcolor = Color.fromInts(250,128,250)
        ctx.surface.replace_colors(Color.fromInts(0,0,0),self.bgcolor)
        #self.font: FontSheet = FileSystem.font_sheet_from_font_dat("./themes/standard_arcade15/font.png.font.dat")
        self.jeff: Bitmap = FileSystem.bitmap_from_image_dat("./themes/standard_arcade15/jeff.png.dat")
        self.jeff_mask = Bitmap(self.jeff.width, self.jeff.height)
        self.jeff_mask.replace_colors(Color.fromInts(0,0,0),self.bgcolor)
        #self.survivor = self.font.make_text_bitmap("Previously on SURVIVOR!", ctx.width, ctx.height, True)
        #self.file_box = Bitmap(1,1)
        #self.dir_iterator = FileSystem.get_direntry_iterator("./")
        self.idle_queue = deque()
    
    def on_tick(self, ctx: LoopContext) -> None:
        is_idle = True

        if ctx.frame_count % 60 == 0:
            is_idle = False

            frame_at_queue = ctx.frame_count

            def print_queue(self: Arcade15Self, ctx: LoopContext) -> None:
                print(len(self.idle_queue))
                print(ctx.frame_count)
                print(frame_at_queue)

            def jeff_queue(self: Arcade15Self, ctx: LoopContext) -> None:
                jeff = self.jeff
                jeff_mask = self.jeff_mask
                surface = ctx.surface
                # clear with mask
                surface.blit_rect(jeff.x, jeff.y, jeff_mask.width, jeff_mask.height, jeff_mask.bytes)
                self.jeff.x = self.jeff.x+1
                self.jeff.y = self.jeff.y+1
                surface.blit_rect(jeff.x, jeff.y, jeff.width, jeff.height, jeff.bytes)

            for i in range(29):
                self.idle_queue.append(print_queue)
                if self.jeff.x + self.jeff.width < ctx.width and self.jeff.y + self.jeff.height < ctx.height:
                    self.idle_queue.append(jeff_queue)

        if is_idle and self.idle_queue:
            self.idle_queue.popleft()(self, ctx)

        # font = self.font
        # jeff = self.jeff
        # text_box = self.survivor
        
        # file_name = ""

        # x_pos = randint(0, width)
        # y_pos = randint(0, height)
        # rect = Bitmap(randint(0, width - x_pos), randint(0, height-y_pos))
        # fill_color = Color.fromInts(randint(0,255), randint(0,255), randint(0,255))
        # rect.fill_rect(0, 0, rect.width, rect.height, fill_color)
        # background.blit_rect(x_pos, y_pos, rect.width, rect.height, rect.get_bytes())

        # #composite text
        # text_box.replace_colors(Color.fromInts(0,0,0), fill_color)
        # background.blit_rect(5,10, text_box.width, text_box.height, text_box.get_bytes())
        
        # background.blit_rect(20+randint(-2,2), 30+randint(-2,2), jeff.width, jeff.height, jeff.get_bytes())

        # if ctx.frame_count % 60 == 0:
        #     file = next(self.dir_iterator, None)
        #     if file is None:
        #         self.dir_iterator = FileSystem.get_direntry_iterator("./")
        #         file = next(self.dir_iterator, None)
        #     file_name = file.name
        #     self.file_box = font.make_text_bitmap(file_name, width-10, 9, False)
        #     self.file_box.replace_colors(Color.fromInts(0,0,0), Color.fromInts(135,232,237))
        #     self.file_box.replace_colors(Color.fromInts(255,255,255), Color.fromInts(100,100,100))

        #     ctx.idle_queue
        # background.blit_rect(5,210, self.file_box.width, self.file_box.height, self.file_box.get_bytes())
    
    def post_loop(self, ctx: LoopContext) -> None:
        pass