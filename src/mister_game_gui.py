import struct
from time import time
from time import sleep
from modules.udp import UdpClient
from modules.udp import SwitchRes
from modules.graphics import Bitmap
from modules.graphics import Color
from random import seed
from random import randint

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

frame_count = 0
frame_limit = switchres.refresh_rate*10
frame_tick = 1000000/switchres.refresh_rate
last_frame = time()
px_offset = 0
while frame_count < frame_limit:
    frame_count = frame_count + 1
    frame_start = time()

    x_pos = randint(0, width)
    y_pos = randint(0, height)
    rect = Bitmap(randint(0, width - x_pos), randint(0, height-y_pos))
    fill_color = Color.fromInts(randint(0,255), randint(0,255), randint(0,255))
    rect.fill_rect(0, 0, rect.width, rect.height, fill_color)
    background.blit_rect(x_pos, y_pos, rect.width, rect.height, rect.get_bytes())
    client.cmd_blit(background.get_bytes())

    sleep_time = (frame_tick - (time()-frame_start))/1000000
    print(f"blit done "+str(sleep_time)+" "+str(time()-frame_start))
    sleep(sleep_time)

client.cmd_close()
print("cmd_close")
print("exiting")
