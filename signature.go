package main

import (
	"crypto/sha256"
	"fmt"
)

func GenerateSignature(requestParams string) string {
	sum := sha256.Sum256([]byte(requestParams))
	sig := fmt.Sprintf("%x", sum)
	return sig
}

func ValidateSignature(requestParams, expectedSignature string) bool {
	sig := GenerateSignature(requestParams)
	return sig == expectedSignature
}
