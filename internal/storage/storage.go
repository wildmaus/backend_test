package storage

type Storage interface {
	User() UserRepository
	Transaction() TxRepository
	Details() DetailsRepository
}
