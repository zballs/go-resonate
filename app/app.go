package app

import (
	"bytes"
	"github.com/tendermint/go-wire"
	"github.com/tendermint/merkleeyes"
	merk "github.com/tendermint/merkleeyes/client"
	tmsp "github.com/tendermint/tmsp/types"
	sm "github.com/zballs/go_resonate/state"
	. "github.com/zballs/go_resonate/types"
	. "github.com/zballs/go_resonate/util"
	"strings"
	"time"
)

const (
	CHAIN   = "res()nate"
	VERSION = "1.0.0"
)

type App struct {
	cli   *merk.Client
	cache *state.State
	state *state.State
}

func NewApp(cli *merk.Client) *App {
	state := state.NewState(cli)
	return &App{
		cli:   cli,
		state: state,
	}
}

// TMSP Requests

func (app *App) Info() string {
	return fmt.Sprintf("res()nate v%s", VERSION)
}

func (app *App) SetOption(key, value string) string {
	if strings.Contains(key, "/") {
		parts := strings.SplitN(key, "/", 2)
		key = parts[1]
	}
	switch key {
	case "chain":
		app.state.SetChain(value)
		return "Success"
	case "account":
		var acc *Account
		var err error
		wire.ReadJSONPtr(&acc, []byte(value), err)
		if err != nil {
			return "Error decoding account: " + err.Error()
		}
		app.state.SetAccount(acc)
		return "Success"
	}
	return "Unrecognized option key " + key
}

func (app *App) DeliverTx(tx []byte) tmsp.Result {
	var action Action
	err := wire.ReadBinaryBytes(tx, &action)
	if err != nil {
		return tmsp.ErrBaseEncodingError.AppendLog("Error decoding tx: " + err.Error())
	}
}
