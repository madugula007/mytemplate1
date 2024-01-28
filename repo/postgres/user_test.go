package repository

import (
	"reflect"
	"testing"
)

func Test_generateColumnsFromStruct(t *testing.T) {
	type args struct {
		instance interface{}
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateColumnsFromStruct(tt.args.instance); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateColumnsFromStruct() = %v, want %v", got, tt.want)
			}
		})
	}
}
