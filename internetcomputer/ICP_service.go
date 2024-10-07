package internetcomputer

import (
	"github.com/aviate-labs/agent-go"
	"github.com/aviate-labs/agent-go/ic"
	"github.com/aviate-labs/agent-go/identity"
)

type WalletService struct {
	agent *agent.Agent
}

type Account struct {
	Account string `ic:"account"`
}

type Balance struct {
	E8S uint64 `ic:"e8s"`
}

func NewWalletService(privateKey, publicKey []byte) (*WalletService, error){
	id, err := identity.NewEd25519Identity(privateKey, publicKey)
	if err != nil{
		return nil, err
	}

	config := agent.Config{
		Identity: id,
	}

	a, err := agent.New(config)
	if err != nil{
		return nil, err
	}

	return &WalletService{agent: a}, nil
}

func(ws *WalletService) GetBalance(accountID string) (uint64, error){
	var balance Balance

	if err := ws.agent.Query(
		ic.LEDGER_PRINCIPAL, "account_balance_dfx",
		[]any{Account{Account: accountID}},
		[]any{&balance},
	); err != nil{
		return 0, err
	}

	return balance.E8S, nil
}