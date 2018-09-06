package config

import (
	"encoding/json"
	"fmt"

	"github.com/uphy/elastalert-mail-gateway/pkgs/mail"
)

type (
	Address struct {
		*mail.Address
	}
	AddressList struct {
		Addresses []*mail.Address
	}
)

func (a Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

func (a *Address) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	addr, err := mail.ParseAddress(s)
	if err != nil {
		return err
	}
	*a = Address{addr}
	return nil
}

func (a AddressList) MarshalJSON() ([]byte, error) {
	list := make([]string, len(a.Addresses))
	for _, aa := range a.Addresses {
		list = append(list, aa.String())
	}
	return json.Marshal(list)
}

func (a *AddressList) UnmarshalJSON(data []byte) error {
	var list []string
	if err := json.Unmarshal(data, &list); err == nil {
		addrs := []*mail.Address{}
		for _, s := range list {
			addr, err := mail.ParseAddress(s)
			if err != nil {
				return err
			}
			addrs = append(addrs, addr)
		}
		*a = AddressList{
			Addresses: addrs,
		}
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		addr, err := mail.ParseAddress(s)
		if err != nil {
			return err
		}
		*a = AddressList{
			Addresses: []*mail.Address{addr},
		}
		return nil
	}
	return fmt.Errorf("unexpected address list: %s", string(data))
}
