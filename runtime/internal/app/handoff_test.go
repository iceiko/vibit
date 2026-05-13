package app

import "testing"

func TestRenderRouteKey(t *testing.T) {
	route := RouteKey{
		Kind:   MessageKindCommand,
		Module: " inventory ",
		Name:   " GrantItem ",
	}

	if got := RenderRouteKey(route); got != "inventory.GrantItem" {
		t.Fatalf("RenderRouteKey() = %q, want %q", got, "inventory.GrantItem")
	}
}

func TestRenderRouteKeyRequiresModuleAndName(t *testing.T) {
	if got := RenderRouteKey(RouteKey{Module: "inventory"}); got != "" {
		t.Fatalf("RenderRouteKey() = %q, want empty string", got)
	}

	if got := RenderRouteKey(RouteKey{Name: "GrantItem"}); got != "" {
		t.Fatalf("RenderRouteKey() = %q, want empty string", got)
	}
}
