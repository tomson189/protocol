/*
	Copyright 2017-2018 OneLedger

	Identities management for any of the associated chains

	TODO: Need to pick a system key for identities. Is a hash of pubkey reasonable?
*/
package id

import (
	"reflect"

	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/err"
	"github.com/Oneledger/protocol/node/log"
	crypto "github.com/tendermint/go-crypto"
	wdata "github.com/tendermint/go-wire/data"
	"golang.org/x/crypto/ripemd160"
)

// Aliases to hide some of the basic underlying types.

type Address = wdata.Bytes // OneLedger address, like Tendermint the hash of the associated PubKey

type PublicKey = crypto.PubKey
type PrivateKey = crypto.PrivKey
type Signature = crypto.Signature

// enum for type
type AccountType int

const (
	UNKNOWN AccountType = iota
	ONELEDGER
	BITCOIN
	ETHEREUM
)

// The persistent collection of all accounts known by this node
type Accounts struct {
	data *data.Datastore
}

func NewAccounts(name string) *Accounts {
	data := data.NewDatastore(name, data.PERSISTENT)

	return &Accounts{
		data: data,
	}
}

func (acc *Accounts) Add(account Account) {
	buffer, _ := comm.Serialize(account)
	acc.data.Store(account.Key(), buffer)
	acc.data.Commit()
}

func (acc *Accounts) Delete(account Account) {
}

func (acc *Accounts) Exists(newType AccountType, name string) bool {
	account := NewAccount(newType, name, PublicKey{})
	value := acc.data.Load(account.Key())
	if value != nil {
		return true
	}
	return false
}

func (acc *Accounts) FindAll() []Account {
	keys := acc.data.List()
	size := len(keys)
	results := make([]Account, size, size)
	for i := 0; i < size; i++ {
		// TODO: This is dangerous...
		account := &AccountOneLedger{}
		base, _ := comm.Deserialize(acc.data.Load(keys[i]), account)
		results[i] = base.(Account)
	}
	return results
}

func (acc *Accounts) Dump() {
	list := acc.FindAll()
	size := len(list)
	for i := 0; i < size; i++ {
		account := list[i]
		log.Info("Account", "Name", account.Name())
		log.Info("Type", "Type", reflect.TypeOf(account))
	}
}

func (acc *Accounts) Find(name string) (Account, err.Code) {
	// TODO: Lookup the identity in the node's database
	return &AccountOneLedger{AccountBase: AccountBase{Name: name}}, 0
}

type AccountKey []byte

// Polymorphism
type Account interface {
	AddPrivateKey(PrivateKey)
	Name() string
	Key() data.DatabaseKey
}

type AccountBase struct {
	Type AccountType

	Key        AccountKey
	Name       string
	PublicKey  PublicKey
	PrivateKey PrivateKey
}

// Hash the public key to get a unqiue hash that can act as a key
func NewAccountKey(key PublicKey) AccountKey {
	hasher := ripemd160.New()

	// TODO: This deosn't seem right?
	bytes, err := key.MarshalJSON()
	if err != nil {
		panic("Unable to Marshal the key into bytes")
	}

	hasher.Write(bytes)

	return hasher.Sum(nil)
}

func NewAccount(newType AccountType, name string, key PublicKey) Account {
	switch newType {

	case ONELEDGER:
		return &AccountOneLedger{
			AccountBase{
				Type:      newType,
				Key:       NewAccountKey(key),
				Name:      name,
				PublicKey: key,
			},
		}

	case BITCOIN:
		return &AccountBitcoin{
			AccountBase{
				Type:      newType,
				Key:       NewAccountKey(key),
				Name:      name,
				PublicKey: key,
			},
		}

	case ETHEREUM:
		return &AccountEthereum{
			AccountBase{
				Type:      newType,
				Key:       NewAccountKey(key),
				Name:      name,
				PublicKey: key,
			},
		}

	default:
		panic("Unknown Type")
	}
}

// Map type to string
func ParseAccountType(typeName string) AccountType {
	switch typeName {
	case "OneLedger":
		return ONELEDGER

	case "Ethereum":
		return ETHEREUM

	case "Bitcoin":
		return BITCOIN
	}
	return UNKNOWN
}

// OneLedger

// Information we need about our own fullnode identities
type AccountOneLedger struct {
	AccountBase
}

func (account *AccountOneLedger) AddPublicKey(key PublicKey) {
	account.PublicKey = key
}

func (account *AccountOneLedger) AddPrivateKey(key PrivateKey) {
	account.PrivateKey = key
}

func (account *AccountOneLedger) Name() string {
	return account.AccountBase.Name
}

func (account *AccountOneLedger) Key() data.DatabaseKey {
	return data.DatabaseKey(account.AccountBase.Name)
}

// Bitcoin

// Information we need for a Bitcoin account
type AccountBitcoin struct {
	AccountBase
}

func (account *AccountBitcoin) AddPublicKey(key PublicKey) {
	account.PublicKey = key
}

func (account *AccountBitcoin) AddPrivateKey(key PrivateKey) {
	account.PrivateKey = key
}

func (account *AccountBitcoin) Name() string {
	return account.AccountBase.Name
}

func (account *AccountBitcoin) Key() data.DatabaseKey {
	return data.DatabaseKey(account.AccountBase.Name)
}

// Ethereum

// Information we need for an Ethereum account
type AccountEthereum struct {
	AccountBase
}

func (account *AccountEthereum) AddPublicKey(key PublicKey) {
	account.PublicKey = key
}

func (account *AccountEthereum) AddPrivateKey(key PrivateKey) {
	account.PrivateKey = key
}

func (account *AccountEthereum) Name() string {
	return account.AccountBase.Name
}

func (account *AccountEthereum) Key() data.DatabaseKey {
	return data.DatabaseKey(account.AccountBase.Name)
}