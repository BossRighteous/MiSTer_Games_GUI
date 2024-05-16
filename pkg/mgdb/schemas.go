package mgdb

type MGDBInfo struct {
	CoreDir     string
	CoreName    string
	CoreSlug    string
	BuildDate   string
	Description string
	IndexDir    string
}

type GameListItem struct {
	GameID int
	Name   string
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
	FileName string
	GameID   int
}

type IndexedRom struct {
	Path     string
	FileName string
	GameID   int
}

type Genre struct {
	GenreID int
	Name    string
}

type Screenshot struct {
	GameID   int
	FileName string
	Bytes    []byte
}

type TitleScreen struct {
	GameID   int
	FileName string
	Bytes    []byte
}
