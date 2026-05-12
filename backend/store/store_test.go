package store

import (
	"context"
	"testing"
)

func TestNewStore_Error(t *testing.T) {
	// Use an invalid connection string to trigger error and cover some lines
	t.Setenv("DATABASE_URL", "invalid-url")
	
	ctx := context.Background()
	_, err := NewStore(ctx)
	
	if err == nil {
		t.Error("Expected error for invalid DATABASE_URL")
	}
}
