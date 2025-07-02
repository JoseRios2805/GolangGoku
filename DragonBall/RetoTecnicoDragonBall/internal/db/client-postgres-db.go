package db

import (
	"RetoTecnicoDragonBall/internal/logs"
	"RetoTecnicoDragonBall/internal/utils"
	"gorm.io/gorm"
)

type IClientPostgres interface {
	GenericClientService(operation func(*gorm.DB) error) error
	WithTransaction(operation func(*gorm.DB) error) error
}

type clientPostgres struct {
	postgresDB IConnectionPostgresBServer
	log        logs.ILogger
}

func NewClientPostgres(postgresCon IConnectionPostgresBServer, lg logs.ILogger) IClientPostgres {
	return &clientPostgres{
		postgresDB: postgresCon,
		log:        lg,
	}
}

func (c *clientPostgres) GenericClientService(operation func(*gorm.DB) error) error {
	db, err := c.postgresDB.GetConnection()
	if err != nil {
		c.log.Error(utils.MSG_ERROR_POSTGRES_CONNECTION, err)
		return err
	}

	if err = operation(db); err != nil {
		c.log.Error("Error during database operation", err)
		return err
	}

	return nil
}

func (c *clientPostgres) WithTransaction(fn func(*gorm.DB) error) error {
	db, err := c.postgresDB.GetConnection()
	if err != nil {
		c.log.Error(utils.MSG_ERROR_POSTGRES_CONNECTION, err)
		return err
	}

	tx := db.Begin()
	if tx.Error != nil {
		c.log.Error("Error starting transaction", tx.Error)
		return tx.Error
	}

	err = fn(tx)
	if err != nil {
		c.log.Error("Error during transaction operation", err)
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
