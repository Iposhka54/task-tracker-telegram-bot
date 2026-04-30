package storage

type Entity[ID comparable] interface {
	GetID() ID
}
