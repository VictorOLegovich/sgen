package jet_storage

import (
	"context"
	"github.com/jackc/pgx"
	. "github.com/victorolegovich/sgen/data"
	qb "github.com/victorolegovich/sgen/data/db/general/query_builder"
)

type JetStorage struct {
	db *pgx.Conn
	qb *qb.QueryBuilder
}

func NewJetStorage(db *pgx.Conn) *JetStorage {
	qBuilder := qb.NewQueryBuilder("jet", "postgresql")
	qBuilder.InitSets(updateSet, insertSet, selectSet)

	return &JetStorage{db: db, qb: qBuilder}
}

func (s *JetStorage) Create(jet Jet) error {
	var err error

	query := s.qb.Insert().SQLString()

	_, err = s.db.Exec(query, &jet.DataID, &jet.NumID, &jet.Frazer)
	if err != nil {
		return err
	}

	return err
}

func (s *JetStorage) One(ID int) (*Jet, error) {
	var err error
	j := &Jet{}

	query := s.qb.Select(qb.Set).Where("ID", "=").Limit(1).SQLString()

	err = s.db.QueryRow(
		query, ID).Scan(&j.ID, &j.DataID, &j.NumID, &j.Frazer)

	return j, err
}

func (s *JetStorage) List() ([]*Jet, error) {
	var err error
	var jList []*Jet
	j := &Jet{}

	query := s.qb.Select(qb.Set).Limit(10).SQLString()

	rows, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(&j.ID, &j.DataID, &j.NumID, &j.Frazer); err != nil {
			return nil, err
		}

		jList = append(jList, j)
	}

	return jList, nil
}

func (s *JetStorage) Update(jet Jet) error {
	var err error

	query := s.qb.Update(qb.Set).Where("ID", "=").SQLString()

	_, err = s.db.Exec(query, &jet.ID, &jet.DataID, &jet.NumID, &jet.Frazer)

	return err
}

func (s *JetStorage) Delete(jet Jet) error {
	var err error
	query := s.qb.Delete().Where("ID", "=").SQLString()

	_, err = s.db.Exec(query, jet.ID)

	return err
}

func (s *JetStorage) GetJetByDataID(DataID int) (Jet, error) {
	var err error
	j := &Jet{}

	query := s.qb.Select(qb.Set).Where("data_id", "=").Limit(1).SQLString()

	err = s.db.QueryRow(
		query, DataID).Scan(&j.ID, &j.DataID, &j.NumID, &j.Frazer)

	return *j, err
}

func (s *JetStorage) GetJetByNumID(NumID int) (Jet, error) {
	var err error
	j := &Jet{}

	query := s.qb.Select(qb.Set).Where("num_id", "=").Limit(1).SQLString()

	err = s.db.QueryRow(
		query, NumID).Scan(&j.ID, &j.DataID, &j.NumID, &j.Frazer)

	return *j, err
}
