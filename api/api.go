package api

import (
	"github.com/zballs/go_resonate/bigchain"
	"github.com/zballs/go_resonate/coala"
	"github.com/zballs/go_resonate/types"
	. "github.com/zballs/go_resonate/util"
	"net"
	"net/http"
	"os"
	"strings"
)

type Api struct {
	logger types.Logger
	priv   *PrivateKey
	pub    *PublicKey
	//serv   *types.HttpService
	serv   *types.SocketService
	user   *types.User
	userId string
}

func NewApi(dir string) *Api {
	logger := types.NewLogger("api")
	serv, err := types.NewSocketService("", dir)
	Check(err)
	return &Api{
		logger: logger,
		serv:   serv,
	}
}

func (api *Api) AddRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/create_project", api.CreateProject)
	mux.HandleFunc("/play_song", api.PlaySong)
	mux.HandleFunc("/user_login", api.UserLogin)
	mux.HandleFunc("/user_register", api.UserRegister)
}

func (api *Api) NewProject(projectTitle string) map[string]interface{} {
	artistName := api.user.Name
	data := coala.NewWork(projectTitle, artistName)
	return data
}

func (api *Api) NewSong(songTitle, projectId, projectPlace string) map[string]interface{} {
	example := "" //what should example be?
	isManifestation := true
	date := DateString()
	playAddr := api.serv.PlayAddr()
	data := coala.NewDigitalManifestation(songTitle, example, isManifestation, projectId, date, projectPlace, playAddr)
	return data
}

func (api *Api) UserRegister(w http.ResponseWriter, req *http.Request) {
	// Should be POST request
	if req.Method != http.MethodPost {
		http.Error(w, Sprintf("Expected POST request; got %s request", req.Method), http.StatusBadRequest)
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
	Println(string(MarshalJSON(t)))
	id, err := bigchain.PostTransaction(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userInfo := NewUserInfo(id, priv, pub)
	WriteJSON(w, userInfo)
}

func (api *Api) UserLogin(w http.ResponseWriter, req *http.Request) {
	// Should be POST request
	if req.Method != http.MethodPost {
		http.Error(w, Sprintf("Expected POST request; got %s request", req.Method), http.StatusBadRequest)
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
	WriteJSON(w, "Logged in!")
}

func (api *Api) PlaySong(w http.ResponseWriter, req *http.Request) {
	// Should be GET request
	if req.Method != http.MethodGet {
		http.Error(w, Sprintf("Expected GET request; got %s request", req.Method), http.StatusBadRequest)
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
	values, err := UrlValues(req)
	if err != nil {
		//..
	}
	songId := values.Get("song_id")
	t, err := bigchain.GetTransaction(songId)
	if err != nil {
		//..
	}
	playAddr := t.GetValue("url").(string)
	conn, err := net.Dial("tcp", playAddr)
	if err != nil {
		//..
	}
	defer conn.Close()
	projectTitle := values.Get("project_title")
	songTitle := values.Get("song_title")
	sig := api.priv.Sign([]byte(projectTitle + songTitle))
	playRequest := types.NewPlayRequest(projectTitle, songTitle, api.pub, sig)
	if err := WriteJSON(conn, playRequest); err != nil {
		//..
	}
	Copy(w, conn)
}

func (api *Api) CreateProject(w http.ResponseWriter, req *http.Request) {
	// Should be POST request
	if req.Method != http.MethodPost {
		http.Error(w, Sprintf("Expected POST request; got %s request", req.Method), http.StatusBadRequest)
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
	// Project title
	projectTitle := form.Value["project_title"][0]
	// Create new project
	data := api.NewProject(projectTitle)
	// Generate and send transaction to IPDB
	t := bigchain.GenerateTransaction(data, api.pub)
	t.Fulfill(api.priv, api.pub)
	projectId, err := bigchain.PostTransaction(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	projectPlace := form.Value["project_place"][0]
	songs := form.File["songs"]
	songIds := make([]string, len(songs))
	for i, song := range songs {
		// Store file
		path := api.serv.Path(projectTitle, song.Filename)
		file, err := os.Create(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_file, err := song.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		Copy(file, _file)
		songTitle := strings.SplitN(song.Filename, ".", 2)[0]
		data = api.NewSong(songTitle, projectId, projectPlace)
		// Generate and send transaction to IPDB
		t = bigchain.GenerateTransaction(data, api.pub)
		songIds[i], err = bigchain.PostTransaction(t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	projectInfo := NewProjectInfo(projectId, songIds)
	WriteJSON(w, projectInfo)
}
