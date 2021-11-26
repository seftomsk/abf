package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/seftomsk/abf/access/storage"
)

const (
	white = "white"
	black = "black"
)

type InMemory struct {
	mu         sync.Mutex
	collection map[string]map[string]map[string]struct{}
}

func (s *InMemory) AddToWList(
	ctx context.Context,
	ip storage.IPEntity) error {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return ctxErr
	}

	if s.collection[white] == nil {
		return storage.ErrInvalidInitialization
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if !storage.ValidEntity(ip) {
		return storage.ErrInvalidEntity
	}

	if s.collection[white][ip.Mask()] == nil {
		s.collection[white][ip.Mask()] = make(map[string]struct{})
	}

	s.collection[white][ip.Mask()][ip.IP()] = struct{}{}

	return nil
}

func (s *InMemory) AddToBList(
	ctx context.Context,
	ip storage.IPEntity) error {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return ctxErr
	}

	if s.collection[black] == nil {
		return storage.ErrInvalidInitialization
	}

	if !storage.ValidEntity(ip) {
		return storage.ErrInvalidEntity
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.collection[black][ip.Mask()] == nil {
		s.collection[black][ip.Mask()] = make(map[string]struct{})
	}

	s.collection[black][ip.Mask()][ip.IP()] = struct{}{}

	return nil
}

func (s *InMemory) DeleteFromWList(
	ctx context.Context,
	ip storage.IPEntity) error {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return ctxErr
	}

	if s.collection[white] == nil {
		return storage.ErrInvalidInitialization
	}

	if !storage.ValidEntity(ip) {
		return storage.ErrInvalidEntity
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.collection[white][ip.Mask()] == nil {
		return fmt.Errorf(
			"deleteFromWList - mask: %w",
			storage.ErrNotFound)
	}

	if _, ok := s.collection[white][ip.Mask()][ip.IP()]; !ok {
		return fmt.Errorf(
			"deleteFromWList - ip: %w",
			storage.ErrNotFound)
	}

	delete(s.collection[white][ip.Mask()], ip.IP())

	if len(s.collection[white][ip.Mask()]) == 0 {
		delete(s.collection[white], ip.Mask())
	}

	return nil
}

func (s *InMemory) DeleteFromBList(
	ctx context.Context,
	ip storage.IPEntity) error {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return ctxErr
	}

	if s.collection[black] == nil {
		return storage.ErrInvalidInitialization
	}

	if !storage.ValidEntity(ip) {
		return storage.ErrInvalidEntity
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.collection[black][ip.Mask()] == nil {
		return fmt.Errorf(
			"deleteFromBList - mask: %w",
			storage.ErrNotFound)
	}

	if _, ok := s.collection[black][ip.Mask()][ip.IP()]; !ok {
		return fmt.Errorf(
			"deleteFromBList - ip: %w",
			storage.ErrNotFound)
	}

	delete(s.collection[black][ip.Mask()], ip.IP())

	if len(s.collection[black][ip.Mask()]) == 0 {
		delete(s.collection[black], ip.Mask())
	}

	return nil
}

func (s *InMemory) IsInBList(
	ctx context.Context,
	ip storage.IPEntity) (bool, error) {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return false, ctxErr
	}

	if s.collection[black] == nil {
		return false, storage.ErrInvalidInitialization
	}

	if !storage.ValidEntity(ip) {
		return false, storage.ErrInvalidEntity
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.collection[black][ip.Mask()] == nil {
		return false, fmt.Errorf(
			"isInBList - mask: %w",
			storage.ErrNotFound)
	}

	if _, ok := s.collection[black][ip.Mask()][ip.IP()]; !ok {
		return false, fmt.Errorf(
			"isInBList - ip: %w",
			storage.ErrNotFound)
	}

	return true, nil
}

func (s *InMemory) IsInWList(
	ctx context.Context,
	ip storage.IPEntity) (bool, error) {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return false, ctxErr
	}

	if s.collection[white] == nil {
		return false, storage.ErrInvalidInitialization
	}

	if !storage.ValidEntity(ip) {
		return false, storage.ErrInvalidEntity
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.collection[white][ip.Mask()] == nil {
		return false, fmt.Errorf(
			"isInWList - mask: %w",
			storage.ErrNotFound)
	}

	if _, ok := s.collection[white][ip.Mask()][ip.IP()]; !ok {
		return false, fmt.Errorf(
			"isInWList - ip: %w",
			storage.ErrNotFound)
	}

	return true, nil
}

func Create() *InMemory {
	collection := make(map[string]map[string]map[string]struct{})
	collection[white] = make(map[string]map[string]struct{})
	collection[black] = make(map[string]map[string]struct{})

	return &InMemory{collection: collection}
}
