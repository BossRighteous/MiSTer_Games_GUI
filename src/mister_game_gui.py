from modules.gui import GUI
from random import seed
from modules.settings import get_ini_settings


seed()
gui = GUI(get_ini_settings())
gui.run()