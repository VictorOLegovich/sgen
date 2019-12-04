package chemical_storage

import (
	"context"
	"github.com/jackc/pgx"
	. "github.com/victorolegovich/sgen/test_data"
	qb "github.com/victorolegovich/sgen/test_data/database/general/query_builder"
)

type ChemicalStorage struct {
	db *pgx.Conn
	qb *qb.QueryBuilder
}

func NewChemicalStorage(db *pgx.Conn) *ChemicalStorage {
	qBuilder := qb.NewQueryBuilder("chemical", "postgresql")

	//you can opt out of using this action
	qBuilder.InitSets(updateSet, insertSet, selectSet)

	return &ChemicalStorage{db: db, qb: qBuilder}
}

func (s *ChemicalStorage) Create(chemical Chemical) error {
	query := s.qb.Insert().SQLString()

	_, err := s.db.Exec(
		context.Background(), query,
		&chemical.ID,
		&chemical.Element,
		&chemical.Position,
		&chemical.BackgroundState,
	)

	return err
}

func (s *ChemicalStorage) ReadOne(ID int) (*Chemical, error) {
	c := &Chemical{}

	query := s.qb.Select(false).Where("ID", "=").Limit(1).SQLString()

	err := s.db.QueryRow(
		context.Background(), query, ID).Scan(
		&c.ID,
		&c.Element,
		&c.Position,
		&c.BackgroundState,
	)

	return c, err
}

func (s *ChemicalStorage) ReadList() ([]Chemical, error) {
	var cList []Chemical
	c := Chemical{}

	query := s.qb.Select(false).Limit(10)

	rows, err := s.db.Query(context.Background(), query.SQLString())

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err := rows.Scan(
			&c.ID,
			&c.Element,
			&c.Position,
			&c.BackgroundState,
		)

		if err != nil {
			return nil, err
		}

		cList = append(cList, c)
	}

	return cList, nil
}

func (s *ChemicalStorage) Update(chemical Chemical) error {
	query := s.qb.Update().Where("ID", "=").SQLString()

	_, err := s.db.Exec(
		context.Background(), query,
		&chemical.ID,
		&chemical.Element,
		&chemical.Position,
		&chemical.BackgroundState,
	)

	return err
}

func (s *ChemicalStorage) Delete(ID int) error {
	query := s.qb.Delete().Where("ID", "=").SQLString()
	_, err := s.db.Exec(context.Background(), query, ID)

	return err
}
