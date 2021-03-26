package node

import (
	"errors"
	"fmt"

	"github.com/tokentransfer/go-MerklePatriciaTree/mpt"

	"github.com/tokentransfer/chain/block"
	"github.com/tokentransfer/chain/store"

	libblock "github.com/tokentransfer/interfaces/block"
	libcore "github.com/tokentransfer/interfaces/core"
	libcrypto "github.com/tokentransfer/interfaces/crypto"
	libstore "github.com/tokentransfer/interfaces/store"
)

type MerkleTree struct {
	mt *mpt.Trie
	cs libcrypto.CryptoService
	ss libstore.KvService
}

func NewMerkleTree(cs libcrypto.CryptoService, ss libstore.KvService) *MerkleTree {
	return &MerkleTree{
		mt: mpt.New(cs, ss),
		cs: cs,
		ss: ss,
	}
}

func (t *MerkleTree) GetRoot() []byte {
	return t.mt.RootHash()
}

func (t *MerkleTree) Commit() error {
	return t.mt.Commit()
}

func (t *MerkleTree) Cancel() error {
	return t.mt.Abort()
}

func (t *MerkleTree) Verify(key []byte) ([]byte, error) {
	return nil, errors.New("unsupport")
}

func (t *MerkleTree) GetData(key []byte) ([]byte, error) {
	return t.mt.Get(key)
}

func (t *MerkleTree) PutData(key, value []byte) error {
	return t.mt.Put(key, value)
}

func (t *MerkleTree) Init(c libcore.Config) error {
	return nil
}

func (t *MerkleTree) Start() error {
	return nil
}

func (t *MerkleTree) Close() error {
	return nil
}

func (t *MerkleTree) PutDatas(keys [][]byte, values [][]byte) error {
	ll := len(keys)
	if ll == 0 || ll != len(values) {
		return errors.New("empty")
	}

	for i := 0; i < ll; i++ {
		key := keys[i]
		value := values[i]
		err := t.PutData(key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *MerkleTree) GetDatas(keys [][]byte) ([][]byte, error) {
	ll := len(keys)
	if ll == 0 {
		return nil, errors.New("null keys")
	}

	values := make([][]byte, ll)
	for i := 0; i < ll; i++ {
		key := keys[i]
		value, err := t.GetData(key)
		if err != nil {
			return nil, err
		}
		values[i] = value
	}
	return values, nil
}

func (t *MerkleTree) Flush() error {
	return t.Commit()
}

func (t *MerkleTree) HasData(key []byte) bool {
	value, err := t.GetData(key)
	if err != nil {
		return false
	}
	if len(value) == 0 {
		return false
	}
	return true
}

// TODO
func (t *MerkleTree) RemoveData(key []byte) error {
	return errors.New("unsupported")
}

func (t *MerkleTree) ListData(f func(key []byte, value []byte) error) error {
	return t.ss.ListData(f)
}

type MerkleService struct {
	index  string
	name   string
	config libcore.Config

	im *MerkleTree // index -> hash
	bm *MerkleTree // block
	tm *MerkleTree // transaction
	sm *MerkleTree // state

	crypto libcrypto.CryptoService
}

func NewMerkleService(c libcore.Config, cs libcrypto.CryptoService) (libstore.MerkleService, error) {
	service := &MerkleService{
		config: c,
		crypto: cs,
	}
	err := service.Init(c)
	if err != nil {
		return nil, err
	}
	err = service.Start()
	if err != nil {
		return nil, err
	}
	return service, nil
}

func (service *MerkleService) Init(c libcore.Config) error {
	service.config = c

	indexdb := &store.LevelService{Name: "index"}
	err := indexdb.Init(c)
	if err != nil {
		return err
	}
	err = indexdb.Start()
	if err != nil {
		return err
	}
	service.im = NewMerkleTree(service.crypto, indexdb)

	blockdb := &store.LevelService{Name: "block"}
	err = blockdb.Init(c)
	if err != nil {
		return err
	}
	err = blockdb.Start()
	if err != nil {
		return err
	}
	service.bm = NewMerkleTree(service.crypto, blockdb)

	txdb := &store.LevelService{Name: "transaction"}
	err = txdb.Init(c)
	if err != nil {
		return err
	}
	err = txdb.Start()
	if err != nil {
		return err
	}
	service.tm = NewMerkleTree(service.crypto, txdb)

	statedb := &store.LevelService{Name: "receipt"}
	err = statedb.Init(c)
	if err != nil {
		return err
	}
	err = statedb.Start()
	if err != nil {
		return err
	}
	service.sm = NewMerkleTree(service.crypto, statedb)
	return nil
}

func (service *MerkleService) Start() error {
	return nil
}

func (service *MerkleService) Close(s ...interface{}) error {
	err := service.im.Close()
	if err != nil {
		return err
	}
	err = service.bm.Close()
	if err != nil {
		return err
	}
	err = service.tm.Close()
	if err != nil {
		return err
	}
	err = service.sm.Close()
	if err != nil {
		return err
	}
	return nil
}

func (service *MerkleService) PutState(state libblock.State, s ...interface{}) error {
	cs := service.crypto

	h, data, err := cs.Raw(state, libcrypto.RawBinary)
	if err != nil {
		return err
	}
	err = service.sm.PutData(h, data)
	if err != nil {
		return err
	}

	stateHash := state.GetHash()
	stateTypeAndKey := getStateKeyWithType(state.GetStateKey(), state.GetStateType())
	stateTypeTypeAndAddress := getStateKeyWithType(state.GetAccount().String(), state.GetStateType())
	stateAddressAndIndexKey := getStateKey(fmt.Sprintf("%s:%d", state.GetAccount().String(), state.GetIndex()))
	stateAddressKey := getStateKey(state.GetAccount().String())

	err = service.im.PutData([]byte(stateTypeAndKey), stateHash)
	if err != nil {
		return err
	}
	err = service.im.PutData([]byte(stateTypeTypeAndAddress), stateHash)
	if err != nil {
		return err
	}
	err = service.im.PutData([]byte(stateAddressAndIndexKey), stateHash)
	if err != nil {
		return err
	}
	err = service.im.PutData([]byte(stateAddressKey), stateHash)
	if err != nil {
		return err
	}
	return nil
}

func (service *MerkleService) GetStateByHash(h libcore.Hash, s ...interface{}) (libblock.State, error) {
	data, err := service.sm.GetData(h)
	if err != nil {
		return nil, ErrorOfNonexists("state", h.String())
	}
	state, err := block.ReadState(data)
	if err != nil {
		return nil, err
	}
	return state, nil
}

func (service *MerkleService) GetStateByTypeAndKey(stateType libblock.StateType, stateKey string, s ...interface{}) (libblock.State, error) {
	stateTypeAndKey := getStateKeyWithType(stateKey, stateType)
	h, err := service.im.GetData([]byte(stateTypeAndKey))
	if err != nil {
		return nil, ErrorOfNonexists("state", stateTypeAndKey)
	}
	return service.GetStateByHash(libcore.Hash(h))
}

func (service *MerkleService) GetStateByTypeAndAddress(stateType libblock.StateType, account libcore.Address, s ...interface{}) (libblock.State, error) {
	stateTypeTypeAndAddress := getStateKeyWithType(account.String(), stateType)
	h, err := service.im.GetData([]byte(stateTypeTypeAndAddress))
	if err != nil {
		return nil, ErrorOfNonexists("state", stateTypeTypeAndAddress)
	}
	return service.GetStateByHash(libcore.Hash(h))
}

func (service *MerkleService) GetStateByAddressAndIndex(account libcore.Address, index uint64, s ...interface{}) (libblock.State, error) {
	stateAddressAndIndexKey := getStateKey(fmt.Sprintf("%s:%d", account.String(), index))
	h, err := service.im.GetData([]byte(stateAddressAndIndexKey))
	if err != nil {
		return nil, ErrorOfNonexists("state", stateAddressAndIndexKey)
	}
	return service.GetStateByHash(libcore.Hash(h))
}

func (service *MerkleService) GetStateByAddress(account libcore.Address, s ...interface{}) (libblock.State, error) {
	stateAddressKey := getStateKey(account.String())
	h, err := service.im.GetData([]byte(stateAddressKey))
	if err != nil {
		return nil, ErrorOfNonexists("state", stateAddressKey)
	}
	return service.GetStateByHash(libcore.Hash(h))
}

func (service *MerkleService) GetStateRoot() libcore.Hash {
	return service.sm.GetRoot()
}

func (service *MerkleService) PutTransaction(txWithData libblock.TransactionWithData, s ...interface{}) error {
	cs := service.crypto

	h := txWithData.GetHash()
	_, data, err := cs.Raw(txWithData, libcrypto.RawBinary)
	if err != nil {
		return err
	}
	err = service.tm.PutData(h, data)
	if err != nil {
		return err
	}

	account := txWithData.GetTransaction().GetAccount()
	accountKey := getTransactionKey(fmt.Sprintf("%s:%d", account.String(), txWithData.GetTransaction().GetIndex()))
	err = service.im.PutData([]byte(accountKey), h)
	if err != nil {
		return err
	}

	return nil
}

func (service *MerkleService) GetTransactionByHash(h libcore.Hash, s ...interface{}) (libblock.TransactionWithData, error) {
	data, err := service.tm.GetData(h)
	if err != nil {
		return nil, ErrorOfNonexists("transaction", h.String())
	}
	txWithData := &block.TransactionWithData{}
	err = txWithData.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}
	return txWithData, nil
}

func (service *MerkleService) GetTransactionByIndex(account libcore.Address, index uint64, s ...interface{}) (libblock.TransactionWithData, error) {
	accountKey := getTransactionKey(fmt.Sprintf("%s:%d", account.String(), index))
	h, err := service.im.GetData([]byte(accountKey))
	if err != nil {
		return nil, ErrorOfNonexists("transaction", fmt.Sprintf("%s, %d", account.String(), index))
	}
	return service.GetTransactionByHash(libcore.Hash(h))
}

func (service *MerkleService) GetTransactionRoot() libcore.Hash {
	return service.tm.GetRoot()
}

func (service *MerkleService) PutBlock(b libblock.Block, s ...interface{}) error {
	cs := service.crypto

	h, data, err := cs.Raw(b, libcrypto.RawBinary)
	if err != nil {
		return err
	}
	err = service.bm.PutData(h, data)
	if err != nil {
		return err
	}
	name := getBlockKey(b.GetIndex())
	err = service.im.PutData([]byte(name), h[:])
	if err != nil {
		return err
	}
	return nil
}

func (service *MerkleService) GetBlockByHash(hash libcore.Hash, s ...interface{}) (libblock.Block, error) {
	data, err := service.bm.GetData(hash)
	if err != nil {
		return nil, ErrorOfNonexists("block", hash.String())
	}
	b := &block.Block{}
	err = b.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (service *MerkleService) GetBlockByIndex(index uint64, s ...interface{}) (libblock.Block, error) {
	name := getBlockKey(index)
	data, err := service.im.GetData([]byte(name))
	if err != nil {
		return nil, ErrorOfNonexists("block", fmt.Sprintf("%d", index))
	}
	h := libcore.Hash(data)
	return service.GetBlockByHash(h)
}

func (service *MerkleService) Commit(s ...interface{}) error {
	err := service.im.Commit()
	if err != nil {
		return err
	}
	err = service.bm.Commit()
	if err != nil {
		return err
	}
	err = service.tm.Commit()
	if err != nil {
		return err
	}
	err = service.sm.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (service *MerkleService) Cancel(s ...interface{}) error {
	err := service.im.Cancel()
	if err != nil {
		return err
	}
	err = service.bm.Cancel()
	if err != nil {
		return err
	}
	err = service.tm.Cancel()
	if err != nil {
		return err
	}
	err = service.sm.Cancel()
	if err != nil {
		return err
	}
	return nil
}

func (service *MerkleService) Verify(key []byte, s ...interface{}) ([]byte, error) {
	return nil, errors.New("unsupport")
}

func getBlockKey(index uint64) string {
	return fmt.Sprintf("block@%d", index)
}

func getTransactionKey(key string) string {
	return fmt.Sprintf("transaction@%s", key)
}

func getStateKey(key string) string {
	return fmt.Sprintf("state@%s", key)
}

func getStateKeyWithType(key string, t libblock.StateType) string {
	return fmt.Sprintf("state@%s@%s", t.String(), key)
}

func ErrorOfNonexists(t string, target string) error {
	return fmt.Errorf("can't find %s: %s", t, target)
}
