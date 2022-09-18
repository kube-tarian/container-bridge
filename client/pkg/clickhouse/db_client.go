package clickhouse

import (
	"context"
	"fmt"
	"log"

	"github.com/kube-tarian/container-bridge/client/pkg/config"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type DBClient struct {
	conn driver.Conn
	conf *config.Config
}

type DBInterface interface {
	InsertEvent(event string)
	FetchEvents() []map[string]interface{}
	Close()
}

func NewDBClient(conf *config.Config) (DBInterface, error) {
	log.Println("Create DB if not exists")
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{conf.DBAddress},
	})
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}

	var settings []struct {
		Name        string  `ch:"name"`
		Value       string  `ch:"value"`
		Changed     uint8   `ch:"changed"`
		Description string  `ch:"description"`
		Min         *string `ch:"min"`
		Max         *string `ch:"max"`
		Readonly    uint8   `ch:"readonly"`
		Type        string  `ch:"type"`
	}
	if err = conn.Select(context.Background(), &settings, "SELECT * FROM system.settings WHERE name LIKE $1 ORDER BY length(name) LIMIT 5", "%max%"); err != nil {
		conn.Close()
		return nil, err
	}
	for _, s := range settings {
		fmt.Printf("name: %s, value: %s, type=%s\n", s.Name, s.Value, s.Type)
	}

	const dbCreate = `CREATE DATABASE IF NOT EXISTS bridge;`
	if err := conn.Exec(context.Background(), dbCreate); err != nil {
		return nil, err
	}
	_ = conn.Close()

	log.Println("Initializing DB client")
	conn, err = clickhouse.Open(&clickhouse.Options{
		Addr: []string{conf.DBAddress},
		Auth: clickhouse.Auth{
			Database: "bridge",
		},
		// Auth: clickhouse.Auth{
		// 	Database: "default",
		// 	Username: "default",
		// 	Password: "",
		// },
		// Compression: &clickhouse.Compression{
		// 	Method: clickhouse.CompressionLZ4,
		// },
		// Settings: clickhouse.Settings{
		// 	"max_execution_time": 60,
		// },
		//Debug: true,
	})
	if err != nil {
		return nil, err
	}

	const ddlSetExperimental = `SET allow_experimental_object_type=1;`
	if err := conn.Exec(context.Background(), ddlSetExperimental); err != nil {
		return nil, err
	}

	const ddlCreateTable = `CREATE table IF NOT EXISTS container_bridge(event JSON) ENGINE = MergeTree ORDER BY tuple();`
	if err := conn.Exec(context.Background(), ddlCreateTable); err != nil {
		return nil, err
	}

	return &DBClient{conn: conn, conf: conf}, nil
}

func (c *DBClient) InsertEvent(event string) {
	log.Printf("Inserting event: %v", event)
	insertStmt := fmt.Sprintf("INSERT INTO container_bridge FORMAT JSONAsObject %v", event)
	if err := c.conn.Exec(context.Background(), insertStmt); err != nil {
		log.Printf("Insert failed, %v", err)
	}
}

func (c *DBClient) FetchEvents() []map[string]interface{} {
	log.Printf("Fetching events")
	events := []map[string]interface{}{}
	insertStmt := "select event from container_bridge;"
	rows, err := c.conn.Query(context.Background(), insertStmt)
	if err != nil {
		log.Printf("Insert failed, %v", err)
	}
	for rows.Next() {
		event := map[string]interface{}{}
		if err := rows.Scan(&event); err != nil {
			log.Printf("Rows scan failed: %v", err)
			return events
		}
		fmt.Printf("row: event=%v\n", event)
		events = append(events, event)
	}
	rows.Close()
	if rows.Err() != nil {
		log.Printf("Fetching rows failed: %v", rows.Err())
	}
	return events
}

func (c *DBClient) Close() {
	_ = c.conn.Close()
}
