package psql

import (
	"backend_test/internal/storage"
	"backend_test/pkg/client/postgressql"
)

type psqlstorage struct {
	client            postgressql.Client
	userRepository    *UserRepository
	txRepository      *TxRepository
	detailsRepository *DetailsRepository
}

func NewStorage(client postgressql.Client) *psqlstorage {
	return &psqlstorage{
		client: client,
	}
}

func (s *psqlstorage) User() storage.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}
	s.userRepository = &UserRepository{s.client}
	return s.userRepository
}

func (s *psqlstorage) Transaction() storage.TxRepository {
	if s.txRepository != nil {
		return s.txRepository
	}
	s.txRepository = &TxRepository{client: s.client}
	return s.txRepository
}

func (s *psqlstorage) Details() storage.DetailsRepository {
	if s.detailsRepository != nil {
		return s.detailsRepository
	}
	s.detailsRepository = &DetailsRepository{client: s.client}
	return s.detailsRepository
}
