package memory

import (
	"exam-api/domain"
	"fmt"
	"sync"
)

// This lines checks if Store implements domain.Storage
// It will fail at build time if not
var _ domain.Storage = (*Store)(nil)

type Store struct {
	products map[string]domain.Product
	// We are using a Read-Write Mutex here
	// This guarantees us when we lock and unlock it that either
	// At most one goroutine is writing in the map and none are reading or;
	// No goroutine is writing and any number are reading
	mu sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		products: make(map[string]domain.Product),
		mu:       sync.RWMutex{},
	}
}

func (s *Store) Save(product domain.Product) (string, bool, error) {
	// Lock - writer's lock
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.products[product.GetHash()]
	if ok {
		return product.GetHash(), true, nil
	}
	s.products[product.GetHash()] = product

	return product.GetHash(), false, nil
}

func (s *Store) Get(id string) (domain.Product, bool, error) {
	// RLock - reader's lock
	s.mu.RLock()
	defer s.mu.RUnlock()

	p, ok := s.products[id]
	if !ok {
		return domain.Product{}, false, nil
	}
	return p, true, nil
}

func (s *Store) Update(id string, diff domain.Product) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.products[id]
	if !ok {
		return false, nil
	}

	newProduct := s.products[id]

	if diff.Name != newProduct.Name {
		newProduct.Name = diff.Name
	}
	if diff.Price != newProduct.Price {
		newProduct.Price = diff.Price
	}
	if diff.Stock != newProduct.Stock {
		newProduct.Stock = diff.Stock
	}

	s.products[id] = newProduct
	return true, nil
}

func (s *Store) Delete(id string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p, ok := s.products[id]
	fmt.Println(p)
	if !ok {
		return false, nil
	}

	delete(s.products, id)
	return true, nil
}
