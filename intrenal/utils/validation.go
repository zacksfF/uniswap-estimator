package utils

import (
	"regexp"
	"strings"
)

var (
	// Ethereum address regex (0x followed by 40 hex characters)
	ethAddressRegex = regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)
)

// IsValidEthereumAddress validates Ethereum address format
func IsValidEthereumAddress(address string) bool {
	if len(address) != 42 {
		return false
	}
	return ethAddressRegex.MatchString(address)
}

// NormalizeAddress converts address to lowercase (Ethereum standard)
func NormalizeAddress(address string) string {
	return strings.ToLower(address)
}

// IsValidAmount checks if amount string represents a positive integer
func IsValidAmount(amount string) bool {
	if amount == "" || amount == "0" {
		return false
	}
	
	// Check if it's a valid big integer
	_, err := ParseBigInt(amount)
	return err == nil
}