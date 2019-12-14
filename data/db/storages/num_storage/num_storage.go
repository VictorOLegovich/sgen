package num_storage

import (
	"context"
	"github.com/jackc/pgx"
	. "github.com/victorolegovich/sgen/data"
	qb "github.com/victorolegovich/sgen/data/db/general/query_builder"
	"github.com/victorolegovich/sgen/data/db/storages/jet_storage"
)

type NumStorage struct {
	db *pgx.Conn
	qb *qb.QueryBuilder
}

func NewNumStorage(db *pgx.Conn) *NumStorage {
	qBuilder := qb.NewQueryBuilder("num", "postgresql")
	qBuilder.InitSets(updateSet, insertSet, selectSet)

	return &NumStorage{db: db, qb: qBuilder}
}

func (s *NumStorage) Create(num Num) error {
	var err error

	query := s.qb.Insert().SQLString()
	_, err = s.db.Exec(query, &num.DataID)

	jetStorage := jet_storage.NewJetStorage(nil)

	if err = jetStorage.Create(num.Jet); err != nil {
		return err
	}

	return err
}

func (s *NumStorage) one(ID int) (*Num, error) {
	var err error
	n := &Num{}

	query := s.qb.Select(qb.Set).Where("ID", "=").Limit(1).SQLString()

	err = s.db.QueryRow(
		query, ID).Scan(&n.ID, &n.DataID)

	return n, err
}

func (s *NumStorage) list() ([]*Num, error) {
	var err error
	var nList []*Num
	n := &Num{}

	query := s.qb.Select(qb.Set).Limit(10).SQLString()

	rows, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(&n.ID, &n.DataID); err != nil {
			return nil, err
		}

		nList = append(nList, n)
	}

	return nList, nil
}

func (s *NumStorage) Update(num Num) error {
	var err error

	query := s.qb.Update(qb.Set).Where("ID", "=").SQLString()
	_, err = s.db.Exec(query, &num.ID, &num.DataID)

	jetStorage := jet_storage.NewJetStorage(nil)

	if err = jetStorage.Update(num.Jet); err != nil {
		return err
	}

	return err
}

func (s *NumStorage) Delete(num Num) error {
	var err error

	query := s.qb.Delete().Where("ID", "=").SQLString()
	_, err = s.db.Exec(query, &num.ID, &num.DataID)

	jetStorage := jet_storage.NewJetStorage(nil)

	if err = jetStorage.Delete(num.Jet); err != nil {
		return err
	}

	return err
}

func (s *NumStorage) GetNumByDataID(DataID int) (Num, error) {
	var err error
	n := &Num{}

	query := s.qb.Select(qb.Set).Where("data_id", "=").Limit(1).SQLString()

	err = s.db.QueryRow(
		query, DataID).Scan(&n.ID, &n.DataID)

	return *n, err
}

func (s *NumStorage) One(NumID int) (*Num, error) {
	num, err := s.one(NumID)
	if err != nil {
		return nil, err
	}

	jetStorage := jet_storage.NewJetStorage(nil)

	if num.Jet, err = jetStorage.GetJetByNumID(NumID); err != nil {
		return nil, err
	}

	return num, nil
}

func (s *NumStorage) List() ([]*Num, error) {
	numList, err := s.list()
	if err != nil {
		return nil, err
	}

	jetStorage := jet_storage.NewJetStorage(nil)

	for _, num := range numList {
		if num.Jet, err = jetStorage.GetJetByNumID(num.ID); err != nil {
			return nil, err
		}
	}

	return numList, nil
}
