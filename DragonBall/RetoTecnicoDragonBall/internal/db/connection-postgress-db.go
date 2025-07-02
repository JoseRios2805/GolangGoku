package db

import (
	"RetoTecnicoDragonBall/internal/logs"
	"database/sql"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"net/url"
	"os"
	"strconv"
	"sync"
)

type IConnectionPostgresBServer interface {
	GetConnection() (*gorm.DB, error)
	Close() error
}

type connectionPostgres struct {
	log       logs.ILogger
	instance  *gorm.DB
	sqlDB     *sql.DB
	initError error
	once      sync.Once
}

func NewPostgresConnection(lg logs.ILogger) IConnectionPostgresBServer {
	return &connectionPostgres{
		log:  lg,
		once: sync.Once{},
	}
}

func (c *connectionPostgres) GetConnection() (*gorm.DB, error) {

	// Inicialización única
	c.once.Do(func() {
		c.instance, c.initError = c.connectionPostgresDBServer()
		if c.instance != nil {
			var err error
			c.sqlDB, err = c.instance.DB()
			if err != nil {
				c.log.Error("Error to get DB", err)
				c.initError = err
			}
		}
	})

	// Verificar si la conexión sigue activa
	if err := c.healthCheckAndReconnect(); err != nil {
		return nil, err
	}

	if c.initError != nil || c.instance == nil {
		msg := "database connection failed to initialize"
		if c.initError != nil {
			msg = fmt.Sprintf("failed to initialize the connection: %v", c.initError)
		}
		return nil, fmt.Errorf(msg)
	}

	return c.instance, nil
}

func (c *connectionPostgres) Close() error {
	if c.sqlDB != nil {
		c.log.Info("Closing Postgrest-DB connection")
		return c.sqlDB.Close()
	}
	return nil
}

func (c *connectionPostgres) healthCheckAndReconnect() error {
	// Verificar si la conexión sigue activa
	if c.instance != nil && c.sqlDB != nil && c.initError == nil {
		if err := c.sqlDB.Ping(); err != nil {
			c.log.Error("Connection lost, attempting to reconnect", err)

			if closeErr := c.Close(); closeErr != nil {
				c.log.Error("Error closing previous connection", closeErr)
			}
			// Reconectar
			newInstance, newErr := c.connectionPostgresDBServer()
			if newErr != nil {
				c.log.Error("Error reconnecting", newErr)
				return newErr
			}
			c.instance = newInstance
			c.sqlDB, _ = c.instance.DB()
			c.initError = nil
		}
	}

	return nil
}

func (c *connectionPostgres) connectionPostgresDBServer() (*gorm.DB, error) {
	query := url.Values{}
	query.Add("database", os.Getenv("DATABASE_NAME"))

	portDBStr := os.Getenv("PORT_DB")

	portDB, _ := strconv.Atoi(portDBStr)
	// Build the DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("HOST_DB"), portDB, os.Getenv("USER_DB"), os.Getenv("PASS_DB"), os.Getenv("DATABASE_NAME"))

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		SkipDefaultTransaction: true,
	})

	if err != nil {
		c.log.Error("Error open postgres-DB connection", err)
		return nil, err
	}

	// Get the underlying sql.DB object
	sqlDB, err := db.DB()
	if err != nil {
		c.log.Error("Error to get postgres.DB", err)
		return nil, err
	}

	// Configuración del pool de conexiones
	sqlDB.SetMaxOpenConns(12) // Número máximo de conexiones abiertas simultáneamente
	sqlDB.SetMaxIdleConns(5)  // Número máximo de conexiones inactivas que pueden estar en el pool

	return db, nil
}
