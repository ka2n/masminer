package gpu

import (
	"github.com/ka2n/masminer/inspect/metal/gpu/ethos"
	"github.com/ka2n/masminer/inspect/metal/gpu/gpustat"
)

// Stat : GPUの情報をベンダーごとに集めます
func Stat() ([]gpustat.GPUStat, error) {
	return ethos.Stat()
}
