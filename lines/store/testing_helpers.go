package store

import (
	"testing"
)

// IsolatedIntegrationTest is an integration test wrapper for running tests.
// It sets up the test environment, runs the test function, and rolls back the transaction.
func IsolatedIntegrationTest(t *testing.T, stores []IntegrationTestStore, testFunc func(*testing.T)) {
	for _, store := range stores {
		err := store.BeginTransaction()
		if err != nil {
			t.Fatalf("Failed to start transaction for integration test: %v", err)
		}
	}

	defer func() {
		for _, store := range stores {
			err := store.RollbackTransaction()
			if err != nil {
				t.Fatalf("Failed to rollback transaction for integration test: %v", err)
			}
		}
	}()
	testFunc(t)
}
