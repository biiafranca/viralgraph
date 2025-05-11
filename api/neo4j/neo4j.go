package neo4j

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/config"
)

var Driver neo4j.DriverWithContext

func init() {

	if os.Getenv("DOCKER_ENV") != "true" {
		envErr := godotenv.Load("../.env")
		if envErr != nil {
			log.Fatal("Error loading .env file")
		}
	}

	uri := os.Getenv("NEO4J_URI")
	user := os.Getenv("NEO4J_USER")
	password := os.Getenv("NEO4J_PASSWORD")

	driver, driverErr := neo4j.NewDriverWithContext(
		uri,
		neo4j.BasicAuth(user, password, ""),
		func(config *config.Config) {
			config.SocketConnectTimeout = 5 * time.Second
			config.MaxConnectionPoolSize = 10 // limit threads
		},
	)
	if driverErr != nil {
		log.Fatalf("Failed to connect to Neo4j: %v", driverErr)
	}

	Driver = driver
}

func GetSession() neo4j.SessionWithContext {
	return Driver.NewSession(context.Background(), neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
}
