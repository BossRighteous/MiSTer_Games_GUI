package mgdb

type MGDBInfo struct {
	CollectionName     string
	GamesFolder        string
	SupportedSystemIds string
	BuildDate          string
	Description        string
}

type Game struct {
	GameID      int
	Name        string
	IsIndexed   int
	GenreId     int
	Genre       string
	Rating      string
	ReleaseDate string
	Developer   string
	Publisher   string
	Players     string
	Description string
}

type RDBRom struct {
	ROMName string
	Name    string
	CRC32   int
	Size    int
	Serial  string
	GameID  int
}

type GamelistRom struct {
	FileName string
	GameID   int
}

type IndexedRom struct {
	Path     string
	FileName string
	FileExt  string
	GameID   int
}

type Genre struct {
	GenreID int
	Name    string
}

type Screenshot struct {
	GameID   int
	FilePath string
	Bytes    []byte
}

type TitleScreen struct {
	GameID   int
	FilePath string
	Bytes    []byte
}
