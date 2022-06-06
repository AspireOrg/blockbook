package gasp

import (
	"bytes"

	"github.com/martinboehm/btcd/wire"
	"github.com/martinboehm/btcutil/chaincfg"
	"github.com/trezor/blockbook/bchain"
	"github.com/trezor/blockbook/bchain/coins/btc"
	"github.com/trezor/blockbook/bchain/coins/utils"
)

const (
	MainnetMagic wire.BitcoinNet = 0xe2cccfe4
	TestnetMagic wire.BitcoinNet = 0xe3cccfe4
)

var (
	MainNetParams chaincfg.Params
	TestNetParams chaincfg.Params
)

func init() {
	// Mainnet address
	MainNetParams = chaincfg.MainNetParams
	MainNetParams.Net = MainnetMagic
	MainNetParams.PubKeyHashAddrID = []byte{38}
	MainNetParams.ScriptHashAddrID = []byte{97}
	MainNetParams.Bech32HRPSegwit = "gasp"

	err := chaincfg.Register(&MainNetParams)
	if err != nil {
		panic(err)
	}

	// Testnet address
	TestNetParams = chaincfg.MainNetParams
	TestNetParams.Net = TestnetMagic
	TestNetParams.PubKeyHashAddrID = []byte{37}
	TestNetParams.ScriptHashAddrID = []byte{38}
	TestNetParams.Bech32HRPSegwit = "tgasp"

	err = chaincfg.Register(&TestNetParams)
	if err != nil {
		panic(err)
	}
}

// GaspParser handle
type GaspParser struct {
	*btc.BitcoinParser
}

// NewGaspParser returns new GaspParser instance
func NewGaspParser(params *chaincfg.Params, c *btc.Configuration) *GaspParser {
	return &GaspParser{BitcoinParser: btc.NewBitcoinParser(params, c)}
}

func GetChainParams(chain string) *chaincfg.Params {
	switch chain {
	case "testnet":
		return &TestNetParams
	default:
		return &MainNetParams
	}
}

// ParseBlock parses raw block to our Block struct
func (p *GaspParser) ParseBlock(b []byte) (*bchain.Block, error) {
	r := bytes.NewReader(b)
	w := wire.MsgBlock{}
	h := wire.BlockHeader{}
	err := h.Deserialize(r)
	if err != nil {
		return nil, err
	}

	err = utils.DecodeTransactions(r, 0, wire.WitnessEncoding, &w)
	if err != nil {
		return nil, err
	}

	txs := make([]bchain.Tx, len(w.Transactions))
	for ti, t := range w.Transactions {
		txs[ti] = p.TxFromMsgTx(t, false)
	}

	return &bchain.Block{
		BlockHeader: bchain.BlockHeader{
			Size: len(b),
			Time: h.Timestamp.Unix(),
		},
		Txs: txs,
	}, nil
}
