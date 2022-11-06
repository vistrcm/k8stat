package main

import "github.com/dgraph-io/badger/v3"

type storage struct {
	*badger.DB
}

func (s *storage) Save(name string, stat *Stat) error {
	err := s.Update(func(tnx *badger.Txn) error {
		// TODO: need to save whole structure
		value, err := floatToBytes(stat.averageMem)
		if err != nil {
			return err
		}
		e := badger.NewEntry([]byte(name), value)
		err = tnx.SetEntry(e)
		return err
	})
	return err
}

func newStorage(path string) (*storage, error) {
	db, err := badger.Open(badger.DefaultOptions(path))
	return &storage{db}, err
}
