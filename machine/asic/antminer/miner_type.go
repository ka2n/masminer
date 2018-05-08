package antminer

import (
	"strings"

	"github.com/ka2n/masminer/machine"
)

// List of supported MinerType
const (
	MinerTypeL3P machine.MinerType = "Antminer L3+"
	MinerTypeX3  machine.MinerType = "Antminer X3"
)

const (
	algoScrypt      = "scrypt"
	algoCryptonight = "cryptonight"
)

// MinerTypeFromString returns MinerType
func MinerTypeFromString(s string) (machine.MinerType, error) {
	switch {
	case strings.Contains(s, "X3"):
		return MinerTypeX3, nil
	case strings.Contains(s, "L3+"):
		return MinerTypeL3P, nil
	}
	return machine.MinerTypeUnknown, nil
}

// Algos returns list of supported algo
func Algos(m machine.MinerType) []string {
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
