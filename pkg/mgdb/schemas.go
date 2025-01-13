package mgdb

type MGDBInfo struct {
	CollectionName     string
	GamesFolder        string
	SupportedSystemIds string
	BuildDate          string
	MGDBVersion        string
	Description        string
}

type Game struct {
	GameID          int
	Name            string
	IsIndexed       int
	GenreId         int
	Genre           string
	Rating          string
	ReleaseDate     string
	DeveloperID     int
	Developer       string
	PublisherID     int
	Publisher       string
	Players         string
	Description     string
	ExternalID      string
	ScreenshotHash  string
	TitleScreenHash string
}

type SlugRom struct {
	Slug               string
	GameID             int
	SupportedSystemIds string
}

type RomCrc struct {
	CRC32 string
	Slug  string
}

type IndexedRom struct {
	Path               string
	FileName           string
	FileExt            string
	GameID             int
	SupportedSystemIds string
}

type Genre struct {
	GenreID int
	Name    string
}

type Developer struct {
	DeveloperID int
	Name        string
}

type Publisher struct {
	PublisherID int
	Name        string
}

type ImageBlob struct {
	Hash  string
	Bytes []byte
}
