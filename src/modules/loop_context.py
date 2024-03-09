from typing import Callable, Dict, List, TypeVar
from modules.graphics import Bitmap


CtxSelf = TypeVar("CtxSelf", bound="LoopContext")
class LoopContext():
    frame_count: int
    frame_tick: float
    frame_time_remaining: float
    width: int
    height: int
    surface: Bitmap

    def __init__(self, width: int, height: int, refresh_rate: int) -> None:
        self.frame_tick = 1 / refresh_rate
        self.width = width
        self.height = height
        self.reset()
    
    def reset(self) -> None:
        self.frame_count = 0
        self.frame_time_remaining = 0
        self.surface = Bitmap(self.width,self.height)