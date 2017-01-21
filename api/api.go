package api

import (
	"fmt"
	// "github.com/gorilla/websocket"
	"github.com/zballs/go_resonate/bigchain"
	"github.com/zballs/go_resonate/coala"
	"github.com/zballs/go_resonate/types"
	. "github.com/zballs/go_resonate/util"
	"net/http"
)

type Api struct {
	logger types.Logger
	priv   *PrivateKey
	pub    *PublicKey
	user   *types.User
	userId string
}

func NewApi() *Api {
	return &Api{
		logger: types.NewLogger("api"),
	}
}

func (api *Api) AddRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/login", api.Login)
	mux.HandleFunc("/new_project", api.NewProject)
	mux.HandleFunc("/register_user", api.RegisterUser)
}

func Respond(w http.ResponseWriter, response interface{}) {
	json := MarshalJSON(response)
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(json)
	Check(err)
}

func (api *Api) RegisterUser(w http.ResponseWriter, req *http.Request) {
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
	user := types.NewUser(email, password, username, _type)
	// Generate keypair from password
	priv, pub := GenerateKeypair(password)
	data := make(map[string]interface{})
	// Sign user bytes
	json := MarshalJSON(user)
	data["user_signature"] = priv.Sign(json).ToB58()
	// send request to IPDB
	t := bigchain.GenerateTransaction(data, pub)
	t.Fulfill(priv, pub)
	fmt.Println(string(MarshalJSON(t)))
	id, err := bigchain.PostTransaction(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userInfo := NewUserInfo(id, priv, pub)
	Respond(w, userInfo)
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
	// PrivKey
	priv := new(PrivateKey)
	keystr := values.Get("private_key")
	if err = priv.FromB58(keystr); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// User Id
	userId := values.Get("user_id")
	// Check that transaction with id exists
	t, err := bigchain.GetTransaction(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// Get user signature
	data := t.GetData()
	sig := new(Signature)
	if err := sig.FromB58(data["user_signature"].(string)); err != nil {
		http.Error(w, Error("Failed to verify user signature").Error(), http.StatusUnauthorized)
		return
	}
	// User
	email := values.Get("email")
	password := values.Get("password")
	_type := values.Get("type")
	username := values.Get("username")
	// New user
	user := types.NewUser(email, username, password, _type)
	json := MarshalJSON(user)
	pub := priv.Public()
	if !pub.Verify(json, sig) {
		http.Error(w, Error("Failed to verify user signature").Error(), http.StatusUnauthorized)
		return
	}
	api.priv = priv
	api.pub = pub
	api.user = user
	api.userId = userId
	Respond(w, "Logged in!")
}

func (api *Api) NewProject(w http.ResponseWriter, req *http.Request) {
	// Should be POST request
	if req.Method != http.MethodPost {
		http.Error(w, fmt.Sprintf("Expected POST request; got %s request", req.Method), http.StatusBadRequest)
		return
	}
	// Make sure we're logged in
	if api.priv == nil {
		http.Error(w, "Privkey is not set", http.StatusUnauthorized)
		return
	}
	if api.pub == nil {
		http.Error(w, "Pubkey is not set", http.StatusUnauthorized)
		return
	}
	if api.user == nil {
		http.Error(w, "User is not set", http.StatusUnauthorized)
		return
	}
	if api.userId == "" {
		http.Error(w, "User Id is not set", http.StatusUnauthorized)
		return
	}
	// Get request data
	form, err := MultipartForm(req)
	if err != nil {
		http.Error(w, "Failed to read request data", http.StatusBadRequest)
		return
	}
	// Name of the project
	name := form.Value["name"][0]
	// Create new work
	data := coala.NewWork("", name, api.user.Name)
	// Generate and send transaction to IPDB
	t := bigchain.GenerateTransaction(data, api.pub)
	t.Fulfill(api.priv, api.pub)
	projectId, err := bigchain.PostTransaction(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Just one song for now
	// .. Eventually iterate through list of songs
	name = form.Value["names"][0]
	place := form.Value["place"][0]
	// Should we have example, as specified in coalaip??
	// .. What should url be?
	data = coala.NewDigitalManifestation("", name, "", true, projectId, DateString(), place, req.RemoteAddr)
	// Generate and send transaction to IPDB
	t = bigchain.GenerateTransaction(data, api.pub)
	songId, err := bigchain.PostTransaction(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO: save project_id and song_ids to disk
	projectInfo := NewProjectInfo(projectId, []string{songId})
	Respond(w, projectInfo)
}
