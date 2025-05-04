package persistence

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const DATABASE_URL = "postgres://postgres:password@192.168.1.7:5432/homing_agent"
//var Pool pgxpool.Pool
var Message string;
func init(){
	fmt.Printf("Init app %s\n","postgrsql")
	Message = "test me"
}

func Pool() error{
	pool, err := pgxpool.New(context.Background(),DATABASE_URL)
	if err != nil {
		return err
	}
	defer pool.Close()
	fmt.Printf("pooled %s\n","tested")
	rows, err := pool.Query(context.Background(), "SELECT name, host FROM agents")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var host string
		err := rows.Scan(&name, &host)
		if err != nil {
			return err
		}
		fmt.Printf("Data loaded %s >> %s\n", name, host)
	}
	return nil
}


func Start() error {
	conn, err := pgx.Connect(context.Background(), DATABASE_URL)
	if err != nil {
		fmt.Printf("Failed to connection")
		return err
	}
	defer conn.Close(context.Background())
	fmt.Printf("connnected %s\n", conn.Config().Host)
	rows, err := conn.Query(context.Background(), "SELECT name, host FROM agents")
	if err != nil {
		return err
	}
	defer rows.Close()
	fmt.Printf("query %s\n", conn.Config().User)

	for rows.Next() {
		var name string
		var host string
		err := rows.Scan(&name, &host)
		if err != nil {
			return err
		}
		fmt.Printf("Data loaded %s >> %s\n", name, host)
	}
	return nil
}
