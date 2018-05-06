package antminer

import (
	"strings"

	"github.com/ka2n/masminer/inspect"
)

// List of supported MinerType
const (
	MinerTypeL3P inspect.MinerType = "Antminer L3+"
	MinerTypeX3  inspect.MinerType = "Antminer X3"
)

const (
	algoScrypt      = "scrypt"
	algoCryptonight = "cryptonight"
)

// MinerTypeFromString returns MinerType
func MinerTypeFromString(s string) (inspect.MinerType, error) {
	switch {
	case strings.Contains(s, "X3"):
		return MinerTypeX3, nil
	case strings.Contains(s, "L3+"):
		return MinerTypeL3P, nil
	}
	return inspect.MinerTypeUnknown, nil
}

// Algos returns list of supported algo
func Algos(m inspect.MinerType) []string {
	switch m {
	case MinerTypeL3P:
		return []string{
			algoScrypt,
		}
	case MinerTypeX3:
		return []string{
			algoCryptonight,
		}
	}
	return nil
}
