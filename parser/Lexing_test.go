package parser

import (
	"reflect"
	"testing"
)

func Test_lexText(t *testing.T) {
	type args struct {
		lex *Lexer
	}
	tests := []struct {
		name string
		args args
		want StateFn
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lexText(tt.args.lex); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lexText() = %v, want %v", got, tt.want)
			}
		})
	}
}
