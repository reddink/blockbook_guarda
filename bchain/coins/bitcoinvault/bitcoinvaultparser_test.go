// +build unittest

package bitcoinvault

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/martinboehm/btcutil/chaincfg"
	"github.com/trezor/blockbook/bchain"
	"github.com/trezor/blockbook/bchain/coins/btc"
)

func TestMain(m *testing.M) {
	c := m.Run()
	chaincfg.ResetParams()
	os.Exit(c)
}

func Test_GetAddrDescFromAddress_Mainnet(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "P2SH1",
			args:    args{address: "RDaCMJuJN7B5yykbNKfkTFqgQQPFqqnviw"},
			want:    "a9142f17cd421f2a4228a80ce0ffdcd037b80423057b87",
			wantErr: false,
		},
	}
	parser := NewBitcoinvaultParser(GetChainParams("main"), &btc.Configuration{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.GetAddrDescFromAddress(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAddrDescFromAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			h := hex.EncodeToString(got)
			if !reflect.DeepEqual(h, tt.want) {
				t.Errorf("GetAddrDescFromAddress() = %v, want %v", h, tt.want)
			}
		})
	}
}

func Test_GetAddressesFromAddrDesc_Mainnet(t *testing.T) {
	type args struct {
		script string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		want2   bool
		wantErr bool
	}{
		{
			name:    "P2SH1",
			args:    args{script: "a9140e692ebf126dc9fd945d6f90524a150d7a5921c387"},
			want:    []string{"RAbPYR9GedFf56Fuyzu4cieBW5CG6Ybeoy"},
			want2:   true,
			wantErr: false,
		},
		{
			name:    "OP_RETURN ascii",
			args:    args{script: "6a0461686f6a"},
			want:    []string{"OP_RETURN (ahoj)"},
			want2:   false,
			wantErr: false,
		},
		{
			name:    "OP_RETURN hex",
			args:    args{script: "6a072020f1686f6a20"},
			want:    []string{"OP_RETURN 2020f1686f6a20"},
			want2:   false,
			wantErr: false,
		},
	}

	parser := NewBitcoinvaultParser(GetChainParams("main"), &btc.Configuration{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := hex.DecodeString(tt.args.script)
			got, got2, err := parser.GetAddressesFromAddrDesc(b)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAddressesFromAddrDesc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAddressesFromAddrDesc() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("GetAddressesFromAddrDesc() = %v, want %v", got2, tt.want2)
			}
		})
	}
}

var (
	testTx1       bchain.Tx
	testTxPacked1 = "00030e6d8ba8d7aa20020000000001010914dec60b361360f04cf43b12063ae3bb2532facb640ef030c358453c8f3bc90100000017160014f6eb6d8f7413e5da585b82e2213291d0600e65f7feffffff02d00705000000000017a914b9ee7e29e9fa3209e2ebc232280a3ee046ea66448744ce01000000000017a91442c360a174b3e07e1abca0ad6fe83d22c26d62fb8702473044022075f19d2c93e23936bc271470eb41f86db694ec3127f1de1294464219a1cf220402200ca5926dfd3ff8607e3633cb58b57625014a5d5478a1269c2078d59f3e3145e50121038c629a5b4bbffacb322b36d88e4b00af43a95a2e08b13dc69d61621035f265b6f5250000"
)

func init() {
	testTx1 = bchain.Tx{
		Hex:       "020000000001010914dec60b361360f04cf43b12063ae3bb2532facb640ef030c358453c8f3bc90100000017160014f6eb6d8f7413e5da585b82e2213291d0600e65f7feffffff02d00705000000000017a914b9ee7e29e9fa3209e2ebc232280a3ee046ea66448744ce01000000000017a91442c360a174b3e07e1abca0ad6fe83d22c26d62fb8702473044022075f19d2c93e23936bc271470eb41f86db694ec3127f1de1294464219a1cf220402200ca5926dfd3ff8607e3633cb58b57625014a5d5478a1269c2078d59f3e3145e50121038c629a5b4bbffacb322b36d88e4b00af43a95a2e08b13dc69d61621035f265b6f5250000",
		Blocktime: 1519053456,
		Txid:      "9040141a87ab9f0e40ddb049818ebedc8c6124e02631144e83d2d52d76a3f588",
		LockTime:  9717,
		Version:   2,
		Vin: []bchain.Vin{
			{
				ScriptSig: bchain.ScriptSig{
					Hex: "160014f6eb6d8f7413e5da585b82e2213291d0600e65f7",
				},
				Txid:     "c93b8f3c4558c330f00e64cbfa3225bbe33a06123bf44cf06013360bc6de1409",
				Vout:     1,
				Sequence: 4294967294,
			},
		},
		Vout: []bchain.Vout{
			{
				ValueSat: *big.NewInt(329680),
				N:        0,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "a914b9ee7e29e9fa3209e2ebc232280a3ee046ea664487",
					Addresses: []string{
						"RSEJnn1UMc3vNENY3FpFBSq24juafHiQA9",
					},
				},
			},
			{
				ValueSat: *big.NewInt(118340),
				N:        1,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "a91442c360a174b3e07e1abca0ad6fe83d22c26d62fb87",
					Addresses: []string{
						"RFNCjEw5MPunUwcqqCuh97zxMH4xCQtyEJ",
					},
				},
			},
		},
	}
}

func Test_PackTx(t *testing.T) {
	type args struct {
		tx        bchain.Tx
		height    uint32
		blockTime int64
		parser    *BitcoinvaultParser
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "bitcoinvault-1",
			args: args{
				tx:        testTx1,
				height:    200301,
				blockTime: 1519053456,
				parser:    NewBitcoinvaultParser(GetChainParams("main"), &btc.Configuration{}),
			},
			want:    testTxPacked1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.parser.PackTx(&tt.args.tx, tt.args.height, tt.args.blockTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("packTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			h := hex.EncodeToString(got)
			if !reflect.DeepEqual(h, tt.want) {
				t.Errorf("packTx() = %v, want %v", h, tt.want)
			}
		})
	}
}

func Test_UnpackTx(t *testing.T) {
	type args struct {
		packedTx string
		parser   *BitcoinvaultParser
	}
	tests := []struct {
		name    string
		args    args
		want    *bchain.Tx
		want1   uint32
		wantErr bool
	}{
		{
			name: "bitcoinvault-1",
			args: args{
				packedTx: testTxPacked1,
				parser:   NewBitcoinvaultParser(GetChainParams("main"), &btc.Configuration{}),
			},
			want:    &testTx1,
			want1:   200301,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := hex.DecodeString(tt.args.packedTx)
			got, got1, err := tt.args.parser.UnpackTx(b)
			if (err != nil) != tt.wantErr {
				t.Errorf("unpackTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unpackTx() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("unpackTx() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

type testBlock struct {
	size int
	time int64
	txs  []string
}

var testParseBlockTxs = map[int]testBlock{
	3125: {
		size: 530,
		time: 1572258327,
		txs: []string{
			"efde5fde3eb383565c2a8108cd2859b0fb294ea8936c423804931d4329c21d18",
			"4e68c6914588f327da24baac51b2452afbb80da51f17d6645887b648a315ba1e",
		},
	},
	9718: {
		size: 2195,
		time: 1576260579,
		txs: []string{
			"ec3b5fbaeb44e346a36c66e82c450a2b5d00544e96ab755bbb59408d05bf8a56",
			"5c9574f5b58740976fac14fe6f2213a89c13b71d9810e94c2fcdbf8860a27a7f",
			"06194f7a2699908b79ab24af017cc12649bedade7e8ceb25f091ae5e6857e2f9",
			"c9f17e73e60b3c27776e7d5bebf6ce9f613393489f5cd1872d463e0f9e9ebb8e",
			"d52076f112ef6e70adad5a44de392aeab87df20c6fa35d9f0832f13904c57bd4",
			"c680ae8d6ec7adabe84524f1d7296921e65adef8580bc22cffa741036e1f414c",
			"ade832242fe7df9a0bba6acebd93f5c9f4258e09a9ed7045c3c276ea3e71ab2e",
			"9040141a87ab9f0e40ddb049818ebedc8c6124e02631144e83d2d52d76a3f588",
		},
	},
	23359: {
		size: 195,
		time: 1584445598,
		txs: []string{
			"25223af61f0a6a65ad9f2a98a127abca7c0b16efffdf20dac8a02407f934c811",
		},
	},
	58420: {
		size: 2662,
		time: 1605617694,
		txs: []string{
			"775a586cf0fd7e1a639920b263e0a78c82533321c6ce4ac48fb08dd09f7a154c",
			"8cbe9de71975a2a0dfbd75a8c3a7f3bec61d200bbf89d5c5928f89e271a68ce3",
			"7815dfb7c774c23b95495bf370bc3ff692cc5ba823b53ac263cea85e315a1893",
			"55d6d1329497230bfc8ccc2cdb77567ab96742926f43d6231e731f7067b41961",
			"f71caddab966605632dbd44bb76e67f572bf3671c21fb4d4bf0a8d13ffa92a8e",
			"89adb12aa09ea6ce5b526bd53320dcc7b243157c9296187b8c61d4865d9c1f96",
			"512a48d730784eb1b36aaf2f8f291bf6f7f1ec0a95999042f7ddc3f413aedb03",
			"24c0c22f514a2a5ad105f7cce712f29b697bd225fa32fa83e77afb1d03f97c29",
		},
	},
}

func helperLoadBlock(t *testing.T, height int) []byte {
	name := fmt.Sprintf("block_dump.%d", height)
	path := filepath.Join("testdata", name)

	d, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	d = bytes.TrimSpace(d)

	b := make([]byte, hex.DecodedLen(len(d)))
	_, err = hex.Decode(b, d)
	if err != nil {
		t.Fatal(err)
	}

	return b
}

func TestParseBlock(t *testing.T) {
	p := NewBitcoinvaultParser(GetChainParams("main"), &btc.Configuration{})

	for height, tb := range testParseBlockTxs {
		b := helperLoadBlock(t, height)

		blk, err := p.ParseBlock(b)
		if err != nil {
			t.Fatal(err)
		}

		if blk.Size != tb.size {
			t.Errorf("ParseBlock() block size: got %d, want %d", blk.Size, tb.size)
		}

		if blk.Time != tb.time {
			t.Errorf("ParseBlock() block time: got %d, want %d", blk.Time, tb.time)
		}

		if len(blk.Txs) != len(tb.txs) {
			t.Errorf("ParseBlock() number of transactions: got %d, want %d", len(blk.Txs), len(tb.txs))
		}

		for ti, tx := range tb.txs {
			if blk.Txs[ti].Txid != tx {
				t.Errorf("ParseBlock() transaction %d: got %s, want %s", ti, blk.Txs[ti].Txid, tx)
			}
		}
	}
}
