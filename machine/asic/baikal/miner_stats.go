package baikal

type MinerStats struct {
	Devs    []SGDev
	Pools   []SGPool
	Stats   []SGStat
	Summary SGSummary
	System  MinerStatsSystem
}

type MinerStatsSystem struct {
	TempCPU string
}
