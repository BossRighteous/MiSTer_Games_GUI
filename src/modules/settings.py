from configparser import ConfigParser

class Settings():
    mister_port: int
    mister_ip: str
    modeline: str
    refresh_rate: int
    debug_level: int
    theme: str
    games_path: str

    def __init__(self) -> None:
        self.mister_port = 32100
        self.mister_ip = '127.0.0.1'
        self.modeline = '4.905 256 264 287 312 240 241 244 262'
        self.refresh = 60
        self.debug = 0
        self.theme = 'arcade15'

def get_ini_settings() -> Settings:
    settings = Settings()
    config = ConfigParser()
    config.read('settings.ini')
    if not config.has_section('settings'):
        return settings

    mister_port = config['settings'].get('mister_port', None)
    if mister_port is not None:
        settings.mister_port = config['settings'].getint('mister_port')
    
    mister_ip = config['settings'].get('mister_ip', None)
    if mister_ip is not None:
        settings.mister_ip = config['settings'].get('mister_ip')
    
    modeline = config['settings'].get('modeline', None)
    if modeline is not None:
        settings.modeline = config['settings'].get('modeline')
    
    refresh_rate = config['settings'].get('refresh_rate', None)
    if refresh_rate is not None:
        settings.refresh_rate = config['settings'].getint('refresh_rate')
    
    debug_level = config['settings'].get('refresh_debug_levelrate', None)
    if debug_level is not None:
        settings.debug_level = config['settings'].getint('debug_level')
    
    theme = config['settings'].get('theme', None)
    if theme is not None:
        settings.theme = config['settings'].get('theme')
    
    games_path = config['settings'].get('games_path', None)
    if games_path is not None:
        settings.games_path = config['settings'].get('games_path')
    
    return settings