package baikal

import (
	"github.com/ka2n/baikalver"
	"github.com/ka2n/masminer/machine"
)

// List of supported MinerType
const (
	ModelGP   machine.Model = "Baikal Gaiant+"
	ModelGX10 machine.Model = "Baikal GX10"     // a.k.a. BK-X
	ModelGB   machine.Model = "Baikal Gaiant-B" // a.k.a. BK-B
	ModelN    machine.Model = "Baikal N"        // a.k.a. BK-N+
)

const (
	algoX11             = "x11"
	algoX13             = "x13"
	algoX14             = "x14"
	algoX15             = "x15"
	algoQuark           = "quark"
	algoQubit           = "qubit"
	algoMyriadGroestl   = "myr-gr"
	algoSkein           = "skein"
	algoX11Gost         = "x11gost"
	algoNist5           = "nist5"
	algoBlake256r14     = "blake256r14"
	algoBlake256r8      = "blake256r8"
	algoBlake2b         = "blake2b"
	algoLbry            = "lbry"
	algoPascal          = "pascal"
	algoCryptonight     = "cryptonight"
	algoCryptonightLite = "cryptonight-lite"
)

func modelFromAPIHWV(s string) (machine.Model, error) {
	m, err := baikalver.ModelFromHWV(s)
	if err != nil {
		return machine.ModelUnknown, err
	}
	switch m {
	case baikalver.GX10:
		return ModelGX10, nil
	case baikalver.GB:
		return ModelGB, nil
	case baikalver.GiantP:
		return ModelGP, nil
	case baikalver.GN20, baikalver.GN40:
		return ModelN, nil
	default:
		return machine.ModelUnknown, nil
	}
}

func minerVersionFromFWV(s string) (string, error) {
	return baikalver.VersionFromFWV(s)
}

// Algos returns list of supported algo
func Algos(m machine.Model) []string {
	switch m {
	case ModelGX10:
		// X11 Quark Qubit Myriad-Groestl Skein X11Gost Nist5
		return []string{
			algoX11,
			algoQuark,
			algoQuark,
			algoMyriadGroestl,
			algoSkein,
			algoX11Gost,
			algoNist5,
		}
	case ModelGP:
		return []string{
			algoX11,
			algoX13,
			algoX14,
			algoX15,
			algoQuark,
			algoQubit,
		}
	case ModelGB:
		// blake256r14 blake256r8 blake2b blake2b lbry pascal
		return []string{
			algoBlake256r14,
			algoBlake256r8,
			algoBlake2b,
			algoLbry,
			algoPascal,
		}
	case ModelN:
		// cryptonight cryptonight-lite
		return []string{
			algoCryptonight,
			algoCryptonightLite,
		}
	}
	return nil
}
