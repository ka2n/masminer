package antminer

import (
	"strings"

	"github.com/ka2n/masminer/machine"
)

// List of supported MinerType
const (
	ModelL3P    machine.Model = "Antminer L3+"
	ModelX3     machine.Model = "Antminer X3"
	ModelB3     machine.Model = "Antminer B3"
	ModelZ9Mini machine.Model = "Antminer Z9-Mini"
	ModelZ9     machine.Model = "Antminer Z9"
	ModelE3     machine.Model = "Antminer E3"
)

const (
	algoScrypt      = "scrypt"
	algoCryptonight = "cryptonight"
)

// MinerTypeFromString returns MinerType
func MinerTypeFromString(s string) (machine.Model, error) {
	switch {
	case strings.Contains(s, "X3"):
		return ModelX3, nil
	case strings.Contains(s, "L3+"):
		return ModelL3P, nil
	case strings.Contains(s, "B3"):
		return ModelB3, nil
	case strings.Contains(s, "Z9-Mini"):
		return ModelZ9Mini, nil
	case strings.Contains(s, "Z9"):
		return ModelZ9, nil
	case strings.Contains(s, "E3"):
		return ModelE3, nil
	}
	return machine.ModelUnknown, nil
}

// Algos returns list of supported algo
func Algos(m machine.Model) []string {
	switch m {
	case ModelL3P:
		return []string{
			algoScrypt,
		}
	case ModelX3:
		return []string{
			algoCryptonight,
		}
	case ModelB3:
		return []string{}
	}
	return nil
}
