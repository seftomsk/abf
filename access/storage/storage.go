package storage

import "errors"

var ErrInvalidInitialization = errors.New("invalid initialization struct")

var ErrNotFound = errors.New("not found")

var ErrInvalidEntity = errors.New("invalid entity")

func ValidEntity(ip IPEntity) bool {
	if ip.Mask() == "" || ip.IP() == "" {
		return false
	}

	return true
}

func ValidWholeEntity(ip IPEntity) bool {
	if !ValidEntity(ip) || ip.ID() == "" {
		return false
	}

	return true
}

type IPEntity interface {
	ID() string
	IP() string
	Mask() string
}

type IPAddress struct {
	id   string
	ip   string
	mask string
}

func (ip *IPAddress) ID() string {
	return ip.id
}

func (ip *IPAddress) IP() string {
	return ip.ip
}

func (ip *IPAddress) Mask() string {
	return ip.mask
}

func NewIPAddress(id, ip, mask string) *IPAddress {
	return &IPAddress{
		id:   id,
		ip:   ip,
		mask: mask,
	}
}
