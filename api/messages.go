package api

import (
	"github.com/pkg/errors"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	. "github.com/zballs/go_resonate/util"
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

type Keypair struct {
	Privkey *PrivateKey `json:"private_key"`
	Pubkey  *PublicKey  `json:"public_key"`
}

func NewKeypairB58(pub *PublicKey, priv *PrivateKey) *Keypair {
	return &Keypair{priv, pub}
}

func MessageCreateUser(data *Keypair, err error) *Message {
	return &Message{
		Action: "create_account",
		Data:   data,
		Error:  err,
	}
}

func MessageRemoveUser(err error) *Message {
	return &Message{
		Action: "remove_account",
		Error:  err,
	}
}

func MessageLogin(err error) *Message {
	return &Message{
		Action: "login",
		Error:  err,
	}
}
