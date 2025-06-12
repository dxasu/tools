package main

import (
	"reflect"
	"strings"
	"testing"

	_ "github.com/dxasu/pure/version"
)

func Test_parseTokens(t *testing.T) {
	type args struct {
		expr string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "simple expression",
			args: args{expr: "8000*(1+1.2)/2 * 30 * 1%3"},
			want: []string{"3", "+", "5"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseTokens(tt.args.expr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTokens() = %v, want %v", strings.Join(got, ","), tt.want)
			}
		})
	}
}

func Test_infixToPostfix(t *testing.T) {
	type args struct {
		expr string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{{
		name: "simple expression",
		args: args{expr: "1+600%"},
		want: []string{"3", "+", "5"},
	}, {
		name: "simple expression",
		args: args{expr: "8000*(1+1.2)/2 * 30 * 1%3"},
		want: []string{"3", "+", "5"},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := infixToPostfix(tt.args.expr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("infixToPostfix() = %v, want %v", strings.Join(got, ","), tt.want)
			}
		})
	}
}

func Test_evaluatePostfix(t *testing.T) {
	type args struct {
		tokens []string
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "simple expression",
			args: args{tokens: []string{"3", "5", "+"}},
			want: 8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := evaluatePostfix(tt.args.tokens); got != tt.want {
				t.Errorf("evaluatePostfix() = %v, want %v", got, tt.want)
			}
		})
	}
}
