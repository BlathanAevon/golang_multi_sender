package wallet

import (
	"crypto/ecdsa"
	"encoding/hex"
	"strings"

	"github.com/BlathanAevon/MultiSender/internal/client"
	"github.com/BlathanAevon/MultiSender/tools"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	Address    common.Address
}

func (w *Wallet) SendNative(address string, c *client.Rpc, amount float64) (common.Hash, error) {

	nonce, err := c.GetNonce(w.Address)

	if err != nil {
		return common.Hash{}, err
	}

	value := tools.FloatToWei(amount)
	gasLimit := uint64(21000)

	gasPrice, err := c.GetGasPrice()

	if err != nil {
		return common.Hash{}, err
	}

	toAddress := common.HexToAddress(address)

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	chainId, err := c.GetChainId()

	if err != nil {
		return common.Hash{}, err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), w.PrivateKey)

	if err != nil {
		return common.Hash{}, err
	}

	hash, err := c.SendTx(signedTx)

	if err != nil {
		return common.Hash{}, err
	}

	return hash, nil

}

func NewWallet(key string) (*Wallet, error) {

	key = strings.Replace(key, "0x", "", 1)

	privateKeyBytes, err := hex.DecodeString(key)

	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)

	w := &Wallet{PrivateKey: privateKey}

	if err != nil {
		return nil, err
	}

	public := privateKey.PublicKey

	w.Address = crypto.PubkeyToAddress(public)

	return w, nil

}
