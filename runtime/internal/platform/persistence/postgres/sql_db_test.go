package postgres

import "testing"

func TestOpenSQLDBFromPoolRejectsNilPool(t *testing.T) {
	_, err := OpenSQLDBFromPool(nil)
	if err == nil {
		t.Fatal("OpenSQLDBFromPool(nil) error = nil, want missing pool error")
	}
}
