package main

import (
	"testing"
)

func TestGenerateSignature(t *testing.T) {
	type args struct {
		requestParams string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "valid request",
			args: args{requestParams: "BUGBOUNTY231QvE8dZshpKhaOmHY173252254717155146485f4a6fcf-9048-4a0b-afc2-ed92d60fb1bf"},
			want: "c950a78220e1fc9bfc2514729ef9b74ad734feec260a3ed374ae6c063d2be40b",
		},
		{
			name: "empty string",
			args: args{requestParams: ""},
			want: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name: "return signature",
			args: args{requestParams: "DECLINED"+"orderID"+"merchantOrderID" + "5f4a6fcf-9048-4a0b-afc2-ed92d60fb1bf"},
			want: "64128a36497a30bbd6eda8e5cf026acbf99d6f40c28f1954b3a7bcfcc64ce53e",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateSignature(tt.args.requestParams); got != tt.want {
				t.Errorf("GenerateSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateSignature(t *testing.T) {
	type args struct {
		requestParams     string
		expectedSignature string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid signature",
			args: args{
				requestParams:     "BUGBOUNTY231QvE8dZshpKhaOmHY173252254717155146485f4a6fcf-9048-4a0b-afc2-ed92d60fb1bf",
				expectedSignature: "c950a78220e1fc9bfc2514729ef9b74ad734feec260a3ed374ae6c063d2be40b",
			},
			want: true,
		},
		{
			name: "invalid signature",
			args: args{
				requestParams:     "BUGBOUNTY231QvE8dZshpKhaOmHY173252254717155146485f4a6fcf-9048-4a0b-afc2-ed92d60fb1bf",
				expectedSignature: "c950a78220e1fc9bfc2514729ef9b74ad734feec260a3ee40b",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateSignature(tt.args.requestParams, tt.args.expectedSignature); got != tt.want {
				t.Errorf("validateSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}
