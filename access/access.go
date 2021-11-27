package access

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/google/uuid"

	"github.com/seftomsk/abf/access/storage"
)

var ErrInvalidStorage = errors.New("you must provide a storage")

var ErrNotFound = errors.New("not found")

var ErrEmptyIp = errors.New("empty ip")

type ErrParseIp struct {
	err error
	msg string
}

func (e *ErrParseIp) Error() string {
	return e.msg
}

func (e *ErrParseIp) Unwrap() error {
	return e.err
}

type IStorage interface {
	AddToWList(ctx context.Context, ip storage.IPEntity) error
	AddToBList(ctx context.Context, ip storage.IPEntity) error
	DeleteFromWList(ctx context.Context, ip storage.IPEntity) error
	DeleteFromBList(ctx context.Context, ip storage.IPEntity) error
	IsInWList(ctx context.Context, ip storage.IPEntity) (bool, error)
	IsInBList(ctx context.Context, ip storage.IPEntity) (bool, error)
}

type IPAccess struct {
	storage IStorage
}

type IpDTO struct {
	IP string
}

func NewIPAccess(storage IStorage) *IPAccess {
	return &IPAccess{storage: storage}
}

func (a *IPAccess) parseIPAndMask(ip string) (string, string, error) {
	if len(ip) <= 0 {
		return "", "", ErrEmptyIp
	}

	ipAddress, ipNet, err := net.ParseCIDR(ip)
	if err != nil {
		return "", "", &ErrParseIp{
			err: err,
			msg: "parseIPAndMask",
		}
	}

	byteMask := ipNet.Mask
	mask := fmt.Sprintf(
		"%d.%d.%d.%d",
		byteMask[0],
		byteMask[1],
		byteMask[2],
		byteMask[3])

	return ipAddress.String(), mask, nil
}

func (a *IPAccess) AddToWList(ctx context.Context, dto IpDTO) error {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return ctxErr
	}

	if a.storage == nil {
		return ErrInvalidStorage
	}

	ip, mask, err := a.parseIPAndMask(dto.IP)
	if err != nil {
		return fmt.Errorf("addToWList: %w", err)
	}

	ipAddress := storage.NewIPAddress(uuid.NewString(), ip, mask)

	err = a.storage.AddToWList(ctx, ipAddress)
	if err != nil {
		return fmt.Errorf("addToWList: %w", err)
	}

	return nil
}

func (a *IPAccess) AddToBList(ctx context.Context, dto IpDTO) error {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return ctxErr
	}

	if a.storage == nil {
		return ErrInvalidStorage
	}

	ip, mask, err := a.parseIPAndMask(dto.IP)
	if err != nil {
		return fmt.Errorf("addToBList: %w", err)
	}

	ipAddress := storage.NewIPAddress(uuid.NewString(), ip, mask)

	err = a.storage.AddToBList(ctx, ipAddress)
	if err != nil {
		return fmt.Errorf("addToBList: %w", err)
	}

	return nil
}

func (a *IPAccess) DeleteFromWList(ctx context.Context, dto IpDTO) error {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return ctxErr
	}

	if a.storage == nil {
		return ErrInvalidStorage
	}

	ip, mask, err := a.parseIPAndMask(dto.IP)
	if err != nil {
		return fmt.Errorf("deleteFromWList: %w", err)
	}

	ipAddress := storage.NewIPAddress("", ip, mask)

	err = a.storage.DeleteFromWList(ctx, ipAddress)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return fmt.Errorf("deleteFromWList: %w", ErrNotFound)
		}
		return fmt.Errorf("deleteFromWList: %w", err)
	}

	return nil
}

func (a *IPAccess) DeleteFromBList(ctx context.Context, dto IpDTO) error {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return ctxErr
	}

	if a.storage == nil {
		return ErrInvalidStorage
	}

	ip, mask, err := a.parseIPAndMask(dto.IP)
	if err != nil {
		return fmt.Errorf("deleteFromBList: %w", err)
	}

	ipAddress := storage.NewIPAddress("", ip, mask)

	err = a.storage.DeleteFromBList(ctx, ipAddress)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return fmt.Errorf("deleteFromBList: %w", ErrNotFound)
		}
		return fmt.Errorf("deleteFromBList: %w", err)
	}

	return nil
}

func (a *IPAccess) IsInWList(ctx context.Context, dto IpDTO) (bool, error) {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return false, ctxErr
	}

	if a.storage == nil {
		return false, ErrInvalidStorage
	}

	ip, mask, err := a.parseIPAndMask(dto.IP)
	if err != nil {
		return false, fmt.Errorf("isInWList: %w", err)
	}

	ipAddress := storage.NewIPAddress("", ip, mask)

	exists, err := a.storage.IsInWList(ctx, ipAddress)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return false, fmt.Errorf("isInWList: %w", ErrNotFound)
		}
		return false, fmt.Errorf("isInWList: %w", err)
	}

	return exists, nil
}

func (a *IPAccess) IsInBList(ctx context.Context, dto IpDTO) (bool, error) {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return false, ctxErr
	}

	if a.storage == nil {
		return false, ErrInvalidStorage
	}

	ip, mask, err := a.parseIPAndMask(dto.IP)
	if err != nil {
		return false, fmt.Errorf("isInBList: %w", err)
	}

	ipAddress := storage.NewIPAddress("", ip, mask)

	exists, err := a.storage.IsInBList(ctx, ipAddress)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return false, fmt.Errorf("isInBList: %w", ErrNotFound)
		}
		return false, fmt.Errorf("isInBList: %w", err)
	}

	return exists, nil
}
