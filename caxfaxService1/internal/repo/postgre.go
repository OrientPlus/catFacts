package repo

import (
	"caxfaxService1/internal/entity"
	"caxfaxService1/internal/repo/sqlc"
	"caxfaxService1/pkg/logger"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type RepoImpl struct {
	db     *sql.DB
	pg     *sqlc.Queries
	logger logger.Logger
}

func NewRepo(config entity.Config) (Repo, error) {
	repo := RepoImpl{}
	repo.logger = config.Logger

	connectionURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		config.DbUser, config.DbPassword, config.DbHost, config.DbPort, config.DbName)

	db, err := sql.Open("postgres", connectionURL)
	if err != nil {
		repo.logger.Errorf("не удалось подключиться к БД: %s", err.Error())
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		repo.logger.Errorf("проверка соединения с БД провалилась: %s", err.Error())
		return nil, err
	}

	repo.db = db
	repo.pg = sqlc.New(db)

	return &repo, nil
}

func (repo *RepoImpl) Close() {
	repo.db.Close()
}

func (repo *RepoImpl) AddFact(ctx context.Context, fact entity.Fact) (int32, error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		repo.logger.Errorf("ошибка создания транзакции %s", err.Error())
		return 0, err
	}

	id, err := repo.pg.WithTx(tx).AddFact(context.Background(), sqlc.AddFactParams{
		Message: sql.NullString{
			String: fact.Message,
			Valid:  true,
		},
		Length: sql.NullInt32{
			Int32: fact.Length,
			Valid: true,
		},
	})
	if err != nil {
		repo.logger.Errorf("ошибка добавления записи в БД: %s", err.Error())
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		repo.logger.Errorf("ошибка завершения транзакции %s", err.Error())
		return 0, err
	}

	return id, nil
}

func (repo *RepoImpl) GetFact(ctx context.Context, factID int32) (entity.Fact, error) {
	factDTO, err := repo.pg.GetFactByID(ctx, factID)
	if err != nil {
		repo.logger.Errorf("ошибка получения данных из БД: %s", err.Error())
		return entity.Fact{}, err
	}

	return entity.Fact{
		Message: factDTO.Message.String,
		Length:  factDTO.Length.Int32,
	}, nil
}

func (repo *RepoImpl) UpdateFact(ctx context.Context, factID int32, fact entity.Fact) error {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		repo.logger.Errorf("ошибка создания транзакции %s", err.Error())
		return err
	}

	err = repo.pg.WithTx(tx).UpdateFactByID(ctx, sqlc.UpdateFactByIDParams{
		Message: sql.NullString{
			String: fact.Message,
			Valid:  true,
		},
		Length: sql.NullInt32{
			Int32: fact.Length,
			Valid: true,
		},
		ID: factID,
	})
	if err != nil {
		repo.logger.Errorf("ошибка обновления записи в БД: %s", err.Error())
		return err
	}

	err = tx.Commit()
	if err != nil {
		repo.logger.Errorf("ошибка завершения транзакции %s", err.Error())
		return err
	}

	return nil
}

func (repo *RepoImpl) DeleteFact(ctx context.Context, factID int32) error {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		repo.logger.Errorf("ошибка создания транзакции %s", err.Error())
		return err
	}

	err = repo.pg.WithTx(tx).DeleteFactByID(ctx, factID)
	if err != nil {
		repo.logger.Errorf("ошибка удаления записи из БД: %s", err.Error())
		return err
	}

	err = tx.Commit()
	if err != nil {
		repo.logger.Errorf("ошибка завершения транзакции %s", err.Error())
		return err
	}

	return nil
}
