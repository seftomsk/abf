package access

import (
	"fmt"
	"net"
)

type IStorage interface {
	AddToWhiteList(ip, mask string)
	AddToBlackList(ip, mask string)
	GetAll() map[string]map[string][]string
}

type IPAccess struct {
	storage IStorage
}

func NewIPAccess(storage IStorage) *IPAccess {
	return &IPAccess{storage: storage}
}

func (a *IPAccess) AddToWhiteList(ip string) error {
	if len(ip) <= 0 {
		return fmt.Errorf("addToWhiteList: ip address is epmty")
	}

	ipAddress, ipNet, err := net.ParseCIDR(ip)
	if err != nil {
		return fmt.Errorf("addToWhiteList: %w", err)
	}

	byteMask := ipNet.Mask
	mask := fmt.Sprintf(
		"%d.%d.%d.%d",
		byteMask[0],
		byteMask[1],
		byteMask[2],
		byteMask[3])

	a.storage.AddToWhiteList(ipAddress.String(), mask)

	return nil
}

func (a *IPAccess) GetAll() map[string]map[string][]string {
	return a.storage.GetAll()
}

func (a *IPAccess) DeleteFromWhiteList(ip string) error {
	return nil
}

func (a *IPAccess) AddToBlackList(ip string) error {
	if len(ip) <= 0 {
		return fmt.Errorf("addToBlackList: ip address is epmty")
	}

	ipAddress, ipNet, err := net.ParseCIDR(ip)
	if err != nil {
		return fmt.Errorf("addToBlackList: %w", err)
	}

	byteMask := ipNet.Mask
	mask := fmt.Sprintf(
		"%d.%d.%d.%d",
		byteMask[0],
		byteMask[1],
		byteMask[2],
		byteMask[3])

	a.storage.AddToBlackList(ipAddress.String(), mask)

	return nil
}

func (a *IPAccess) DeleteFromBlackList(ip string) error {
	return nil
}
