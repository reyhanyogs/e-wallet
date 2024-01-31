package component

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/reyhanyogs/e-wallet/internal/config"
)

func GetDatabaseConn(config *config.Config) *sql.DB {
	dsn := fmt.Sprintf(
		"host=%s "+
			"port=%s "+
			"user=%s "+
			"password=%s "+
			"dbname=%s "+
			"sslmode=disable",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.Name,
	)

	connection, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("error when open connection %s", err.Error())
	}

	err = connection.Ping()
	if err != nil {
		log.Fatalf("error when open connection %s", err.Error())
	}

	return connection
}
