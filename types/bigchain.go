package types

import (
	"fmt"
	. "github.com/zballs/go_resonate/util"
	"io/ioutil"
	"net/http"
	// "net/url"
)

// GET
func GetTransactionStatus(endpoint, tx_id string) (string, error) {
	url := fmt.Sprintf("%s/transactions/%s/status", endpoint, tx_id)
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	mp := make(map[string]interface{})
	FromJSON(data, mp)
	status := mp["status"].(string)
	return status, nil
}

func GetTransaction(endpoint, tx_id string) (*Transaction, error) {
	url := fmt.Sprintf("%s/transactions/%s", endpoint, tx_id)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	t := new(Transaction)
	data, err := ReadJSON(response.Body, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// POST

// BigchainDB transaction type
// https://docs.bigchaindb.com/projects/py-driver/en/latest/handcraft.htmls

const (
	// For ed25519
	BITMASK            = 32
	FULFILLMENT_LENGTH = PUBKEY_LENGTH + SIGNATURE_LENGTH
	TYPE               = "fulfillment"
	TYPE_ID            = 4

	// Operation types
	CREATE   = "CREATE"
	GENESIS  = "GENESIS"
	TRANSFER = "TRANSFER"

	// Regex
	FULFILLMENT_REGEX = `^cf:([1-9a-f][0-9a-f]{0,3}|0):[a-zA-Z0-9_-]*$`
)

type Transaction struct {
	Id      string `json:"id"`
	Tx      *Tx    `json:"tx"`
	Version int    `json:"version"`
}

func NewTransaction(tx *Tx, version int) *Transaction {
	transaction := &Transaction{
		Tx:      tx,
		Version: version,
	}
	conditions := transaction.Tx.Conditions
	sigs := make([]string, len(conditions))
	for i, c := range conditions {
		sigs[i] = c.Cond.Details.Signature
		c.Cond.Details.Signature = ""
	}
	json := ToJSON(transaction)
	sum := Checksum256(json)
	transaction.Id = BytesToHex(sum)
	for i, c := range conditions {
		c.Cond.Details.Signature = sigs[i]
	}
	return transaction
}

func (t *Transaction) Fulfill(priv *PrivateKey, pub *PublicKey) {
	json := ToJSON(t)
	sig := priv.Sign(json)
	data := append(pub.Bytes(), sig.Bytes()...)
	b64 := Base64RawURL(data)
	f := fmt.Sprintf("cf:%x:%s", TYPE_ID, b64)
	fulfillments := t.Tx.Fulfillments
	for i, _ := range fulfillments {
		fulfillments[i].Fulfill = f
	}
	t.Tx.Fulfillments = fulfillments //necessary?
}

type Tx struct {
	Asset        *Asset       `json:"asset"`
	Conditions   Conditions   `json:"conditions"`
	Fulfillments Fulfillments `json:"fulfillments"`
	Metadata     *Metadata    `json:"metadata"`
	Operation    string       `json:"operation"`
}

func NewTx(asset *Asset, conditions Conditions, fulfillments Fulfillments, meta *Metadata, op Operation) *Tx {
	return &Tx{
		Asset:        asset,
		Conditions:   conditions,
		Fulfillments: fulfillments,
		Metadata:     meta,
		Operation:    op,
	}
}

func NewCreateTx(asset *Asset, conditions Conditions, fulfillments Fulfillments, meta *Metadata) *Tx {
	return NewTx(asset, conditions, fulfillments, meta, CREATE)
}

type Asset struct {
	Data       map[string]interface{} `json:"data"` //--> coalaip model
	Divisible  bool                   `json:"divisible"`
	Id         string                 `json:"id"`
	Refillable bool                   `json:"refillable"`
	Updatable  bool                   `json:"updatable"`
}

func NewAsset(data map[string]interface{}, divisible, refillable, updatable bool) *Asset {
	id := Uuid4()
	return &Asset{
		Data:       data,
		Divisible:  divisible,
		Id:         id,
		Refillable: refillable,
		Updatable:  updatable,
	}
}

type Condition struct {
	Amount      int          `json:"amount"`
	CID         int          `json:"cid"`
	Cond        *Cond        `json:"condition"`
	OwnersAfter []*PublicKey `json:"owners_after"`
}

type Conditions []*Condition

func NewCondition(amount, cid int, details *Details, ownersAfter []*PublicKey) *Condition {
	sig := details.Signature
	details.Signature = ""
	json := ToJSON(details)
	sum := Checksum256(json)
	b64 := Base64RawURL(sum)
	uri := fmt.Sprintf("cc:%x:%x:%s:%d", TYPE_ID, BITMASK, b64, FULFILLMENT_LENGTH)
	details.Signature = sig
	return &Condition{
		Amount: amount,
		CID:    cid,
		Cond: &Cond{
			Uri:     uri,
			Details: details,
		},
		OwnersAfter: ownersAfter,
	}
}

type Cond struct {
	Uri     string   `json:"uri"`
	Details *Details `json:"details"`
}

type Details struct {
	Bitmask   int        `json:"bitmask"`
	PublicKey *PublicKey `json:"public_key"`
	Signature string     `json:"signature"`
	Type      string     `json:"type"`
	TypeId    int        `json:"type_id"`
}

func NewDetails(pub *PublicKey) *Details {
	return &Details{
		Bitmask:   BITMASK,
		PublicKey: pub,
		Type:      TYPE,
		TypeId:    TYPE_ID,
	}
}

type Fulfillment struct {
	FID          int                    `json:"fid"`
	Fulfill      string                 `json:"fulfillment"`
	Input        map[string]interface{} `json:"input"`
	OwnersBefore []*PublicKey           `json:"owners_before"`
}

type Fulfillments []*Fulfillment

func NewFulfillment(fid int, ownersBefore []*PublicKey) *Fulfillment {
	return &Fulfillment{
		FID:          fid,
		OwnersBefore: ownersBefore,
	}
}

type Metadata struct {
	Data map[string]interface{} `json:"data"`
	Id   string                 `json:"id"`
}

func NewMetadata(data map[string]interface{}) *Metadata {
	id := Uuid4()
	return &Metadata{
		Data: data,
		Id:   id,
	}
}
