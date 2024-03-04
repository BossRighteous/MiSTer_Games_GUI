import struct
from socket import socket
from socket import AF_INET
from socket import SOCK_DGRAM

class SwitchRes():
    pixel_clock: float = 4.905
    hactive: int = 256
    hbegin: int = 264
    hend: int = 287
    htotal: int = 312
    vactive: int = 240
    vbegin: int = 241
    vend: int = 244
    vtotal: int = 262
    refresh_rate: int = 60

class UdpClient():
    UDP_IP = "127.0.0.1"
    UDP_PORT = 32100
    MTU_BLOCK_SIZE = 1470

    sock: socket
    width: int
    height: int
    frame: int

    def __init__(self):
        self.sock = socket(AF_INET, # Internet
            SOCK_DGRAM) # UDP        
        #self.sock.bind((UDP_IP, UDP_PORT))
        self.width = 256
        self.height = 240
        self.frame = 0

    def send_packet(self, packet: bytearray):
        self.sock.sendto(packet, (self.UDP_IP, self.UDP_PORT))

    def cmd_init(self):
        buffer: bytearray = bytearray(4)
        struct.pack_into('<B', buffer, 0, 2) # CMD
        struct.pack_into('<B', buffer, 1, 0) # lz4 compression flag
        struct.pack_into('<B', buffer, 2, 0) # sound rate flag
        struct.pack_into('<B', buffer, 3, 0) # sound channel
        self.send_packet(buffer)
        #print([i for i in buffer])

    def cmd_close(self):
        buffer: bytearray = bytearray(1)
        struct.pack_into('>B', buffer, 0, 1) # CMD
        self.send_packet(buffer)
        #print([i for i in buffer])

    def cmd_switchres(self, switchres: SwitchRes):
        self.width = switchres.hactive
        self.height = switchres.vactive
        buffer: bytearray = bytearray(26)
        struct.pack_into('<B', buffer, 0, 3) # CMD
        struct.pack_into('<d', buffer, 1, switchres.pixel_clock)
        struct.pack_into('<H', buffer, 9, switchres.hactive)
        struct.pack_into('<H', buffer, 11, switchres.hbegin)
        struct.pack_into('<H', buffer, 13, switchres.hend)
        struct.pack_into('<H', buffer, 15, switchres.htotal)
        struct.pack_into('<H', buffer, 17, switchres.vactive)
        struct.pack_into('<H', buffer, 19, switchres.vbegin)
        struct.pack_into('<H', buffer, 21, switchres.vend)
        struct.pack_into('<H', buffer, 23, switchres.vtotal)
        struct.pack_into('<B', buffer, 25, 0) # interlace
        self.send_packet(buffer)
        #print([i for i in buffer])

    def cmd_blit(self, frame_buffer: bytearray):
        self.frame = self.frame + 1
        buffer: bytearray = bytearray(9)
        struct.pack_into('<B', buffer, 0, 6) # CMD
        struct.pack_into('<I', buffer, 1, self.frame)
        struct.pack_into('<H', buffer, 5, 0) # vsyncAuto
        struct.pack_into('<B', buffer, 7, 0) # lz4 blockSize & 0xff
        struct.pack_into('<B', buffer, 8, 0) # lz4 blockSize >> 8
        self.send_packet(buffer)
        #print([i for i in buffer])
        self.send_mtu(frame_buffer)
    
    def send_mtu(self, buffer: bytearray):
        bytes_to_send = len(buffer)
        chunk_max_size = self.MTU_BLOCK_SIZE
        chunk_size: int = 0
        offset: int = 0
        while bytes_to_send > 0:
            chunk_size = chunk_max_size if bytes_to_send > chunk_max_size else bytes_to_send
            bytes_to_send = bytes_to_send - chunk_size
            self.send_packet(buffer[offset : offset+chunk_size])
            #print([i for i in buffer[offset : offset+chunk_size]])
            offset += chunk_size

    def wait_for_ack(self):
        # Untested Placeholder
        while True:
            data, addr = self.sock.recvfrom(self.MTU_BLOCK_SIZE)
            print("received ack")
