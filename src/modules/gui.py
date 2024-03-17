from time import time, sleep
from modules.input import Inputs
from modules.loop_context import LoopContext
from modules.settings import Settings
from modules.theme import Theme
from modules.udp import UdpClient
from modules.switchres import SwitchRes
from themes.arcade15 import Arcade15

class GUI():
    # Settings
    settings: Settings
    theme: Theme

    # Runtime
    client: UdpClient
    switchres: SwitchRes
    ctx: LoopContext
    inputs: Inputs

    def __init__(self, settings: Settings) -> None:
        self.settings = settings
        self.run_loop = False

        # ignore theme setting for now

        self.client: UdpClient = UdpClient(settings.mister_port, settings.mister_ip)
        self.switchres: SwitchRes = SwitchRes(settings.modeline, settings.refresh_rate)
        self.ctx = LoopContext(self.switchres.hactive, self.switchres.vactive, self.switchres.refresh_rate)
        self.inputs = Inputs()
        self.theme = Arcade15(self.ctx, self.inputs)
    
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
        self.theme.pre_loop()
    
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
                #print(f"blit done "+str(sleep_time)+" "+str(frame_time))
                sleep(sleep_time)
        except KeyboardInterrupt:
            self.run_loop = False

    def _on_tick(self):
            self.inputs.poll()
            self.theme.on_tick()

    def _post_loop(self) -> None:
        # clean up anything in case we go sleep and re-run or something
        self.theme.post_loop()
        self.ctx.reset()