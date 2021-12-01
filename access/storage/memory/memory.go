package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/seftomsk/abf/access/storage"
)

type InMemory struct {
	mu        sync.Mutex
	whiteList map[string]map[string]struct{}
	blackList map[string]map[string]struct{}
}

func (s *InMemory) AddToWList(
	ctx context.Context,
	ip storage.IPEntity) error {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return ctxErr
	}

	if s.whiteList == nil {
		return storage.ErrInvalidInitialization
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if !storage.ValidEntity(ip) {
		return storage.ErrInvalidEntity
	}

	if s.whiteList[ip.Mask()] == nil {
		s.whiteList[ip.Mask()] = make(map[string]struct{})
	}

	s.whiteList[ip.Mask()][ip.IP()] = struct{}{}

	return nil
}

func (s *InMemory) AddToBList(
	ctx context.Context,
	ip storage.IPEntity) error {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return ctxErr
	}

	if s.blackList == nil {
		return storage.ErrInvalidInitialization
	}

	if !storage.ValidEntity(ip) {
		return storage.ErrInvalidEntity
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.blackList[ip.Mask()] == nil {
		s.blackList[ip.Mask()] = make(map[string]struct{})
	}

	s.blackList[ip.Mask()][ip.IP()] = struct{}{}

	return nil
}

func (s *InMemory) DeleteFromWList(
	ctx context.Context,
	ip storage.IPEntity) error {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return ctxErr
	}

	if s.whiteList == nil {
		return storage.ErrInvalidInitialization
	}

	if !storage.ValidEntity(ip) {
		return storage.ErrInvalidEntity
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.whiteList[ip.Mask()] == nil {
		return fmt.Errorf(
			"deleteFromWList - mask: %w",
			storage.ErrNotFound)
	}

	if _, ok := s.whiteList[ip.Mask()][ip.IP()]; !ok {
		return fmt.Errorf(
			"deleteFromWList - ip: %w",
			storage.ErrNotFound)
	}

	delete(s.whiteList[ip.Mask()], ip.IP())

	if len(s.whiteList[ip.Mask()]) == 0 {
		delete(s.whiteList, ip.Mask())
	}

	return nil
}

func (s *InMemory) DeleteFromBList(
	ctx context.Context,
	ip storage.IPEntity) error {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return ctxErr
	}

	if s.blackList == nil {
		return storage.ErrInvalidInitialization
	}

	if !storage.ValidEntity(ip) {
		return storage.ErrInvalidEntity
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.blackList[ip.Mask()] == nil {
		return fmt.Errorf(
			"deleteFromBList - mask: %w",
			storage.ErrNotFound)
	}

	if _, ok := s.blackList[ip.Mask()][ip.IP()]; !ok {
		return fmt.Errorf(
			"deleteFromBList - ip: %w",
			storage.ErrNotFound)
	}

	delete(s.blackList[ip.Mask()], ip.IP())

	if len(s.blackList[ip.Mask()]) == 0 {
		delete(s.blackList, ip.Mask())
	}

	return nil
}

func (s *InMemory) IsInBList(
	ctx context.Context,
	ip storage.IPEntity) (bool, error) {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return false, ctxErr
	}

	if s.blackList == nil {
		return false, storage.ErrInvalidInitialization
	}

	if !storage.ValidEntity(ip) {
		return false, storage.ErrInvalidEntity
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.blackList[ip.Mask()] == nil {
		return false, fmt.Errorf(
			"isInBList - mask: %w",
			storage.ErrNotFound)
	}

	if _, ok := s.blackList[ip.Mask()][ip.IP()]; !ok {
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

	if s.whiteList == nil {
		return false, storage.ErrInvalidInitialization
	}

	if !storage.ValidEntity(ip) {
		return false, storage.ErrInvalidEntity
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.whiteList[ip.Mask()] == nil {
		return false, fmt.Errorf(
			"isInWList - mask: %w",
			storage.ErrNotFound)
	}

	if _, ok := s.whiteList[ip.Mask()][ip.IP()]; !ok {
		return false, fmt.Errorf(
			"isInWList - ip: %w",
			storage.ErrNotFound)
	}

	return true, nil
}

func Create() *InMemory {
	return &InMemory{
		whiteList: make(map[string]map[string]struct{}),
		blackList: make(map[string]map[string]struct{}),
	}
}
