package api

import (
	"github.com/pkg/errors"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	. "github.com/zballs/go_resonate/util"
)

const (
	CREATE_USER = "create_user"
	LOGIN       = "login"
	REMOVE_USER = "remove_user"
)

func ResultToError(result interface{}) error {
	switch result.(type) {
	case *ctypes.ResultTMSPQuery:
		tmResult := result.(*ctypes.ResultTMSPQuery).Result
		if tmResult.Code == 0 {
			return nil
		}
		return errors.New(tmResult.Error())
	case *ctypes.ResultBroadcastTx:
		_result := result.(*ctypes.ResultBroadcastTx)
		if _result.Code == 0 {
			return nil
		}
		return errors.New(_result.Log)
	default:
		return errors.New("Unrecognized result type")
	}
}

type Message struct {
	Action string      `json:"action"`
	Data   interface{} `json:"data, omitempty"`
	Error  error       `json:"error, omitempty"`
}

type UserAccount struct {
	Id      string      `json:"id"`
	Privkey *PrivateKey `json:"private_key"`
	Pubkey  *PublicKey  `json:"public_key"`
}

func NewUserAccount(id string, priv *PrivateKey, pub *PublicKey) *UserAccount {
	return &UserAccount{
		Id:      id,
		Privkey: priv,
		Pubkey:  pub,
	}
}

func MessageCreateUser(data *UserAccount, err error) *Message {
	return &Message{
		Action: CREATE_USER,
		Data:   data,
		Error:  err,
	}
}

func MessageRemoveUser(err error) *Message {
	return &Message{
		Action: REMOVE_USER,
		Error:  err,
	}
}

func MessageLogin(err error) *Message {
	return &Message{
		Action: LOGIN,
		Error:  err,
	}
}
