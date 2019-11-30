package testdata_storage

import (
	. "github.com/victorolegovich/sgen/test_data"
	qb "github.com/victorolegovich/sgen/test_data/database/general/query_builder"
)

type TestDataStorage struct {
	DataBase string
	qb       *qb.QueryBuilder
}

func NewTestDataStorage(DB string) *TestDataStorage {
	queryBuilder := qb.NewQueryBuilder("TestData", updateSet, insertSet, selectSet, "PostgreSQL")
	return &TestDataStorage{DataBase: DB, qb: queryBuilder}
}

func (s *TestDataStorage) Create(testdata TestData) error {
	query := s.qb.Insert()
	println(query)
	return nil
}

func (s *TestDataStorage) ReadOne(ID int) (TestData, error) {
	query := s.qb.Select(false).Where("ID", "=")
	println(query)
	return TestData{}, nil
}

func (s *TestDataStorage) ReadList() ([]TestData, error) {
	query := s.qb.Select(false)
	println(query)
	return []TestData{}, nil
}

func (s *TestDataStorage) Update(field, value string) error {
	query := s.qb.Update(field).Where("ID", "=")
	println(query)
	return nil
}

func (s *TestDataStorage) UpdateSeveral(testdata TestData) error {
	query := s.qb.UpdateSeveral().Where("ID", "=")
	println(query)
	return nil
}

func (s *TestDataStorage) Delete(ID int) error {
	query := s.qb.Delete().Where("ID", "=")
	println(query)
	return nil
}
