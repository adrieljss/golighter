package platform

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DbName          string
	PoolSize        int
	PoolMaxIdleTime time.Duration
}

func ConnectDB(config PgConfig) *pgxpool.Pool {
	pgxConfig, err := pgxpool.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%d/%s", config.User, config.Password, config.Host, config.Port, config.DbName))

	if err != nil {
		log.Fatal("error parsing database configuration: ", err)
	}

	db, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}

	// try to ping the database until it's ready
	attempts := 5
	sleepTime := 5 * time.Second
	for i := 0; i < attempts; i++ {
		err = db.Ping(context.Background())
		if err != nil {
			log.Errorf("failed to ping db, trying again in 5 seconds... [attempt %d/%d]", i+1, attempts)
			if i == attempts-1 {
				log.Fatalf("failed to ping db after %d attempts", attempts)
			}
		} else {
			break
		}
		time.Sleep(sleepTime)
	}

	log.Info("connected to the database")

	return db
}
