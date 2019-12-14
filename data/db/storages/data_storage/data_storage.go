package data_storage

import (
	"context"
	"github.com/jackc/pgx"
	. "github.com/victorolegovich/sgen/data"
	qb "github.com/victorolegovich/sgen/data/db/general/query_builder"
	"github.com/victorolegovich/sgen/data/db/storages/jet_storage"
	"github.com/victorolegovich/sgen/data/db/storages/num_storage"
)

type DataStorage struct {
	db *pgx.Conn
	qb *qb.QueryBuilder
}

func NewDataStorage(db *pgx.Conn) *DataStorage {
	qBuilder := qb.NewQueryBuilder("data", "postgresql")
	qBuilder.InitSets(updateSet, insertSet, selectSet)

	return &DataStorage{db: db, qb: qBuilder}
}

func (s *DataStorage) Create(data Data) error {
	var err error

	query := s.qb.Insert().SQLString()
	_, err = s.db.Exec(query, &data.Family)

	numStorage := num_storage.NewNumStorage(nil)
	jetStorage := jet_storage.NewJetStorage(nil)

	if err = numStorage.Create(data.Num); err != nil {
		return err
	}

	if err = jetStorage.Create(data.Jet); err != nil {
		return err
	}

	return err
}

func (s *DataStorage) one(ID int) (*Data, error) {
	var err error
	d := &Data{}

	query := s.qb.Select(qb.Set).Where("ID", "=").Limit(1).SQLString()

	err = s.db.QueryRow(
		query, ID).Scan(&d.ID, &d.Family)

	return d, err
}

func (s *DataStorage) list() ([]*Data, error) {
	var err error
	var dList []*Data
	d := &Data{}

	query := s.qb.Select(qb.Set).Limit(10).SQLString()

	rows, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(&d.ID, &d.Family); err != nil {
			return nil, err
		}

		dList = append(dList, d)
	}

	return dList, nil
}

func (s *DataStorage) Update(data Data) error {
	var err error

	query := s.qb.Update(qb.Set).Where("ID", "=").SQLString()
	_, err = s.db.Exec(query, &data.ID, &data.Family)

	numStorage := num_storage.NewNumStorage(nil)
	jetStorage := jet_storage.NewJetStorage(nil)

	if err = numStorage.Update(data.Num); err != nil {
		return err
	}

	if err = jetStorage.Update(data.Jet); err != nil {
		return err
	}

	return err
}

func (s *DataStorage) Delete(data Data) error {
	var err error

	query := s.qb.Delete().Where("ID", "=").SQLString()
	_, err = s.db.Exec(query, &data.ID, &data.Family)

	numStorage := num_storage.NewNumStorage(nil)
	jetStorage := jet_storage.NewJetStorage(nil)

	if err = numStorage.Delete(data.Num); err != nil {
		return err
	}

	if err = jetStorage.Delete(data.Jet); err != nil {
		return err
	}

	return err
}

func (s *DataStorage) One(DataID int) (*Data, error) {
	data, err := s.one(DataID)
	if err != nil {
		return nil, err
	}

	numStorage := num_storage.NewNumStorage(nil)
	jetStorage := jet_storage.NewJetStorage(nil)

	if data.Num, err = numStorage.GetNumByDataID(DataID); err != nil {
		return nil, err
	}

	if data.Jet, err = jetStorage.GetJetByDataID(DataID); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *DataStorage) List() ([]*Data, error) {
	dataList, err := s.list()
	if err != nil {
		return nil, err
	}

	numStorage := num_storage.NewNumStorage(nil)
	jetStorage := jet_storage.NewJetStorage(nil)

	for _, data := range dataList {
		if data.Num, err = numStorage.GetNumByDataID(data.ID); err != nil {
			return nil, err
		}
		if data.Jet, err = jetStorage.GetJetByDataID(data.ID); err != nil {
			return nil, err
		}
	}

	return dataList, nil
}
