import struct
from time import time
from time import sleep
from modules.udp import UdpClient
from modules.udp import SwitchRes
from modules.graphics import Bitmap
from modules.graphics import Color
from modules.graphics import FontSheet
from random import seed
from random import randint
from modules.filesystem import FileSystem

seed()

client: UdpClient = UdpClient()
switchres: SwitchRes = SwitchRes()
width = switchres.hactive
height = switchres.vactive

client.cmd_init()
sleep(0.25)
client.cmd_switchres(switchres)
sleep(0.25)

background = Bitmap(width, height)
background.fill_rect(0, 0, width, height, Color.fromInts(128,128,128))

font: FontSheet = FileSystem.font_sheet_from_font_dat("./themes/standard_arcade15/font.png.font.dat")
jeff: Bitmap = FileSystem.bitmap_from_image_dat("./themes/standard_arcade15/jeff.png.dat")

dir_iterator = FileSystem.get_direntry_iterator("./")
file_name = ""


frame_count = 0
frame_limit = switchres.refresh_rate*10
frame_tick = 1/switchres.refresh_rate
last_frame = time()
px_offset = 0
while frame_count < frame_limit:
    frame_count = frame_count + 1
    frame_start = time()
    client.cmd_blit(background.get_bytes())

    x_pos = randint(0, width)
    y_pos = randint(0, height)
    rect = Bitmap(randint(0, width - x_pos), randint(0, height-y_pos))
    fill_color = Color.fromInts(randint(0,255), randint(0,255), randint(0,255))
    rect.fill_rect(0, 0, rect.width, rect.height, fill_color)
    background.blit_rect(x_pos, y_pos, rect.width, rect.height, rect.get_bytes())

    #composite text
    text_box = font.make_text_bitmap("Previously on SURVIVOR!", width, height, True)
    text_box.replace_colors(Color.fromInts(0,0,0), fill_color)
    background.blit_rect(5,10, text_box.width, text_box.height, text_box.get_bytes())
    
    background.blit_rect(20+randint(-2,2), 30+randint(-2,2), jeff.width, jeff.height, jeff.get_bytes())
    #break

    if frame_count % 60 == 0:
        file = next(dir_iterator, None)
        if file is None:
            dir_iterator = FileSystem.get_direntry_iterator("./")
            file = next(dir_iterator, None)
        file_name = file.name
    file_box = font.make_text_bitmap(file_name, width-10, 9, False)
    file_box.replace_colors(Color.fromInts(0,0,0), Color.fromInts(135,232,237))
    file_box.replace_colors(Color.fromInts(255,255,255), Color.fromInts(100,100,100))
    background.blit_rect(5,210, file_box.width, file_box.height, file_box.get_bytes())
    frame_time = time()-frame_start
    sleep_time = frame_tick - frame_time
    sleep_time = sleep_time if sleep_time > 0 else 0
    print(f"blit done "+str(sleep_time)+" "+str(frame_time))
    sleep(sleep_time)

client.cmd_close()
print("cmd_close")
print("exiting")
