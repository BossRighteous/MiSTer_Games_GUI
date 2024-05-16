package mgdb

import (
	"database/sql"
	"os"

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
