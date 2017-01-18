package state

import (
	"fmt"
	tmsp "github.com/tendermint/tmsp/types"
	"github.com/zballs/go_resonate/bigchain"
	"github.com/zballs/go_resonate/types"
	. "github.com/zballs/go_resonate/util"
)

// If the action is invalid, TMSP error will be returned.
func ExecuteAction(state *State, action types.Action, isCheckTx bool) tmsp.Result {
	chain := state.GetChain()
	// Validate Input Basic
	result := action.Input.ValidateBasic()
	if result.IsErr() {
		return result
	}
	var acc *types.Account
	if action.Type == types.CREATE_ACCOUNT {
		// Create new account..
		// Must have input pubKey
		pub := action.Input.PubKey
		if pub == nil {
			return tmsp.ErrUnauthorized.SetLog("no pubkey in CREATE_ACCOUNT action")
		}
		if action.Data == nil {
			return tmsp.ErrUnauthorized.SetLog("no data in CREATE_ACCOUNT action")
		}
		id := action.Data["user_id"].(string)
		accSig, ok := action.Data["account_signature"].(*Signature)
		if !ok {
			//..
		}
		t, err := bigchain.GetTransaction(id)
		if err != nil {
			//..
		}
		data := t.GetData()
		userSig, ok := data["user_signature"].(*Signature)
		if !ok {
			//..
		}
		// account_sig should be valid signature of the user_sig
		if !pub.Verify(userSig.Bytes(), accSig) {
			return tmsp.ErrUnauthorized.SetLog("user/account verification failed")
		}
		acc = types.NewAccount(pub)
	} else {
		// Get input account
		addr := action.Input.Address
		acc = state.GetAccount(addr)
		if acc == nil {
			return tmsp.ErrBaseUnknownAddress
		}
		if pub := action.Input.PubKey; pub != nil {
			acc.PubKey = pub
		}
	}
	// Validate input, advanced
	signBytes := action.SignBytes(chain)
	result = validateInputAdvanced(acc, signBytes, action.Input)
	if result.IsErr() {
		return result.PrependLog("in validateInputAdvanced()")
	}
	if isCheckTx {
		// CheckTx does not set state
		// Ok, we are done
		return tmsp.OK
	}
	// Increment sequence and create checkpoint
	acc.Sequence += 1
	accCopy := acc.Copy()
	// Run the action.
	cache := state.CacheWrap()
	switch action.Type {
	case types.CREATE_ACCOUNT:
		// Set address to new account
		addr := acc.PubKey.Address()
		cache.SetAccount(addr, acc)
		result = tmsp.OK
	case types.REMOVE_ACCOUNT:
		if action.Data == nil {
			result = tmsp.ErrUnauthorized.SetLog("no data in REMOVE_ACCOUNT action")
		} else {
			id := action.Data["user_id"].(string)
			accSig, ok := action.Data["account_signature"].(*Signature)
			if !ok {
				//..
			}
			t, err := bigchain.GetTransaction(id)
			if err != nil {
				//..
			}
			data := t.GetData()
			userSig, ok := data["user_signature"].(*Signature)
			if !ok {
				//..
			}
			// account_sig should be valid signature of the user_sig
			if !acc.PubKey.Verify(userSig.Bytes(), accSig) {
				result = tmsp.ErrUnauthorized.SetLog("user/account verification failed")
			} else {
				// Set address to nil
				addr := acc.PubKey.Address()
				cache.SetAccount(addr, nil)
				result = tmsp.OK
			}
		}
	default:
		result = tmsp.ErrUnknownRequest
		// TODO: add more actions
	}
	if result.IsOK() {
		logger.Info("Success")
		cache.CacheSync()
	} else {
		logger.Info("AppTx failed", "error", result)
		cache.SetAccount(action.Input.Address, accCopy)
	}
	return result
}

func validateInputAdvanced(acc *types.Account, signBytes []byte, in *types.ActionInput) (res tmsp.Result) {
	if in == nil {
		// shouldn't happen
	}
	// Check sequence
	seq := acc.Sequence
	if seq+1 != in.Sequence {
		return tmsp.ErrBaseInvalidSequence.AppendLog(
			fmt.Sprintf("Got %v, expected %v. (acc.seq=%v)", in.Sequence, seq+1, acc.Sequence))
	}
	// Check signature
	if !acc.PubKey.Verify(signBytes, in.Signature) {
		return tmsp.ErrBaseInvalidSignature.AppendLog(
			fmt.Sprintf("SignBytes: %X", signBytes))
	}
	return tmsp.OK
}
