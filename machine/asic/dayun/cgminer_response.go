package dayun

import (
	"encoding/json"
	"unicode"

	"github.com/ka2n/masminer/minerapi"
)

type SummaryResponse struct {
	minerapi.ResponseCommon
	Summary []Summary `json:"SUMMARY"`
}

type StatsResponse struct {
	minerapi.ResponseCommon
	Stats []Stat `json:"STATS"`
}

type VersionResponse struct {
	minerapi.ResponseCommon
	Version []Version `json:"VERSION"`
}

type PoolsResponse struct {
	minerapi.ResponseCommon
	Pools []Pool `json:"POOLS"`
}

type DevsResponse struct {
	minerapi.ResponseCommon
	Devs []Dev `json:"DEVS"`
}

type Summary struct {
	Elapsed            int     `json:"Elapsed"`
	MHSAv              float64 `json:"MHS av"`
	MHS5S              float64 `json:"MHS 5s"`
	MHS1M              float64 `json:"MHS 1m"`
	MHS5M              float64 `json:"MHS 5m"`
	MHS15M             float64 `json:"MHS 15m"`
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

type Stat struct {
	STATS                    int     `json:"STATS"`
	ID                       string  `json:"ID"`
	Elapsed                  int     `json:"Elapsed"`
	Calls                    int     `json:"Calls"`
	Wait                     float64 `json:"Wait"`
	Max                      float64 `json:"Max"`
	Min                      float64 `json:"Min"`
	MHS30S                   float64 `json:"MHS 30S,omitempty"`
	MHS5M                    float64 `json:"MHS 5m,omitempty"`
	FanNunber                int     `json:"Fan Nunber,omitempty"`
	FanIn                    int     `json:"Fan In,omitempty"`
	FanOut                   int     `json:"Fan Out,omitempty"`
	Frequency                int     `json:"Frequency,omitempty"`
	TemperatureCore          int     `json:"Temperature Core,omitempty"`
	AutoFrequency            bool    `json:"AutoFrequency,omitempty"`
	CheckNonceSubmit         bool    `json:"CheckNonceSubmit,omitempty"`
	AutoStopMinerEnable      bool    `json:"AutoStopMiner_Enable,omitempty"`
	AutoStopMinerTemperature int     `json:"AutoStopMiner_Temperature,omitempty"`
	AutoStopMinerFanIn       int     `json:"AutoStopMiner_Fan_In,omitempty"`
	AutoStopMinerFanOut      int     `json:"AutoStopMiner_Fan_Out,omitempty"`
	AutoRestartEnable        bool    `json:"AutoRestart_Enable,omitempty"`
	AutoRestartHashrate      int     `json:"AutoRestart_Hashrate,omitempty"`
	AutoRestartFailedRate    int     `json:"AutoRestart_FailedRate,omitempty"`
	CH                       []DevCh
	PoolCalls                int  `json:"Pool Calls,omitempty"`
	PoolAttempts             int  `json:"Pool Attempts,omitempty"`
	PoolWait                 int  `json:"Pool Wait,omitempty"`
	PoolMax                  int  `json:"Pool Max,omitempty"`
	PoolMin                  int  `json:"Pool Min,omitempty"`
	PoolAv                   int  `json:"Pool Av,omitempty"`
	WorkHadRollTime          bool `json:"Work Had Roll Time,omitempty"`
	WorkCanRoll              bool `json:"Work Can Roll,omitempty"`
	WorkHadExpire            bool `json:"Work Had Expire,omitempty"`
	WorkRollTime             int  `json:"Work Roll Time,omitempty"`
	WorkDiff                 int  `json:"Work Diff,omitempty"`
	MinDiff                  int  `json:"Min Diff,omitempty"`
	MaxDiff                  int  `json:"Max Diff,omitempty"`
	MinDiffCount             int  `json:"Min Diff Count,omitempty"`
	MaxDiffCount             int  `json:"Max Diff Count,omitempty"`
	TimesSent                int  `json:"Times Sent,omitempty"`
	BytesSent                int  `json:"Bytes Sent,omitempty"`
	TimesRecv                int  `json:"Times Recv,omitempty"`
	BytesRecv                int  `json:"Bytes Recv,omitempty"`
	NetBytesSent             int  `json:"Net Bytes Sent,omitempty"`
	NetBytesRecv             int  `json:"Net Bytes Recv,omitempty"`
}

type Version struct {
	CGMiner string `json:"CGMiner"`
	API     string `json:"API"`
}

type Pool struct {
	POOL                int     `json:"POOL"`
	URL                 string  `json:"URL"`
	USER                string  `json:"USER"`
	PASSWORD            string  `json:"PASSWORD"`
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
	Diff1Shares         int     `json:"Diff1 Shares"`
	ProxyType           string  `json:"Proxy Type"`
	Proxy               string  `json:"Proxy"`
	DifficultyAccepted  int     `json:"Difficulty Accepted"`
	DifficultyRejected  int     `json:"Difficulty Rejected"`
	DifficultyStale     int     `json:"Difficulty Stale"`
	LastShareDifficulty int     `json:"Last Share Difficulty"`
	HasStratum          bool    `json:"Has Stratum"`
	StratumActive       bool    `json:"Stratum Active"`
	StratumURL          string  `json:"Stratum URL"`
	HasGBT              bool    `json:"Has GBT"`
	BestShare           int     `json:"Best Share"`
	PoolRejected        float64 `json:"Pool Rejected%"`
	PoolStale           float64 `json:"Pool Stale%"`
}

type devBase struct {
	ASC                      int     `json:"ASC"`
	Name                     string  `json:"Name"`
	ID                       int     `json:"ID"`
	Enabled                  string  `json:"Enabled"`
	Status                   string  `json:"Status"`
	Accepted                 int     `json:"Accepted"`
	Rejected                 int     `json:"Rejected"`
	HardwareErrors           int     `json:"Hardware Errors"`
	Utility                  float64 `json:"Utility"`
	LastSharePool            int     `json:"Last Share Pool"`
	LastShareTime            int     `json:"Last Share Time"`
	TotalMH                  float64 `json:"Total MH"`
	Diff1Work                int     `json:"Diff1 Work"`
	DifficultyAccepted       int     `json:"Difficulty Accepted"`
	DifficultyRejected       int     `json:"Difficulty Rejected"`
	LastShareDifficulty      int     `json:"Last Share Difficulty"`
	LastValidWork            int     `json:"Last Valid Work"`
	DeviceHardware           float64 `json:"Device Hardware%"`
	DeviceRejected           float64 `json:"Device Rejected%"`
	DeviceElapsed            int     `json:"Device Elapsed"`
	MHS30S                   float64 `json:"MHS 30S"`
	MHS5M                    float64 `json:"MHS 5m"`
	FanNunber                int     `json:"Fan Nunber"`
	FanIn                    int     `json:"Fan In"`
	FanOut                   int     `json:"Fan Out"`
	Frequency                int     `json:"Frequency"`
	TemperatureCore          int     `json:"Temperature Core"`
	AutoFrequency            bool    `json:"AutoFrequency"`
	CheckNonceSubmit         bool    `json:"CheckNonceSubmit"`
	AutoStopMinerEnable      bool    `json:"AutoStopMiner_Enable"`
	AutoStopMinerTemperature int     `json:"AutoStopMiner_Temperature"`
	AutoStopMinerFanIn       int     `json:"AutoStopMiner_Fan_In"`
	AutoStopMinerFanOut      int     `json:"AutoStopMiner_Fan_Out"`
	AutoRestartEnable        bool    `json:"AutoRestart_Enable"`
	AutoRestartHashrate      int     `json:"AutoRestart_Hashrate"`
	AutoRestartFailedRate    int     `json:"AutoRestart_FailedRate"`
}

type Dev struct {
	devBase
	CH []DevCh
}

func (dev *Dev) UnmarshalJSON(data []byte) error {
	var b devBase
	if err := json.Unmarshal(data, &b); err != nil {
		return nil
	}
	dev.devBase = b

	type rtype map[string]json.RawMessage
	var rv rtype
	if err := json.Unmarshal(data, &rv); err != nil {
		return err
	}

	for k, v := range rv {
		if len(k) >= 3 && k[0] == 'C' && k[1] == 'H' && unicode.IsDigit(rune(k[2])) {
			var ch DevCh
			if err := json.Unmarshal([]byte(v), &ch); err != nil {
				return err
			}
			dev.CH = append(dev.CH, ch)
		}
	}

	return nil
}

type DevCh struct {
	Temperature int           `json:"Temperature"`
	MHS30S      int           `json:"MHS 30S"`
	Status      []DevChStatus `json:"status"`
}

type DevChStatus struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Accept int `json:"accept"`
	Failed int `json:"failed"`
	Reject int `json:"reject"`
}
