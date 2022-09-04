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

func NewDBClient(conf *config.Config) (*DBClient, error) {
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

	return &DBClient{conn: conn}, nil
}

func (c *DBClient) InsertEvent(metrics string) {
	log.Printf("Inserting event: %v", metrics)
	insertStmt := fmt.Sprintf("INSERT INTO container_bridge FORMAT JSONAsObject %v", metrics)
	if err := c.conn.Exec(context.Background(), insertStmt); err != nil {
		log.Printf("Insert failed, %v", err)
	}
}

func (c *DBClient) Close() {
	_ = c.conn.Close()
}
