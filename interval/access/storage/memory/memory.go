package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/seftomsk/abf/interval/access/storage"
)

const (
	white = "white"
	black = "black"
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

func (s *InMemory) deleteFromList(
	ctx context.Context,
	ip storage.IPEntity,
	list string) error {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return ctxErr
	}

	if s.whiteList == nil || s.blackList == nil {
		return storage.ErrInvalidInitialization
	}

	if !storage.ValidEntity(ip) {
		return storage.ErrInvalidEntity
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	var currentList map[string]map[string]struct{}
	if list == white {
		currentList = s.whiteList
	} else {
		currentList = s.blackList
	}

	if currentList[ip.Mask()] == nil {
		return fmt.Errorf("mask: %w", storage.ErrNotFound)
	}

	if _, ok := currentList[ip.Mask()][ip.IP()]; !ok {
		return fmt.Errorf("ip: %w", storage.ErrNotFound)
	}

	delete(currentList[ip.Mask()], ip.IP())

	if len(currentList[ip.Mask()]) == 0 {
		delete(currentList, ip.Mask())
	}

	return nil
}

func (s *InMemory) DeleteFromWhiteList(
	ctx context.Context,
	ip storage.IPEntity) error {
	return s.deleteFromList(ctx, ip, white)
}

func (s *InMemory) DeleteFromBlackList(
	ctx context.Context,
	ip storage.IPEntity) error {
	return s.deleteFromList(ctx, ip, black)
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
