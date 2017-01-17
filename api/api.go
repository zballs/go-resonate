package api

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/tendermint/go-wire" //TODO: move to util
	tndr "github.com/tendermint/tendermint/types"
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
	user         *types.User
}

func NewApi(remote string) *Api {
	return &Api{
		blocks: make(chan *tndr.Block),
		logger: types.NewLogger("api"),
		proxy:  types.NewProxy(remote, "/websocket"),
	}
}

func (api *Api) AddRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/create_account", api.CreateAccount)
	mux.HandleFunc("/remove_account", api.RemoveAccount)
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
	}
	email := values.Get("email")
	password := values.Get("password")
	username := values.Get("username")
	// New user
	user := types.NewUser(email, username)
	// Generate keypair from password
	priv, pub := GenerateKeypair(password)
	// Sign user bytes
	sig := priv.Sign(ToJSON(user))
	data := sig.Bytes()
	// Create action
	action := types.NewAction(data, types.CREATE_USER)
	// Prepare and sign action
	action.Prepare(pub, 1) // pass sequence=1
	action.Sign(priv, CHAIN)
	// Broadcast tx
	result, err := api.proxy.BroadcastTx("sync", action.Tx())
	if err != nil {
		Respond(w, MessageCreateAccount(nil, err))
		return
	}
	if err = ResultToError(result); err != nil {
		Respond(w, MessageCreateAccount(nil, err))
		return
	}
	keypair := NewKeypair(pub, priv)
	Respond(w, MessageCreateUser(keypair, nil))
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
	// User info
	email := values.Get("email")
	username := values.Get("username")
	user := types.NewUser(email, username)
	// PrivKey
	priv := new(PrivateKey)
	keystr := values.Get("private_key")
	if err = priv.FromB58(keystr); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// Sign user bytes
	sig := priv.Sign(ToJSON(user))
	data := sig.Bytes()
	// * * * * * * * * * * * * * * * * * * * * * * * * * * *
	// TODO: query bigchaindb node; verify that user exists
	// * * * * * * * * * * * * * * * * * * * * * * * * * * *
	// Query account
	// Does requiring user to enter public_key through interface provide a security benefit?
	pub := new(PublicKey)
	keystr = values.Get("public_key")
	if err = pub.FromB58(keystr); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}
	accKey := state.AccountKey(pub.Address())
	query := KeyQuery(acckey)
	result, err := api.proxy.TMSPQuery(query)
	if err == nil {
		err = ResultToError(result)
		if err == nil {
			err = wire.ReadBinaryBytes(result.Result.Data, &api.user.Account)
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
	var privAcc *types.PrivateAccount
	if privAcc = api.user; privAcc == nil {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}
	// Should be a POST request
	if req.Method != http.MethodPost {
		http.Error(w, fmt.Sprintf("Expected POST request; got %s request", req.Method), http.StatusBadRequest)
		return
	}
	// Create action
	action := types.NewAction(nil, types.REMOVE_USER)
	// Prepare and sign action
	action.Prepare(privAcc.PubKey, privAcc.Sequence)
	action.Sign(privAcc.PrivKey, CHAIN)
	// Broadcast tx
	result, err := api.proxy.BroadcastTx("sync", action.Tx())
	if err == nil {
		err = ResultToError(result)
	}
	Respond(w, MessageRemoveAccount(err))
}
