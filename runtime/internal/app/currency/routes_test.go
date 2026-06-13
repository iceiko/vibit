package currency

import (
	"testing"

	"github.com/iceiko/vibit/runtime/internal/app"
	currencymodule "github.com/iceiko/vibit/runtime/internal/modules/currency"
)

func TestCurrencyWalletRoutesUseExpectedKindsAndNames(t *testing.T) {
	tests := []struct {
		name string
		got  app.RouteKey
		want app.RouteKey
	}{
		{
			name: "ensure player wallet command",
			got:  EnsurePlayerWalletRoute(),
			want: app.RouteKey{Kind: app.MessageKindCommand, Module: currencymodule.ModuleName, Name: CommandEnsurePlayerWallet},
		},
		{
			name: "get own wallet query",
			got:  GetOwnWalletRoute(),
			want: app.RouteKey{Kind: app.MessageKindQuery, Module: currencymodule.ModuleName, Name: QueryGetOwnWallet},
		},
		{
			name: "list own wallet balances query",
			got:  ListOwnWalletBalancesRoute(),
			want: app.RouteKey{Kind: app.MessageKindQuery, Module: currencymodule.ModuleName, Name: QueryListOwnWalletBalances},
		},
		{
			name: "grant currency command",
			got:  GrantCurrencyRoute(),
			want: app.RouteKey{Kind: app.MessageKindCommand, Module: currencymodule.ModuleName, Name: CommandGrantCurrency},
		},
		{
			name: "spend currency command",
			got:  SpendCurrencyRoute(),
			want: app.RouteKey{Kind: app.MessageKindCommand, Module: currencymodule.ModuleName, Name: CommandSpendCurrency},
		},
		{
			name: "list own wallet transactions query",
			got:  ListOwnWalletTransactionsRoute(),
			want: app.RouteKey{Kind: app.MessageKindQuery, Module: currencymodule.ModuleName, Name: QueryListOwnWalletTransactions},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.got != tc.want {
				t.Fatalf("route = %#v, want %#v", tc.got, tc.want)
			}
			if rendered := app.RenderRouteKey(tc.got); rendered == "" {
				t.Fatalf("RenderRouteKey(%#v) = empty, want currency route id", tc.got)
			}
		})
	}
}
