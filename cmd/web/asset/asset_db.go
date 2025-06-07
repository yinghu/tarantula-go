package main

import (
	"fmt"

	"gameclustering.com/internal/bootstrap"
)

const (
	ASSET_INDEX_SQL_SCHEMA    string = "CREATE TABLE IF NOT EXISTS asset_index (system_id BIGINT NOT NULL,name VARCHAR(100) NOT NULL,file_index VARCHAR(100) NOT NULL,upload_time TIMESTAMP DEFAULT NOW(),PRIMARY KEY(system_id,name))"
	ASSET_INDEX_INSERT_UPDATE string = "INSERT INTO asset_index AS a (system_id,name,file_index) VALUES($1,$2,$3) ON CONFLICT (system_id,name) DO UPDATE SET a.file_index = $4 WHERE a.system_id = $5 AND a.name = $6;"
)

func (s *AssetService) createSchema() error {
	_, err := s.Sql.Exec(bootstrap.METRICS_SQL_SCHEMA)
	if err != nil {
		return err
	}
	_, err = s.Sql.Exec(ASSET_INDEX_SQL_SCHEMA)
	return err
}

func (s *AssetService) saveAssetIndex(aindex AssetIndex) error {
	rt, err := s.Sql.Exec(ASSET_INDEX_INSERT_UPDATE, aindex.systemId, aindex.name, aindex.fileIndex, aindex.fileIndex, aindex.systemId, aindex.name)
	fmt.Printf("CHANGE %d\n", rt)
	if err != nil {
		fmt.Printf("Erro : %s\n", err.Error())
	}
	return err
}
