package persistence

import (
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
)
func loop(rows pgx.Rows) error{
	var name string
	var hash string
	err := rows.Scan(&name, &hash)
	if err != nil {
		return err
	}
	fmt.Printf("Data loaded %s >> %s\n", name, hash)
	return nil
}
func TestPool(t *testing.T) {
	pool := Postgresql{Url: "postgres://postgres:password@192.168.1.7:5432/tarantula_user"}
	err := pool.Create()
	defer pool.Close()
	if err != nil {
		t.Errorf("SQL error %s", err.Error())
	}
	query := "SELECT name, hash FROM login WHERE name=$1 and system_id=$2"
	pool.Query(loop,query, "root", 1)
	pool.Query(loop,query, "root", 1)
	pool.Query(loop,query, "root", 1)
	pool.Query(loop,query, "root", 1)

	query1 := "SELECT name, hash FROM login"
	pool.Query(loop,query1)

	ct, _ := pool.Exec("INSERT INTO login (name,hash,system_id) VALUES($1,$2,$3)","adrun","hash",4)
	fmt.Printf("CT %d : ",ct)
}
