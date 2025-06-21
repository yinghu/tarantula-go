package main

import (
	"errors"
	"github.com/jackc/pgx/v5"
)

const (
	ASSET_INDEX_SQL_SCHEMA    string = "CREATE TABLE IF NOT EXISTS asset_index (system_id BIGINT NOT NULL,name VARCHAR(100) NOT NULL,file_index VARCHAR(100) NOT NULL,upload_time TIMESTAMP DEFAULT NOW(),PRIMARY KEY(system_id,name))"
	ASSET_INDEX_INSERT_UPDATE string = "INSERT INTO asset_index AS a (system_id,name,file_index,upload_time) VALUES($1,$2,$3,$4) ON CONFLICT (system_id,name) DO UPDATE SET file_index = $5, upload_time = $6  WHERE a.system_id = $7 AND a.name = $8"
	ASSET_INDEX_SELECT        string = "SELECT file_index FROM asset_index AS a WHERE a.system_id = $1 AND a.name = $2"
)

func (s *AssetService) createSchema() error {
	
	_, err := s.Sql.Exec(ASSET_INDEX_SQL_SCHEMA)
	return err
}

func (s *AssetService) loadAssetIndex(aindex *AssetIndex) error {
	err := s.Sql.Query(func(row pgx.Rows) error {
		var fileIndex string
		err := row.Scan(&fileIndex)
		if err != nil {
			return err
		}
		if fileIndex == "" {
			return errors.New("index not found")
		}
		aindex.fileIndex = fileIndex
		return nil
	}, ASSET_INDEX_SELECT, aindex.systemId, aindex.name)
	if err != nil {
		return err
	}
	return nil
}

func (s *AssetService) saveAssetIndex(aindex AssetIndex) error {
	_, err := s.Sql.Exec(ASSET_INDEX_INSERT_UPDATE, aindex.systemId, aindex.name, aindex.fileIndex, aindex.uploadTime, aindex.fileIndex, aindex.uploadTime, aindex.systemId, aindex.name)
	if err != nil {
		return err
	}
	return nil
}
