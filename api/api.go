package api

import (
	"fmt"
	// "github.com/gorilla/websocket"
	"github.com/tendermint/go-wire"
	tndr "github.com/tendermint/tendermint/types"
	"github.com/zballs/go_resonate/bigchain"
	"github.com/zballs/go_resonate/state"
	"github.com/zballs/go_resonate/types"
	. "github.com/zballs/go_resonate/util"
	"net/http"
)

const CHAIN = "res()nate"

type Api struct {
	blocks       chan *tndr.Block
	latestHeight int
	logger       types.Logger
	privAcc      *types.PrivateAccount
	proxy        *types.Proxy
	userId       string
}

func NewApi(remote string) *Api {
	return &Api{
		blocks: make(chan *tndr.Block),
		logger: types.NewLogger("api"),
		proxy:  types.NewProxy(remote, "/websocket"),
	}
}

func (api *Api) AddRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/create_user", api.CreateUser)
	mux.HandleFunc("/remove_user", api.RemoveUser)
	mux.HandleFunc("/login", api.Login)
}

func Respond(w http.ResponseWriter, response interface{}) {
	json := ToJSON(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (api *Api) CreateUser(w http.ResponseWriter, req *http.Request) {
	// Should be POST request
	if req.Method != http.MethodPost {
		http.Error(w, fmt.Sprintf("Expected POST request; got %s request", req.Method), http.StatusBadRequest)
		return
	}
	// Get request data
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, "Failed to read request data", http.StatusBadRequest)
		return
	}
	email := values.Get("email")
	password := values.Get("password")
	_type := values.Get("type")
	username := values.Get("username")
	// New user
	user := types.NewUser(email, username, _type)
	// Generate keypair from password
	priv, pub := GenerateKeypair(password)
	data := make(map[string]interface{})
	// Sign user bytes
	sig := priv.Sign(ToJSON(user))
	data["user_signature"] = sig
	// send request to IPDB
	t := bigchain.NewUserTransaction(data, pub)
	t.Fulfill(priv, pub) // should we fulfill here?
	data["user_id"], err = bigchain.PostUserTransaction(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Sign user signature to get account signature
	data["user_signature"] = nil
	data["account_signature"] = priv.Sign(sig.Bytes())
	// Create action
	action := types.NewAction(data, types.CREATE_ACCOUNT)
	// Prepare and sign action
	action.Prepare(pub, 1) // pass sequence=1
	action.Sign(priv, CHAIN)
	// Broadcast tx
	result, err := api.proxy.BroadcastTx("sync", action.Tx())
	if err != nil {
		Respond(w, MessageCreateUser(nil, err))
		return
	}
	if err = ResultToError(result); err != nil {
		Respond(w, MessageCreateUser(nil, err))
		return
	}
	id := data["user_id"].(string)
	userAccount := NewUserAccount(id, priv, pub)
	Respond(w, MessageCreateUser(userAccount, nil))
}

func (api *Api) Login(w http.ResponseWriter, req *http.Request) {
	// Should be POST request
	if req.Method != http.MethodPost {
		http.Error(w, fmt.Sprintf("Expected POST request; got %s request", req.Method), http.StatusBadRequest)
		return
	}
	// Get request data
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, "Failed to read request data", http.StatusBadRequest)
		return
	}
	userId := values.Get("user_id")
	status, err := bigchain.GetTransactionStatus(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	} else if status != "valid" {
		http.Error(w, Error("Invalid user_id").Error(), http.StatusUnauthorized)
		return
	}
	// PrivKey
	priv := new(PrivateKey)
	keystr := values.Get("private_key")
	if err = priv.FromB58(keystr); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// Query Tendermint account
	// Does requiring user to enter public_key through interface provide a security benefit?
	pub := new(PublicKey)
	keystr = values.Get("public_key")
	if err = pub.FromB58(keystr); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}
	accKey := state.AccountKey(pub.Address())
	query := KeyQuery(accKey)
	result, err := api.proxy.TMSPQuery(query)
	if err == nil {
		if err = ResultToError(result); err == nil {
			if err = wire.ReadBinaryBytes(result.Result.Data, &api.privAcc.Account); err == nil {
				api.privAcc.PrivKey = priv
				api.userId = userId
			}
		}
	}
	Respond(w, MessageLogin(err))
	// Start ws
}

func (api *Api) RemoveUser(w http.ResponseWriter, req *http.Request) {
	// Should be POST request
	if req.Method != http.MethodPost {
		http.Error(w, fmt.Sprintf("Expected POST request; got %s request", req.Method), http.StatusBadRequest)
		return
	}
	// Make sure we're logged in
	if api.privAcc == nil || api.userId == "" {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}
	// Get request data
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, "Failed to read request data", http.StatusBadRequest)
		return
	}
	email := values.Get("email")
	_type := values.Get("type")
	username := values.Get("username")
	user := types.NewUser(email, username, _type)
	data := make(map[string]interface{})
	priv := api.privAcc.PrivKey
	sig := priv.Sign(ToJSON(user))
	data["account_signature"] = priv.Sign(sig.Bytes())
	data["user_id"] = api.userId
	// Create action
	action := types.NewAction(data, types.REMOVE_ACCOUNT)
	// Prepare and sign action
	pub := api.privAcc.PubKey
	seq := api.privAcc.Sequence
	action.Prepare(pub, seq)
	action.Sign(priv, CHAIN)
	// Broadcast tx
	result, err := api.proxy.BroadcastTx("sync", action.Tx())
	if err == nil {
		err = ResultToError(result)
	}
	Respond(w, MessageRemoveUser(err))
}
