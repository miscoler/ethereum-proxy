package blocktsx

import "encoding/json"

type Transaction struct {
	Hash string `json:"hash"`
}

type Block struct {
	Transactions []json.RawMessage `json:"transactions"`
}

type BlockStored struct {
	TransactionsByHash map[string]*json.RawMessage
	Transactions       []json.RawMessage
}
