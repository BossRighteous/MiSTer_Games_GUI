from abc import ABC, abstractmethod
from typing import Callable
from modules.loop_context import LoopContext


class Theme(ABC):
    @abstractmethod
    def pre_loop(self, ctx: LoopContext) -> None:
        pass

    @abstractmethod
    def on_tick(self, ctx: LoopContext) -> None:
        pass

    @abstractmethod
    def post_loop(self, ctx: LoopContext) -> None:
        pass