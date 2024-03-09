from time import time, sleep
from modules.loop_context import LoopContext
from modules.theme import Theme
from modules.udp import UdpClient, SwitchRes
import argparse

from themes.arcade15 import Arcade15

class GUI():
    # Settings
    port: int
    ip: int
    modeline: str
    refresh: int
    debug: int
    run_loop: bool
    theme: Theme

    # Runtime
    client: UdpClient
    switchres: SwitchRes
    ctx: LoopContext

    def __init__(self) -> None:
        self.port = 32100
        self.ip = '192.168.0.168' #'127.0.0.1'
        self.modeline = ''
        self.refresh = 60
        self.debug = 0
        self.theme = Arcade15()

        self._parse_env()
        self._parse_args()
        self.run_loop = False

        self.client: UdpClient = UdpClient(self.port, self.ip)
        self.switchres: SwitchRes = SwitchRes(self.modeline, self.refresh)
        self.ctx = LoopContext(self.switchres.hactive, self.switchres.vactive, self.switchres.refresh_rate)

    def _parse_env(self) -> None:
        # Stub
        pass

    def _parse_args(self,) -> None:
        parser = argparse.ArgumentParser(description='Low-resolution analog friendly MiSTer Python script GUI for your game library')
        parser.add_argument('--port', '-p', type=int, help="Optional GroovyMister UDP bind port", required=False)
        parser.add_argument('--ip', '-i', type=str, help="Optional GroovyMister Host IP", required=False)
        parser.add_argument('--modeline', '-m', type=str, help="Optional Modeline for analog sync. Default is 256x240p@60hz", required=False)
        parser.add_argument('--refresh', '-r', type=int, help="Optional Refresh rate in hz or fps", required=False)
        parser.add_argument('--debug', '-d', type=int, help="Optional Debug level", required=False)
        args = parser.parse_args()
        if args.port is not None:
            self.port = args.port
        if args.ip is not None:
            self.ip = args.ip
        if args.modeline is not None:
            self.modeline = args.modeline
        if args.refresh is not None:
            self.refresh = args.refresh
        if args.debug is not None:
            self.port = args.debug

    
    def run(self) -> None:
        self.run_loop = True
        self._start_client_session()
        self._pre_loop()
        self._loop()
        self._post_loop()
        self._end_client_session()

    def _start_client_session(self) -> None:
        self.client.cmd_init()
        sleep(0.25)
        self.client.cmd_switchres(self.switchres)
        sleep(0.25)
    
    def _end_client_session(self) -> None:
        self.client.cmd_close()
        print('cmd_close')
    
    def _pre_loop(self) -> None:
        self.theme.pre_loop(self.ctx)
    
    def _loop(self) -> None:
        ctx = self.ctx
        frame_limit = self.switchres.refresh_rate*10
        try:
            while self.run_loop and self.ctx.frame_count < frame_limit:
                frame_start = time()
                self.client.cmd_blit(ctx.surface.get_bytes())

                ctx.frame_count = ctx.frame_count + 1
                ctx.frame_time_remaining = ctx.frame_tick - (time() - frame_start)
                self._on_tick()

                frame_time = time()-frame_start
                sleep_time = ctx.frame_tick - frame_time
                sleep_time = sleep_time if sleep_time > 0 else 0
                print(f"blit done "+str(sleep_time)+" "+str(frame_time))
                sleep(sleep_time)
        except KeyboardInterrupt:
            self.run_loop = False

    def _on_tick(self):
            self.theme.on_tick(self.ctx)

    def _post_loop(self) -> None:
        # clean up anything in case we go sleep and re-run or something
        self.theme.post_loop(self.ctx)
        self.ctx.reset()