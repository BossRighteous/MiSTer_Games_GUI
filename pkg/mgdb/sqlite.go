package mgdb

import (
	"database/sql"
	"errors"
	"fmt"
	"image"
	"os"

	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/img"
	_ "github.com/mattn/go-sqlite3"
)

type DBNilError struct{}

func (err *DBNilError) Error() string {
	return "MGDB DB is nil"
}

type MGDBClient struct {
	path string
	db   *sql.DB
}

func (mgdb *MGDBClient) GetMGDBInfo() (MGDBInfo, error) {
	var info MGDBInfo
	if mgdb.db == nil {
		return info, &DBNilError{}
	}
	err := mgdb.db.QueryRow(
		"select CollectionName, GamesFolder, SupportedSystemIds, "+
			"BuildDate, Description "+
			"from MGDBInfo",
	).Scan(
		&info.CollectionName,
		&info.GamesFolder,
		&info.SupportedSystemIds,
		&info.BuildDate,
		&info.Description,
	)
	if err != nil {
		return info, err
	}
	return info, nil
}

func (mgdb *MGDBClient) GetGameList() ([]Game, error) {
	var gameList []Game
	if mgdb.db == nil {
		return gameList, &DBNilError{}
	}
	rows, err := mgdb.db.Query("select GameID, Name from Game where IsIndexed = 1 order by Name ASC")
	if err != nil {
		return gameList, err
	}
	defer rows.Close()
	for rows.Next() {
		game := Game{}
		err := rows.Scan(&game.GameID, &game.Name)
		if err != nil {
			return gameList, err
		}
		gameList = append(gameList, game)
	}
	err = rows.Err()
	if err != nil {
		return gameList, err
	}
	return gameList, nil
}

func (mgdb *MGDBClient) GetGame(gameID int) (Game, error) {
	var game Game
	if mgdb.db == nil {
		return game, &DBNilError{}
	}
	err := mgdb.db.QueryRow(
		"select Game.GameID, Game.Name, Genre.Name as GenreName, Game.GenreID, "+
			"Game.IsIndexed, Game.Rating, "+
			"Game.ReleaseDate, Game.DeveloperID, Developer.Name as DeveloperName, "+
			"Game.PublisherID, Publisher.Name as PublisherName, Game.Players, Game.Description, "+
			"Game.ScreenshotHash, Game.TitleScreenHash "+
			"from Game JOIN Genre on Genre.GenreID = Game.GenreID "+
			"JOIN Developer on Developer.DeveloperID = Game.DeveloperID "+
			"JOIN Publisher on Publisher.PublisherID = Game.PublisherID "+
			"where Game.GameID = ?", gameID,
	).Scan(
		&game.GameID,
		&game.Name,
		&game.Genre,
		&game.GenreId,
		&game.IsIndexed,
		&game.Rating,
		&game.ReleaseDate,
		&game.DeveloperID,
		&game.Developer,
		&game.PublisherID,
		&game.Publisher,
		&game.Players,
		&game.Description,
		&game.ScreenshotHash,
		&game.TitleScreenHash,
	)
	if err != nil {
		return game, err
	}
	return game, nil
}

func (mgdb *MGDBClient) GetIndexedRoms(gameID int) ([]IndexedRom, error) {
	var romList []IndexedRom
	if mgdb.db == nil {
		return romList, &DBNilError{}
	}
	rows, err := mgdb.db.Query("select FileName, FileExt, Path from IndexedRom where GameID = ? order by FileName ASC", gameID)
	if err != nil {
		return romList, err
	}
	defer rows.Close()
	for rows.Next() {
		rom := IndexedRom{GameID: gameID}
		err := rows.Scan(&rom.FileName, &rom.FileExt, &rom.Path)
		if err != nil {
			return romList, err
		}
		romList = append(romList, rom)
	}
	err = rows.Err()
	if err != nil {
		return romList, err
	}
	return romList, nil
}

func (mgdb *MGDBClient) GetGameImage(imgHash string) (image.Image, error) {
	if mgdb.db == nil {
		return nil, &DBNilError{}
	}

	var imgBlob ImageBlob
	fmt.Println("querying ImageBlob Hash", imgHash)
	err := mgdb.db.QueryRow(
		fmt.Sprintf("select Hash, Bytes from ImageBlob where Hash = '%v'", imgHash),
	).Scan(
		&imgBlob.Hash,
		&imgBlob.Bytes,
	)
	if err != nil {
		return nil, err
	}

	var screenImg *image.Image
	if imgBlob.Bytes != nil && len(imgBlob.Bytes) > 0 {
		screenImg, err = img.DecodeImageBytes(&imgBlob.Bytes)
		if err != nil {
			return nil, err
		}
	}
	return *screenImg, nil
}

func (mgdb *MGDBClient) FlushGamesIndex() error {
	if _, err := mgdb.db.Exec("UPDATE Game SET IsIndexed = 0"); err != nil {
		fmt.Println("FlushGamesIndex UPDATE Error")
		return err
	}
	if _, err := mgdb.db.Exec("DELETE FROM IndexedRom"); err != nil {
		fmt.Println("FlushGamesIndex DELETE Error")
		return err
	}
	if _, err := mgdb.db.Exec("DELETE FROM SQLITE_SEQUENCE WHERE name='IndexedRom'"); err != nil {
		fmt.Println("FlushGamesIndex DELETE SQLITE_SEQUENCE Error")
	}
	return nil
}

func (mgdb *MGDBClient) FindGameIdFromFilename(filename string) (int, error) {
	slug := SlugifyString(filename)
	var rom SlugRom
	err := mgdb.db.QueryRow(
		fmt.Sprintf("select Slug, GameID from SlugRom where Slug = '%v'", slug),
	).Scan(
		&rom.Slug,
		&rom.GameID,
	)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return 0, err
		}
		// Otherwise, no real error, just no match
		return 0, nil
	}
	return rom.GameID, nil
}

func (mgdb *MGDBClient) IndexGameRom(rom IndexedRom) (bool, error) {
	// Got a match, index it
	if _, err := mgdb.db.Exec(fmt.Sprintf("UPDATE Game SET IsIndexed = 1 WHERE GameID = %v AND IsIndexed = 0", rom.GameID)); err != nil {
		return false, err
	}

	stmt, err := mgdb.db.Prepare(
		"insert into IndexedRom(" +
			"Path, FileName, FileExt, GameID, SupportedSystemIds" +
			") values (?, ?, ?, ?, ?)",
	)
	if err != nil {
		return false, err
	}
	if _, err = stmt.Exec(
		rom.Path,
		rom.FileName,
		rom.FileExt,
		rom.GameID,
		rom.SupportedSystemIds,
	); err != nil {
		return false, err
	}
	return true, nil
}

func OpenMGDB(path string) (*MGDBClient, error) {
	mgdb := &MGDBClient{}
	mgdb.path = path

	// Check path for valid file
	_, err := os.Stat(path)
	if err != nil {
		return mgdb, err
	}

	// Open it
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return mgdb, err
	}
	mgdb.db = db

	// Check the MGDBInfo Table for validity
	_, err = mgdb.GetMGDBInfo()
	if err != nil {
		return mgdb, err
	}

	return mgdb, nil
}
