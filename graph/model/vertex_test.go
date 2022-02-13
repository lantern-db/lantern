package model

import (
	"testing"
)

func TestEmptyVertex_IsNil(t *testing.T) {
	type fields struct {
		key        Key
		expiration Expiration
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		{
			"valid_case",
			fields{"key", Expiration(0)},
			true,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := EmptyVertex{
				key:        tt.fields.key,
				expiration: tt.fields.expiration,
			}
			got, err := e.IsNil()
			if (err != nil) != tt.wantErr {
				t.Errorf("IsNil() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsNil() got = %v, want %v", got, tt.want)
			}
		})
	}
}
