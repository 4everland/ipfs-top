package arweave

import (
	"encoding/hex"
	arseedSchema "github.com/everFinance/arseeding/schema"
	"github.com/everFinance/arseeding/sdk/schema"
	"github.com/everFinance/goar"
	"github.com/everFinance/goar/types"
	"github.com/everFinance/goether"
	lru "github.com/hashicorp/golang-lru"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"math/big"
)

type SDK struct {
	ItemSigner ItemSignerFactory
	Cli        *ArSeedCli
	NoFee      bool
}

func NewSDK(arseedUrl, mnemonic string) (*SDK, error) {
	return &SDK{
		ItemSigner: NewItemSignerFactory(mnemonic),
		Cli:        New(arseedUrl),
	}, nil
}

func (s *SDK) SendData(uid string, data []byte, option *schema.OptionItem) (order *arseedSchema.RespOrder, size int, err error) {
	var bundleItem types.BundleItem
	if option != nil {
		bundleItem, err = s.ItemSigner.GetItemSigner(uid).CreateAndSignItem(data, option.Target, option.Anchor, option.Tags)
	} else {
		bundleItem, err = s.ItemSigner.GetItemSigner(uid).CreateAndSignItem(data, "", "", nil)
	}
	if err != nil {
		return
	}
	size = len(bundleItem.ItemBinary)
	order, err = s.Cli.SubmitItem(bundleItem.ItemBinary, "0")

	return
}

func (s *SDK) GetWalletBalance(address string) (*big.Float, error) {
	return s.Cli.ACli.GetWalletBalance(address)
}

func NewItemSignerFactory(mnemonic string) ItemSignerFactory {
	cache, _ := lru.New(100)

	return &singer{mnemonic: mnemonic, cache: cache}
}

type ItemSignerFactory interface {
	GetItemSigner(secret string) *goar.ItemSigner
}

type singer struct {
	mnemonic string
	cache    *lru.Cache
}

func (s *singer) GetItemSigner(secret string) *goar.ItemSigner {
	if itemSigner, ok := s.cache.Get(secret); ok {
		return itemSigner.(*goar.ItemSigner)
	}

	seed := bip39.NewSeed(s.mnemonic, secret)
	privateKey, _ := bip32.NewMasterKey(seed)

	privateKeyHex := hex.EncodeToString(privateKey.Key)
	signer, _ := goether.NewSigner(privateKeyHex)
	itemSigner, _ := goar.NewItemSigner(signer)
	s.cache.Add(secret, itemSigner)
	return itemSigner
}
