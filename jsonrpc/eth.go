package jsonrpc

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/umbracle/ethgo"
)

// Eth is the eth namespace
type Eth struct {
	c       *Client
	chainId *big.Int
}

// Eth returns the reference to the eth namespace
func (c *Client) Eth() *Eth {
	return c.endpoints.e
}

func (e *Eth) GetNodeInfo() (string, error) {
	var res string
	err := e.c.Call("web3_clientVersion", &res)
	return res, err
}

// GetCode returns the code of a contract
func (e *Eth) GetCode(addr ethgo.Address, block ethgo.BlockNumberOrHash) (string, error) {
	var res string
	if err := e.c.Call("eth_getCode", &res, addr, block.Location()); err != nil {
		return "", err
	}
	return res, nil
}

// Accounts returns a list of addresses owned by client.
func (e *Eth) Accounts() ([]ethgo.Address, error) {
	var out []ethgo.Address
	if err := e.c.Call("eth_accounts", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetStorageAt returns the value from a storage position at a given address.
func (e *Eth) GetStorageAt(addr ethgo.Address, slot ethgo.Hash, block ethgo.BlockNumberOrHash) (ethgo.Hash, error) {
	var hash ethgo.Hash
	err := e.c.Call("eth_getStorageAt", &hash, addr, slot, block.Location())
	return hash, err
}

// BlockNumber returns the number of most recent block.
func (e *Eth) BlockNumber() (uint64, error) {
	var out string
	if err := e.c.Call("eth_blockNumber", &out); err != nil {
		return 0, err
	}
	return parseUint64orHex(out)
}

// GetBlockByNumber returns information about a block by block number.
func (e *Eth) GetBlockByNumber(i ethgo.BlockNumber, full bool) (*ethgo.Block, error) {
	var b *ethgo.Block
	if err := e.c.Call("eth_getBlockByNumber", &b, i.String(), full); err != nil {
		return nil, err
	}
	return b, nil
}

// GetBlockByHash returns information about a block by hash.
func (e *Eth) GetBlockByHash(hash ethgo.Hash, full bool) (*ethgo.Block, error) {
	var b *ethgo.Block
	if err := e.c.Call("eth_getBlockByHash", &b, hash, full); err != nil {
		return nil, err
	}
	return b, nil
}

// GetFilterChanges returns the filter changes for log filters
func (e *Eth) GetFilterChanges(id string) ([]*ethgo.Log, error) {
	var raw string
	err := e.c.Call("eth_getFilterChanges", &raw, id)
	if err != nil {
		return nil, err
	}
	var res []*ethgo.Log
	if err := json.Unmarshal([]byte(raw), &res); err != nil {
		return nil, err
	}
	return res, nil
}

// GetTransactionByHash returns a transaction by his hash
func (e *Eth) GetTransactionByHash(hash ethgo.Hash) (*ethgo.Transaction, error) {
	var txn *ethgo.Transaction
	err := e.c.Call("eth_getTransactionByHash", &txn, hash)
	return txn, err
}

// GetFilterChangesBlock returns the filter changes for block filters
func (e *Eth) GetFilterChangesBlock(id string) ([]ethgo.Hash, error) {
	var raw string
	err := e.c.Call("eth_getFilterChanges", &raw, id)
	if err != nil {
		return nil, err
	}
	var res []ethgo.Hash
	if err := json.Unmarshal([]byte(raw), &res); err != nil {
		return nil, err
	}
	return res, nil
}

// NewFilter creates a new log filter
func (e *Eth) NewFilter(filter *ethgo.LogFilter) (string, error) {
	var id string
	err := e.c.Call("eth_newFilter", &id, filter)
	return id, err
}

// NewBlockFilter creates a new block filter
func (e *Eth) NewBlockFilter() (string, error) {
	var id string
	err := e.c.Call("eth_newBlockFilter", &id, nil)
	return id, err
}

// UninstallFilter uninstalls a filter
func (e *Eth) UninstallFilter(id string) (bool, error) {
	var res bool
	err := e.c.Call("eth_uninstallFilter", &res, id)
	return res, err
}

// SendRawTransaction sends a signed transaction in rlp format.
func (e *Eth) SendRawTransaction(data []byte) (ethgo.Hash, error) {
	var hash ethgo.Hash
	hexData := "0x" + hex.EncodeToString(data)
	err := e.c.Call("eth_sendRawTransaction", &hash, hexData)
	return hash, err
}

// SendTransaction creates new message call transaction or a contract creation.
func (e *Eth) SendTransaction(txn *ethgo.Transaction) (ethgo.Hash, error) {
	var hash ethgo.Hash
	err := e.c.Call("eth_sendTransaction", &hash, txn)
	return hash, err
}

// GetTransactionReceipt returns the receipt of a transaction by transaction hash.
func (e *Eth) GetTransactionReceipt(hash ethgo.Hash) (*ethgo.Receipt, error) {
	var receipt *ethgo.Receipt
	err := e.c.Call("eth_getTransactionReceipt", &receipt, hash)
	return receipt, err
}

// GetNonce returns the nonce of the account
func (e *Eth) GetNonce(addr ethgo.Address, blockNumber ethgo.BlockNumberOrHash) (uint64, error) {
	var nonce string
	if err := e.c.Call("eth_getTransactionCount", &nonce, addr, blockNumber.Location()); err != nil {
		return 0, err
	}
	return parseUint64orHex(nonce)
}

// GetBalance returns the balance of the account of given address.
func (e *Eth) GetBalance(addr ethgo.Address, blockNumber ethgo.BlockNumberOrHash) (*big.Int, error) {
	var out string
	if err := e.c.Call("eth_getBalance", &out, addr, blockNumber.Location()); err != nil {
		return nil, err
	}
	b, ok := new(big.Int).SetString(out[2:], 16)
	if !ok {
		return nil, fmt.Errorf("failed to convert to big.int")
	}
	return b, nil
}

// GasPrice returns the current price per gas in wei.
func (e *Eth) GasPrice() (uint64, error) {
	var out string
	if err := e.c.Call("eth_gasPrice", &out); err != nil {
		return 0, err
	}
	return parseUint64orHex(out)
}

// Call executes a new message call immediately without creating a transaction on the block chain.
func (e *Eth) Call(msg *ethgo.CallMsg, block ethgo.BlockNumber) (string, error) {
	var out string
	if err := e.c.Call("eth_call", &out, msg, block.String()); err != nil {
		return "", err
	}
	return out, nil
}

// EstimateGasContract estimates the gas to deploy a contract
func (e *Eth) EstimateGasContract(bin []byte) (uint64, error) {
	var out string
	msg := map[string]interface{}{
		"data": "0x" + hex.EncodeToString(bin),
	}
	if err := e.c.Call("eth_estimateGas", &out, msg); err != nil {
		return 0, err
	}
	return parseUint64orHex(out)
}

// EstimateGas generates and returns an estimate of how much gas is necessary to allow the transaction to complete.
func (e *Eth) EstimateGas(msg *ethgo.CallMsg) (uint64, error) {
	var out string
	if err := e.c.Call("eth_estimateGas", &out, msg); err != nil {
		return 0, err
	}
	return parseUint64orHex(out)
}

// GetLogs returns an array of all logs matching a given filter object
func (e *Eth) GetLogs(filter *ethgo.LogFilter) ([]*ethgo.Log, error) {
	var out []*ethgo.Log
	if err := e.c.Call("eth_getLogs", &out, filter); err != nil {
		return nil, err
	}
	return out, nil
}

// ChainID returns the id of the chain
func (e *Eth) ChainID() (*big.Int, error) {
	if e.chainId != nil {
		return e.chainId, nil
	}
	var out string
	if err := e.c.Call("eth_chainId", &out); err != nil {
		return nil, err
	}
	chainId := parseBigInt(out)
	e.chainId = chainId
	return chainId, nil
}

/**
 var methods = [
        new Method({
            name: 'getNodeInfo',
            call: 'web3_clientVersion'
        }),
        new Method({
            name: 'getProtocolVersion',
            call: 'eth_protocolVersion',
            params: 0
        }),
        new Method({
            name: 'getCoinbase',
            call: 'eth_coinbase',
            params: 0
        }),
        new Method({
            name: 'isMining',
            call: 'eth_mining',
            params: 0
        }),
        new Method({
            name: 'getHashrate',
            call: 'eth_hashrate',
            params: 0,
            outputFormatter: utils.hexToNumber
        }),
        new Method({
            name: 'isSyncing',
            call: 'eth_syncing',
            params: 0,
            outputFormatter: formatter.outputSyncingFormatter
        }),
        new Method({
            name: 'getGasPrice',
            call: 'eth_gasPrice',
            params: 0,
            outputFormatter: formatter.outputBigNumberFormatter
        }),
        new Method({
            name: 'getAccounts',
            call: 'eth_accounts',
            params: 0,
            outputFormatter: utils.toChecksumAddress
        }),
        new Method({
            name: 'getBlockNumber',
            call: 'eth_blockNumber',
            params: 0,
            outputFormatter: utils.hexToNumber
        }),
        new Method({
            name: 'getBalance',
            call: 'eth_getBalance',
            params: 2,
            inputFormatter: [formatter.inputAddressFormatter, formatter.inputDefaultBlockNumberFormatter],
            outputFormatter: formatter.outputBigNumberFormatter
        }),
        new Method({
            name: 'getStorageAt',
            call: 'eth_getStorageAt',
            params: 3,
            inputFormatter: [formatter.inputAddressFormatter, utils.numberToHex, formatter.inputDefaultBlockNumberFormatter]
        }),
        new Method({
            name: 'getCode',
            call: 'eth_getCode',
            params: 2,
            inputFormatter: [formatter.inputAddressFormatter, formatter.inputDefaultBlockNumberFormatter]
        }),
        new Method({
            name: 'getBlock',
            call: blockCall,
            params: 2,
            inputFormatter: [formatter.inputBlockNumberFormatter, function (val) { return !!val; }],
            outputFormatter: formatter.outputBlockFormatter
        }),
        new Method({
            name: 'getUncle',
            call: uncleCall,
            params: 2,
            inputFormatter: [formatter.inputBlockNumberFormatter, utils.numberToHex],
            outputFormatter: formatter.outputBlockFormatter,
        }),
        new Method({
            name: 'getBlockTransactionCount',
            call: getBlockTransactionCountCall,
            params: 1,
            inputFormatter: [formatter.inputBlockNumberFormatter],
            outputFormatter: utils.hexToNumber
        }),
        new Method({
            name: 'getBlockUncleCount',
            call: uncleCountCall,
            params: 1,
            inputFormatter: [formatter.inputBlockNumberFormatter],
            outputFormatter: utils.hexToNumber
        }),
        new Method({
            name: 'getTransaction',
            call: 'eth_getTransactionByHash',
            params: 1,
            inputFormatter: [null],
            outputFormatter: formatter.outputTransactionFormatter
        }),
        new Method({
            name: 'getTransactionFromBlock',
            call: transactionFromBlockCall,
            params: 2,
            inputFormatter: [formatter.inputBlockNumberFormatter, utils.numberToHex],
            outputFormatter: formatter.outputTransactionFormatter
        }),
        new Method({
            name: 'getTransactionReceipt',
            call: 'eth_getTransactionReceipt',
            params: 1,
            inputFormatter: [null],
            outputFormatter: formatter.outputTransactionReceiptFormatter
        }),
        new Method({
            name: 'getTransactionCount',
            call: 'eth_getTransactionCount',
            params: 2,
            inputFormatter: [formatter.inputAddressFormatter, formatter.inputDefaultBlockNumberFormatter],
            outputFormatter: utils.hexToNumber
        }),
        new Method({
            name: 'sendSignedTransaction',
            call: 'eth_sendRawTransaction',
            params: 1,
            inputFormatter: [null],
            abiCoder: abi
        }),

        new Method({
            name: 'signTransaction',
            call: 'eth_signTransaction',
            params: 1,
            inputFormatter: [formatter.inputTransactionFormatter]
        }),
        new Method({
            name: 'sendTransaction',
            call: 'eth_sendTransaction',
            params: 1,
            inputFormatter: [formatter.inputTransactionFormatter],
            abiCoder: abi
        }),
        new Method({
            name: 'sign',
            call: 'eth_sign',
            params: 2,
            inputFormatter: [formatter.inputSignFormatter, formatter.inputAddressFormatter],
            transformPayload: function (payload) {
                payload.params.reverse();
                return payload;
            }
        }),
        new Method({
            name: 'call',
            call: 'eth_call',
            params: 2,
            inputFormatter: [formatter.inputCallFormatter, formatter.inputDefaultBlockNumberFormatter],
            abiCoder: abi
        }),
        new Method({
            name: 'estimateGas',
            call: 'eth_estimateGas',
            params: 1,
            inputFormatter: [formatter.inputCallFormatter],
            outputFormatter: utils.hexToNumber
        }),
        new Method({
            name: 'submitWork',
            call: 'eth_submitWork',
            params: 3
        }),
        new Method({
            name: 'getWork',
            call: 'eth_getWork',
            params: 0
        }),
        new Method({
            name: 'getPastLogs',
            call: 'eth_getLogs',
            params: 1,
            inputFormatter: [formatter.inputLogFormatter],
            outputFormatter: formatter.outputLogFormatter
        }),
        new Method({
            name: 'getChainId',
            call: 'eth_chainId',
            params: 0,
            outputFormatter: utils.hexToNumber
        }),
        new Method({
            name: 'requestAccounts',
            call: 'eth_requestAccounts',
            params: 0,
            outputFormatter: utils.toChecksumAddress
        }),
        new Method({
            name: 'getProof',
            call: 'eth_getProof',
            params: 3,
            inputFormatter: [formatter.inputAddressFormatter, formatter.inputStorageKeysFormatter, formatter.inputDefaultBlockNumberFormatter],
            outputFormatter: formatter.outputProofFormatter
        }),
        new Method({
            name: 'getPendingTransactions',
            call: 'eth_pendingTransactions',
            params: 0,
            outputFormatter: formatter.outputTransactionFormatter
        }),


**/
