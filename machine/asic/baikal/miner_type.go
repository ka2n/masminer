package baikal

import (
	"github.com/ka2n/baikalver"
	"github.com/ka2n/masminer/inspect"
)

// List of supported MinerType
const (
	MinerTypeGP   inspect.MinerType = "Baikal Gaiant+"
	MinerTypeGX10 inspect.MinerType = "Baikal GX10"     // a.k.a. BK-X
	MinerTypeGB   inspect.MinerType = "Baikal Gaiant-B" // a.k.a. BK-B
	MinerTypeN    inspect.MinerType = "Baikal N"        // a.k.a. BK-N+
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

func minerTypeFromAPIHWV(s string) (inspect.MinerType, error) {
	m, err := baikalver.ModelFromHWV(s)
	if err != nil {
		return inspect.MinerTypeUnknown, err
	}
	switch m {
	case baikalver.GX10:
		return MinerTypeGX10, nil
	case baikalver.GB:
		return MinerTypeGB, nil
	case baikalver.GiantP:
		return MinerTypeGP, nil
	case baikalver.GN20, baikalver.GN40:
		return MinerTypeN, nil
	default:
		return inspect.MinerTypeUnknown, nil
	}
}

func minerVersionFromFWV(s string) (string, error) {
	return baikalver.VersionFromFWV(s)
}

// Algos returns list of supported algo
func Algos(m inspect.MinerType) []string {
	switch m {
	case MinerTypeGX10:
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
	case MinerTypeGP:
		return []string{
			algoX11,
			algoX13,
			algoX14,
			algoX15,
			algoQuark,
			algoQubit,
		}
	case MinerTypeGB:
		// blake256r14 blake256r8 blake2b blake2b lbry pascal
		return []string{
			algoBlake256r14,
			algoBlake256r8,
			algoBlake2b,
			algoLbry,
			algoPascal,
		}
	case MinerTypeN:
		// cryptonight cryptonight-lite
		return []string{
			algoCryptonight,
			algoCryptonightLite,
		}
	}
	return nil
}
