package baikal

import "encoding/json"

type SGMultipleCMDResponse struct {
	ID int `json:"id"`
}

type SGAPIStatus struct {
	Status      string `json:"STATUS"`
	When        int    `json:"When"`
	Code        int    `json:"Code"`
	Msg         string `json:"Msg"`
	Description string `json:"Description"`
}

type SGCommonResponse struct {
	Status []SGAPIStatus `json:"STATUS"`
	ID     int           `json:"id"`
}

type SGSummaryResponse struct {
	SGCommonResponse
	Summary []SGSummary `json:"SUMMARY"`
}

type SGStatsResponse struct {
	SGCommonResponse
	Stats []SGStat `json:"STATS"`
}

type SGVersionResponse struct {
	SGCommonResponse
	Version []SGVersion `json:"VERSION"`
}

type SGPoolsResponse struct {
	SGCommonResponse
	Pools []SGPool `json:"POOLS"`
}

type SGDevsResponse struct {
	SGCommonResponse
	Devs []SGDev `json:"DEVS"`
}

type SGSummary struct {
	Elapsed            int     `json:"Elapsed"`
	MHSAv              float64 `json:"MHS av"`
	MHS5S              float64 `json:"MHS 5s"`
	KHSAv              float64 `json:"KHS av"`
	KHS5S              float64 `json:"KHS 5s"`
	FoundBlocks        int     `json:"Found Blocks"`
	Getworks           int     `json:"Getworks"`
	Accepted           int     `json:"Accepted"`
	Rejected           int     `json:"Rejected"`
	HardwareErrors     int     `json:"Hardware Errors"`
	Utility            float64 `json:"Utility"`
	Discarded          int     `json:"Discarded"`
	Stale              int     `json:"Stale"`
	GetFailures        int     `json:"Get Failures"`
	LocalWork          int     `json:"Local Work"`
	RemoteFailures     int     `json:"Remote Failures"`
	NetworkBlocks      int     `json:"Network Blocks"`
	TotalMH            float64 `json:"Total MH"`
	WorkUtility        float64 `json:"Work Utility"`
	DifficultyAccepted float64 `json:"Difficulty Accepted"`
	DifficultyRejected float64 `json:"Difficulty Rejected"`
	DifficultyStale    float64 `json:"Difficulty Stale"`
	BestShare          float64 `json:"Best Share"`
	DeviceHardware     float64 `json:"Device Hardware%"`
	DeviceRejected     float64 `json:"Device Rejected%"`
	PoolRejected       float64 `json:"Pool Rejected%"`
	PoolStale          float64 `json:"Pool Stale%"`
	LastGetwork        int     `json:"Last getwork"`
}

type SGVersion struct {
	Miner   string `json:"Miner"`
	SGMiner string `json:"SGMiner"`
	API     string `json:"API"`
}

type SGPool struct {
	POOL                int     `json:"POOL"`
	Name                string  `json:"Name"`
	URL                 string  `json:"URL"`
	Profile             string  `json:"Profile"`
	Algorithm           string  `json:"Algorithm"`
	AlgorithmType       string  `json:"Algorithm Type"`
	Description         string  `json:"Description"`
	Status              string  `json:"Status"`
	Priority            int     `json:"Priority"`
	Quota               int     `json:"Quota"`
	LongPoll            string  `json:"Long Poll"`
	Getworks            int     `json:"Getworks"`
	Accepted            int     `json:"Accepted"`
	Rejected            int     `json:"Rejected"`
	Works               int     `json:"Works"`
	Discarded           int     `json:"Discarded"`
	Stale               int     `json:"Stale"`
	GetFailures         int     `json:"Get Failures"`
	RemoteFailures      int     `json:"Remote Failures"`
	User                string  `json:"User"`
	LastShareTime       int     `json:"Last Share Time"`
	Diff1Shares         float64 `json:"Diff1 Shares"`
	ProxyType           string  `json:"Proxy Type"`
	Proxy               string  `json:"Proxy"`
	DifficultyAccepted  float64 `json:"Difficulty Accepted"`
	DifficultyRejected  float64 `json:"Difficulty Rejected"`
	DifficultyStale     float64 `json:"Difficulty Stale"`
	LastShareDifficulty float64 `json:"Last Share Difficulty"`
	HasStratum          bool    `json:"Has Stratum"`
	StratumActive       bool    `json:"Stratum Active"`
	StratumURL          string  `json:"Stratum URL"`
	HasGBT              bool    `json:"Has GBT"`
	BestShare           float64 `json:"Best Share"`
	PoolRejected        float64 `json:"Pool Rejected%"`
	PoolStale           float64 `json:"Pool Stale%"`
}

type SGStat struct {
	STATS     int         `json:"STATS"`
	ID        string      `json:"ID"`
	Elapsed   int         `json:"Elapsed"`
	Calls     int         `json:"Calls"`
	Wait      float64     `json:"Wait"`
	Max       float64     `json:"Max"`
	Min       float64     `json:"Min"`
	ChipCount int         `json:"Chip Count"`
	Clock     int         `json:"Clock"`
	HWV       json.Number `json:"HWV"`
	FWV       json.Number `json:"FWV"`
	Algo      string      `json:"Algo"`
	USBPipe   string      `json:"USB Pipe"`
	USBDelay  string      `json:"USB Delay"`
	USBTmo    string      `json:"USB tmo"`
}

type SGDev struct {
	ASC                 int     `json:"ASC"`
	Name                string  `json:"Name"`
	ID                  int     `json:"ID"`
	Enabled             string  `json:"Enabled"`
	Status              string  `json:"Status"`
	Temperature         float64 `json:"Temperature"`
	MHSAv               float64 `json:"MHS av"`
	MHS5S               float64 `json:"MHS 5s"`
	Accepted            int     `json:"Accepted"`
	Rejected            int     `json:"Rejected"`
	HardwareErrors      int     `json:"Hardware Errors"`
	Utility             float64 `json:"Utility"`
	LastSharePool       int     `json:"Last Share Pool"`
	LastShareTime       int     `json:"Last Share Time"`
	TotalMH             float64 `json:"Total MH"`
	Diff1Work           float64 `json:"Diff1 Work"`
	DifficultyAccepted  float64 `json:"Difficulty Accepted"`
	DifficultyRejected  float64 `json:"Difficulty Rejected"`
	LastShareDifficulty float64 `json:"Last Share Difficulty"`
	NoDevice            bool    `json:"No Device"`
	LastValidWork       int     `json:"Last Valid Work"`
	DeviceHardware      float64 `json:"Device Hardware%"`
	DeviceRejected      float64 `json:"Device Rejected%"`
	DeviceElapsed       int     `json:"Device Elapsed"`
}
