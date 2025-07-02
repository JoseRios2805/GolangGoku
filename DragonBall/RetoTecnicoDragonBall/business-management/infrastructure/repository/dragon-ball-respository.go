package repository

import (
	"RetoTecnicoDragonBall/business-management/domain/entity"
	"RetoTecnicoDragonBall/internal/db"
	"RetoTecnicoDragonBall/internal/logs"
	"gorm.io/gorm"
)

type IDragonBallRepository interface {
	Insert(source *entity.SaveCharacterRequestDTO) (id int, err error)
	SelectByNameInDB(name string) (data *entity.GetCharacterByNameDBResponseEntity, err error)
}

type dragonBallRepository struct {
	_log        logs.ILogger
	_postgresDB db.IClientPostgres
}

func NewDragonBallRepository(postgresDB db.IClientPostgres, log logs.ILogger) IDragonBallRepository {
	return &dragonBallRepository{
		_log:        log,
		_postgresDB: postgresDB,
	}
}

func (p *dragonBallRepository) Insert(source *entity.SaveCharacterRequestDTO) (id int, err error) {
	return id, p._postgresDB.GenericClientService(func(db *gorm.DB) error {
		sql, values := p.buildInsertQuery(source)
		return db.Raw(sql, values...).Scan(&id).Error
	})
}

func (p *dragonBallRepository) SelectByNameInDB(name string) (data *entity.GetCharacterByNameDBResponseEntity, err error) {
	query := `SELECT * FROM get_character_by_name(?)`

	return data, p._postgresDB.GenericClientService(func(db *gorm.DB) error {
		return db.Raw(query, name).Scan(&data).Error
	})
}

func (p *dragonBallRepository) buildInsertQuery(source *entity.SaveCharacterRequestDTO) (string, []interface{}) {
	sql := `SELECT save_character(?, ?, ?)`

	values := []interface{}{
		source.IdCharacter,
		source.Name,
		source.Description,
	}

	return sql, values
}
