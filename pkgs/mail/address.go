package mail

import (
	"encoding/json"
	gomail "net/mail"
)

type Address struct {
	gomail.Address
}

func ParseAddressList(list string) ([]*Address, error) {
	a, err := gomail.ParseAddressList(list)
	if err != nil {
		return nil, err
	}
	aa := []*Address{}
	for _, addr := range a {
		aa = append(aa, &Address{*addr})
	}
	return aa, nil
}

func ParseAddress(address string) (*Address, error) {
	a, err := gomail.ParseAddress(address)
	if err != nil {
		return nil, err
	}
	return &Address{*a}, nil
}

func (a Address) MarshalJSON() ([]byte, error) {
	v := map[string]interface{}{}
	v["name"] = a.Address.Name
	v["address"] = a.Address.Address
	return json.Marshal(v)
}

func (a *Address) UnmarshalJSON(data []byte) error {
	var v map[string]string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	name := v["name"]
	address := v["address"]
	*a = Address{
		gomail.Address{
			Name:    name,
			Address: address,
		},
	}
	return nil
}
