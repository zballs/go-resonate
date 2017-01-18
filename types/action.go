package types

import (
	"fmt"
	"github.com/tendermint/go-wire"
	tndr "github.com/tendermint/tendermint/types"
	tmsp "github.com/tendermint/tmsp/types"
	. "github.com/zballs/go_resonate/util"
)

const (
	CREATE_ACCOUNT = 0x01
	REMOVE_ACCOUNT = 0x02

	// TODO: add more actions
)

type ActionInput struct {
	Address   []byte     `json: "address"`
	PubKey    *PublicKey `json: "pub_key"`
	Sequence  int        `json: "sequence"`
	Signature *Signature `json: "signature"`
}

func (in ActionInput) ValidateBasic() tmsp.Result {
	if len(in.Address) != 20 {
		return tmsp.ErrBaseInvalidInput.AppendLog("Invalid address length")
	}
	if in.Sequence <= 0 {
		return tmsp.ErrBaseInvalidInput.AppendLog("Sequence must be greater than 0")
	}
	if in.Sequence == 1 && in.PubKey == nil {
		return tmsp.ErrBaseInvalidInput.AppendLog("PubKey must be present when Sequence == 1")
	}
	if in.Sequence > 1 && in.PubKey != nil {
		return tmsp.ErrBaseInvalidInput.AppendLog("PubKey must be nil when Sequence > 1")
	}
	return tmsp.OK
}

func (in ActionInput) StringIndented(indent string) string {
	return fmt.Sprintf(`Input{
		%s %s Address: %x
		%s %s Sequence: %v
		%s %s PubKey: %v
		%s}`,
		indent, indent, in.Address,
		indent, indent, in.Sequence,
		indent, indent, in.PubKey,
		indent)
}

func (in ActionInput) String() string {
	return in.StringIndented("")
}

type Action struct {
	Data  map[string]interface{} `json: "data"`
	Input *ActionInput           `json: "input"`
	Type  byte                   `json: "type"`
}

func NewAction(data map[string]interface{}, _type byte) Action {
	return Action{
		Data: data,
		Type: _type,
	}
}

func (a Action) SignBytes(chain string) []byte {
	signBytes := wire.BinaryBytes(chain)
	sig := a.Input.Signature
	a.Input.Signature = nil
	signBytes = append(signBytes, wire.BinaryBytes(a)...)
	a.Input.Signature = sig
	return signBytes
}

func (a Action) Prepare(pub *PublicKey, seq int) {
	a.Input.Sequence = seq
	if a.Input.Sequence == 1 {
		a.Input.PubKey = pub
	}
	a.Input.Address = pub.Address()
}

func (a Action) Sign(priv *PrivateKey, chain string) {
	a.Input.Signature = priv.Sign(a.SignBytes(chain))
}

func (a Action) Id(chain string) []byte {
	signBytes := a.SignBytes(chain)
	return wire.BinaryRipemd160(signBytes)
}

func (a Action) Tx() tndr.Tx {
	return wire.BinaryBytes(a)
}

func (a Action) StringIndented(indent string) string {
	return fmt.Sprintf(`Action{
		%s Type: %v
		%s %v 
		%s Data: %x
		}`,
		indent, a.Type,
		indent, a.Input,
		indent, a.Data)
}

func (a Action) String() string {
	return a.StringIndented("")
}
