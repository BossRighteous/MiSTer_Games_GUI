package mgdb

import (
	"database/sql"
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

func (mgdb MGDBClient) GetMGDBInfo() (MGDBInfo, error) {
	var info MGDBInfo
	if mgdb.db == nil {
		return info, &DBNilError{}
	}
	err := mgdb.db.QueryRow(
		"select CoreDir, CoreName, CoreSlug, "+
			"BuildDate, Description, IndexDir "+
			"from MGDBInfo",
	).Scan(
		&info.CoreDir,
		&info.CoreName,
		&info.CoreSlug,
		&info.BuildDate,
		&info.Description,
		&info.IndexDir,
	)
	if err != nil {
		return info, err
	}
	return info, nil
}

func (mgdb MGDBClient) GetGameList() ([]GameListItem, error) {
	var gameList []GameListItem
	if mgdb.db == nil {
		return gameList, &DBNilError{}
	}
	rows, err := mgdb.db.Query("select GameID, Name from Game order by Name ASC")
	if err != nil {
		return gameList, err
	}
	defer rows.Close()
	for rows.Next() {
		game := GameListItem{}
		err := rows.Scan(&game.GameID, &game.Name)
		if err != nil {
			return gameList, err
		}
		if game.GameID == 0 {
			continue
		}
		gameList = append(gameList, game)
	}
	err = rows.Err()
	if err != nil {
		return gameList, err
	}
	return gameList, nil
}

func (mgdb MGDBClient) GetGame(gameID int) (Game, error) {
	var game Game
	if mgdb.db == nil {
		return game, &DBNilError{}
	}
	err := mgdb.db.QueryRow(
		"select Game.GameID, Game.Name, Genre.Name as GenreName, Game.Rating, "+
			"Game.ReleaseDate, Game.Developer, Game.Publisher, Game.Players, Game.Description "+
			"from Game JOIN Genre on Genre.GenreID = Game.GenreID "+
			"where Game.GameID = ?", gameID,
	).Scan(
		&game.GameID,
		&game.Name,
		&game.Genre,
		&game.Rating,
		&game.ReleaseDate,
		&game.Developer,
		&game.Publisher,
		&game.Players,
		&game.Description,
	)
	if err != nil {
		return game, err
	}
	return game, nil
}

func (mgdb MGDBClient) getGameImage(table string, gameID int) (*image.Image, error) {
	if mgdb.db == nil {
		return nil, &DBNilError{}
	}
	var screen Screenshot
	fmt.Println("querying gameID", gameID)
	err := mgdb.db.QueryRow(
		fmt.Sprintf("select GameID, Bytes from "+table+" where GameID = %v", gameID),
	).Scan(
		&screen.GameID,
		&screen.Bytes,
	)
	if err != nil {
		return nil, err
	}

	var screenImg *image.Image
	if screen.Bytes != nil && len(screen.Bytes) > 0 {
		screenImg, err = img.DecodeImageBytes(&screen.Bytes)
		if err != nil {
			return nil, err
		}
	}
	return screenImg, nil
}

func (mgdb MGDBClient) GetGameScreenshot(gameID int) (*image.Image, error) {
	return mgdb.getGameImage("Screenshot", gameID)
}

func (mgdb MGDBClient) GetGameTitleScreen(gameID int) (*image.Image, error) {
	return mgdb.getGameImage("TitleScreen", gameID)
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
