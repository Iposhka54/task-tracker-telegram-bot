package storage

import (
	"errors"
	"fmt"
	"sync"
)

var (
	ErrAlreadyExists = errors.New("storage: entity already exists")
	ErrNotFound      = errors.New("storage: entity not found")
)

// Storage описывает базовое хранилище для сущностей T с идентификатором ID.
type Storage[T Entity[ID], ID comparable] interface {
	Add(entity T) error
	AddMany(entities []T) error
	GetByID(id ID) (T, error)
	Update(entity T) error
	DeleteByID(id ID) error
	DeleteByIDs(ids []ID) error
	GetAll() []T
}

type InMemoryStorage[T Entity[ID], ID comparable] struct {
	mu   sync.RWMutex
	data map[ID]T
}

func NewInMemoryStorage[T Entity[ID], ID comparable]() *InMemoryStorage[T, ID] {
	return &InMemoryStorage[T, ID]{
		data: make(map[ID]T),
	}
}

func (s *InMemoryStorage[T, ID]) Add(entity T) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := entity.GetID()
	if _, exists := s.data[id]; exists {
		return ErrAlreadyExists
	}

	s.data[id] = entity
	return nil
}

func (s *InMemoryStorage[T, ID]) AddMany(entities []T) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	seen := make(map[ID]struct{}, len(entities))
	for _, entity := range entities {
		id := entity.GetID()
		if _, exists := s.data[id]; exists {
			return fmt.Errorf("%w: %v", ErrAlreadyExists, id)
		}
		if _, exists := seen[id]; exists {
			return fmt.Errorf("%w: duplicate id in batch %v", ErrAlreadyExists, id)
		}
		seen[id] = struct{}{}
	}

	for _, entity := range entities {
		s.data[entity.GetID()] = entity
	}
	return nil
}

func (s *InMemoryStorage[T, ID]) GetByID(id ID) (T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entity, exists := s.data[id]
	if !exists {
		var zero T
		return zero, fmt.Errorf("%w: %v", ErrNotFound, id)
	}

	return entity, nil
}

func (s *InMemoryStorage[T, ID]) Update(entity T) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := entity.GetID()

	if _, exists := s.data[id]; !exists {
		return fmt.Errorf("%w: %v", ErrNotFound, id)
	}

	s.data[id] = entity
	return nil
}

func (s *InMemoryStorage[T, ID]) DeleteByID(id ID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[id]; !exists {
		return fmt.Errorf("%w: %v", ErrNotFound, id)
	}

	delete(s.data, id)
	return nil
}

func (s *InMemoryStorage[T, ID]) DeleteByIDs(ids []ID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, id := range ids {
		if _, exists := s.data[id]; !exists {
			return fmt.Errorf("%w: %v", ErrNotFound, id)
		}
	}

	for _, id := range ids {
		delete(s.data, id)
	}

	return nil
}

func (s *InMemoryStorage[T, ID]) GetAll() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entities := make([]T, 0, len(s.data))
	for _, entity := range s.data {
		entities = append(entities, entity)
	}

	return entities
}
