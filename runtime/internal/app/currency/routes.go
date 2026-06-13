package currency

import (
	"github.com/iceiko/vibit/runtime/internal/app"
	currencymodule "github.com/iceiko/vibit/runtime/internal/modules/currency"
)

const (
	CommandEnsurePlayerWallet      = "EnsurePlayerWallet"
	CommandGrantCurrency           = "GrantCurrency"
	CommandSpendCurrency           = "SpendCurrency"
	QueryGetOwnWallet              = "GetOwnWallet"
	QueryListOwnWalletBalances     = "ListOwnWalletBalances"
	QueryListOwnWalletTransactions = "ListOwnWalletTransactions"
)

// Full route keys:
// - currency.EnsurePlayerWallet
// - currency.GetOwnWallet
// - currency.ListOwnWalletBalances
// - currency.GrantCurrency
// - currency.SpendCurrency
// - currency.ListOwnWalletTransactions

func EnsurePlayerWalletRoute() app.RouteKey {
	return app.RouteKey{Kind: app.MessageKindCommand, Module: currencymodule.ModuleName, Name: CommandEnsurePlayerWallet}
}

func GetOwnWalletRoute() app.RouteKey {
	return app.RouteKey{Kind: app.MessageKindQuery, Module: currencymodule.ModuleName, Name: QueryGetOwnWallet}
}

func ListOwnWalletBalancesRoute() app.RouteKey {
	return app.RouteKey{Kind: app.MessageKindQuery, Module: currencymodule.ModuleName, Name: QueryListOwnWalletBalances}
}

func GrantCurrencyRoute() app.RouteKey {
	return app.RouteKey{Kind: app.MessageKindCommand, Module: currencymodule.ModuleName, Name: CommandGrantCurrency}
}

func SpendCurrencyRoute() app.RouteKey {
	return app.RouteKey{Kind: app.MessageKindCommand, Module: currencymodule.ModuleName, Name: CommandSpendCurrency}
}

func ListOwnWalletTransactionsRoute() app.RouteKey {
	return app.RouteKey{Kind: app.MessageKindQuery, Module: currencymodule.ModuleName, Name: QueryListOwnWalletTransactions}
}
