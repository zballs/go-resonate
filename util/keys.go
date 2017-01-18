package util

import (
	"github.com/tendermint/go-crypto"
	bcrypt "golang.org/x/crypto/bcrypt"
)

const PUBKEY_LENGTH = 32
const PRIVKEY_LENGTH = 64
const SIGNATURE_LENGTH = 64

type PublicKey struct {
	key crypto.PubKeyEd25519
}

func (pub *PublicKey) Bytes() []byte {
	return pub.key[:]
}

func (pub *PublicKey) Address() []byte {
	return pub.key.Address()
}

func (pub *PublicKey) ToHex() string {
	return BytesToHex(pub.Bytes())
}

func (pub *PublicKey) FromHex(hexstr string) error {
	data := BytesFromHex(hexstr)
	if length := len(data); length != PUBKEY_LENGTH {
		return Errorf("Expected public key with length=%d; got length=%d\n", PUBKEY_LENGTH, length)
	}
	copy(pub.Bytes(), data)
	return nil
}

func (pub *PublicKey) ToB58() string {
	return BytesToB58(pub.Bytes())
}

func (pub *PublicKey) FromB58(b58 string) error {
	data := BytesFromB58(b58)
	if length := len(data); length != PUBKEY_LENGTH {
		return Errorf("Expected public key with length=%d; got length=%d\n", PUBKEY_LENGTH, length)
	}
	copy(pub.Bytes(), data)
	return nil
}

func (pub *PublicKey) MarshalJSON() ([]byte, error) {
	return []byte(pub.ToB58()), nil
}

func (pub *PublicKey) UnmarshalJSON(data []byte) error {
	if err := pub.FromB58(string(data)); err != nil {
		return err
	}
	return nil
}

func (pub *PublicKey) Verify(data []byte, sig *Signature) bool {
	return pub.key.VerifyBytes(data, sig.s)
}

type PrivateKey struct {
	key crypto.PrivKeyEd25519
}

func (priv *PrivateKey) PublicKey() *PublicKey {
	key := priv.key.PubKey().(crypto.PubKeyEd25519)
	return &PublicKey{key}
}

func (priv *PrivateKey) ToHex() string {
	return BytesToHex(priv.key[:])
}

func (priv *PrivateKey) FromHex(hexstr string) error {
	data := BytesFromHex(hexstr)
	if length := len(data); length != PRIVKEY_LENGTH {
		return Errorf("Expected private key with length=%d; got length=%d\n", PRIVKEY_LENGTH, length)
	}
	copy(priv.key[:], data)
	return nil
}

func (priv *PrivateKey) ToB58() string {
	return BytesToB58(priv.key[:])
}

func (priv *PrivateKey) FromB58(b58 string) error {
	data := BytesFromB58(b58)
	if length := len(data); length != PRIVKEY_LENGTH {
		return Errorf("Expected private key with length=%d; got length=%d\n", PRIVKEY_LENGTH, length)
	}
	copy(priv.key[:], data)
	return nil
}

func (priv *PrivateKey) MarshalJSON() ([]byte, error) {
	return []byte(priv.ToB58()), nil
}

func (priv *PrivateKey) UnmarshalJSON(data []byte) error {
	if err := priv.FromB58(string(data)); err != nil {
		return err
	}
	return nil
}

func (priv *PrivateKey) Sign(data []byte) (sig *Signature) {
	s := priv.key.Sign(data)
	sig.s = s.(crypto.SignatureEd25519)
	return sig
}

type Signature struct {
	s crypto.SignatureEd25519
}

func (sig *Signature) Bytes() []byte {
	return sig.s[:]
}

func (sig *Signature) FromBytes(data []byte) error {
	if length := len(data); length != SIGNATURE_LENGTH {
		return Errorf("Expected data with length=%d; got length=%d\n", SIGNATURE_LENGTH, length)
	}
	copy(sig.s[:], data)
	return nil
}

func (sig *Signature) ToB58() string {
	return BytesToB58(sig.Bytes())
}

func (sig *Signature) FromB58(b58 string) error {
	data := BytesFromB58(b58)
	if length := len(data); length != PRIVKEY_LENGTH {
		return Errorf("Expected private key with length=%d; got length=%d\n", PRIVKEY_LENGTH, length)
	}
	copy(sig.Bytes(), data)
	return nil
}

func (sig *Signature) ToHex() string {
	return BytesToHex(sig.Bytes())
}

func (sig *Signature) FromHex(hexstr string) error {
	data := BytesFromHex(hexstr)
	if length := len(data); length != PRIVKEY_LENGTH {
		return Errorf("Expected private key with length=%d; got length=%d\n", PRIVKEY_LENGTH, length)
	}
	copy(sig.Bytes(), data)
	return nil
}

func (sig *Signature) MarshalJSON() ([]byte, error) {
	return []byte(sig.ToB58()), nil
}

func (sig *Signature) UnmarshalJSON(data []byte) error {
	if err := sig.FromB58(string(data)); err != nil {
		return err
	}
	return nil
}

// Generate secret from password string

func GenerateSecret(password string) ([]byte, error) {
	secret, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	Check(err)
	return secret, nil
}

// Generate keypair from password string

func GenerateKeypair(password string) (*PrivateKey, *PublicKey) {
	secret, err := GenerateSecret(password)
	Check(err)
	priv, pub := new(PrivateKey), new(PublicKey)
	priv.key = crypto.GenPrivKeyEd25519FromSecret(secret)
	pub.key = priv.key.PubKey().(crypto.PubKeyEd25519)
	return priv, pub
}
