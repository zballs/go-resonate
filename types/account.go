package types

import . "github.com/zballs/go_resonate/util"

// Basic Account type

type Account struct {
	PubKey   *PublicKey `json:"public_key"`
	Sequence int        `json:"sequence"`
}

func NewAccount(pub *PublicKey) *Account {
	return &Account{
		PubKey: pub,
	}
}

func (acc *Account) Copy() *Account {
	accCopy := *acc
	return &accCopy
}

func (acc *Account) Address() []byte {
	return acc.PubKey.Address()
}

// Private Account

type PrivateAccount struct {
	*Account
	PrivKey *PrivateKey
}

func NewPrivateAccount(acc *Account, priv *PrivateKey) *PrivateAccount {
	return &PrivateAccount{acc, priv}
}

type AccountGetter interface {
	GetAccount(addr []byte) *Account
}

type AccountSetter interface {
	SetAccount(addr []byte, acc *Account)
}

type AccountGetterSetter interface {
	GetAccount(addr []byte) *Account
	SetAccount(addr []byte, acc *Account)
}
