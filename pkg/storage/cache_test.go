package storage

import (
	"context"
	"encoding/json"
	"testing"
)

func TestStorage_CacheSet(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		storage *Storage
		key     string
		value   struct {
			Ping string
		}
		wantErr bool
	}{
		{
			name:    "cache-set-test",
			storage: MockStorage(ctx),
			key:     "corba",
			value:   struct{ Ping string }{"pong"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			js, _ := json.Marshal(tt.value)
			if err := tt.storage.CacheSet(ctx, tt.key, js); (err != nil) != tt.wantErr {
				t.Errorf("Storage.CacheSet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCacheGet(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		storage *Storage
		key     string
		wantErr bool
	}{
		{
			name:    "cache-get-test",
			storage: MockStorage(ctx),
			key:     "corba",
			wantErr: false,
		},
		{
			name:    "cache-get-test",
			storage: MockStorage(ctx),
			key:     "shouldfail",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.storage.CacheGet(ctx, tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.CacheGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
