package gpu

import (
	"github.com/ka2n/masminer/machine/metal/gpu/ethos"
	"github.com/ka2n/masminer/machine/metal/gpu/gpustat"
)

// Stat : GPUの情報をベンダーごとに集めます
func Stat() ([]gpustat.GPUStat, error) {
	return ethos.Stat()
}
