class SwitchRes():
    pixel_clock: float
    hactive: int
    hbegin: int
    hend: int
    htotal: int
    vactive: int
    vbegin: int
    vend: int
    vtotal: int

    def __init__(self, modeline: str = "", refresh_rate: int = 60) -> None:
        self.pixel_clock = 4.905
        self.hactive = 256
        self.hbegin = 264
        self.hend = 287
        self.htotal = 312
        self.vactive = 240
        self.vbegin = 241
        self.vend = 244
        self.vtotal = 262
        self.refresh_rate= refresh_rate
        self._parse_modeline(modeline)
    
    def _parse_modeline(self, modeline: str) -> None:
        parts = modeline.split(" ")
        if modeline == "" or len(parts) != 9:
            return
        self.pixel_clock = float(parts[0])
        self.hactive = int(parts[1])
        self.hbegin = int(parts[2])
        self.hend = int(parts[3])
        self.htotal = int(parts[4])
        self.vactive = int(parts[5])
        self.vbegin = int(parts[6])
        self.vend = int(parts[7])
        self.vtotal = int(parts[8])