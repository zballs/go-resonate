package bigchain

import (
	"bytes"
	"fmt"
	. "github.com/zballs/go_resonate/util"
	"io/ioutil"
	"net/http"
)

const IPDB_ENDPOINT = ""

// GET
func GetTransactionStatus(tx_id string) (string, error) {
	url := IPDB_ENDPOINT + "/transactions/" + tx_id + "/status"
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	mp := make(map[string]interface{})
	UnmarshalJSON(data, mp)
	status := mp["status"].(string)
	return status, nil
}

func GetTransaction(tx_id string) (*Transaction, error) {
	url := IPDB_ENDPOINT + "/transactions/" + tx_id
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	t := new(Transaction)
	if err = ReadJSON(response.Body, t); err != nil {
		return nil, err
	}
	return t, nil
}

// POST

// BigchainDB transaction type
// docs.bigchaindb.com/projects/py-driver/en/latest/handcraft.html

const (
	// For ed25519
	BITMASK          = 32
	FULFILLMENT_SIZE = PUBKEY_SIZE + SIGNATURE_SIZE
	TYPE             = "fulfillment"
	TYPE_ID          = 4
	VERSION          = 1

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

func PostTransaction(t *Transaction) (string, error) {
	url := IPDB_ENDPOINT + "/transactions/"
	buf := new(bytes.Buffer)
	if err := ReadJSON(buf, t); err != nil {
		return "", err
	}
	response, err := http.Post(url, "application/json", buf)
	if err != nil {
		return "", err
	}
	rd, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	mp := make(map[string]interface{})
	UnmarshalJSON(rd, mp)
	id := mp["id"].(string)
	return id, nil
}

func GenerateTransaction(data map[string]interface{}, pub *PublicKey) *Transaction {
	asset := NewAsset(data, false, false, false) //should it be updatable?
	details := NewDetails(pub)
	condition := NewCondition(1, 0, details, pub) //what should cid be?
	fulfillment := NewFulfillment(0, pub)         //what should fid be?
	meta := NewMetadata(nil)
	tx := NewTx(
		asset,
		Conditions{condition},
		Fulfillments{fulfillment},
		meta, // what should metadata be?
		CREATE,
	)
	t := NewTransaction(tx, VERSION)
	t.SetValue("@id", t.Id)
	return t
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
	json := MarshalJSON(transaction)
	sum := Checksum256(json)
	transaction.Id = BytesToHex(sum)
	for i, c := range conditions {
		c.Cond.Details.Signature = sigs[i]
	}
	return transaction
}

func (t *Transaction) GetData() map[string]interface{} {
	// for convenience
	return t.Tx.Asset.Data
}

func (t *Transaction) SetData(data map[string]interface{}) {
	t.Tx.Asset.Data = data
}

func (t *Transaction) GetValue(key string) interface{} {
	return t.Tx.Asset.Data[key]
}

func (t *Transaction) SetValue(key string, value interface{}) {
	t.Tx.Asset.Data[key] = value
}

func (t *Transaction) Fulfill(priv *PrivateKey, pub *PublicKey) {
	json := MarshalJSON(t)
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

func NewTx(asset *Asset, conditions Conditions, fulfillments Fulfillments, meta *Metadata, op string) *Tx {
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

func NewCondition(amount, cid int, details *Details, ownersAfter ...*PublicKey) *Condition {
	sig := details.Signature
	details.Signature = ""
	json := MarshalJSON(details)
	sum := Checksum256(json)
	b64 := Base64RawURL(sum)
	uri := fmt.Sprintf("cc:%x:%x:%s:%d", TYPE_ID, BITMASK, b64, FULFILLMENT_SIZE)
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

func NewFulfillment(fid int, ownersBefore ...*PublicKey) *Fulfillment {
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
