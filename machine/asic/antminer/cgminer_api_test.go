package antminer

import (
	"reflect"
	"testing"
)

func Test_parseCGMinerVersion(t *testing.T) {
	type args struct {
		in []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"X3", args{in: []byte("STATUS=S,When=1525262182,Code=22,Msg=CGMiner versions,Description=cgminer 4.9.0|VERSION,CGMiner=4.9.0,API=3.1,Miner=1.3.0.10,CompileTime=Wed Apr 11 23:42:20 CST 2018,Type=Antminer X3|")}, "4.9.0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := parseCGMinerVersion(tt.args.in); got != tt.want {
				t.Errorf("parseCGMinerOutputToVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseHWVersionsFromCGMinerStats(t *testing.T) {
	type args struct {
		in []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"X3", args{[]byte("STATUS=S,When=1525262920,Code=70,Msg=CGMiner stats,Description=cgminer 4.9.0|CGMiner=4.9.0,Miner=1.3.0.10,CompileTime=Wed Apr 11 23:42:20 CST 2018,Type=Antminer X3|STATS=0,ID=X30,Elapsed=616151,Calls=0,Wait=0.000000,Max=0.000000,Min=99999999.000000,GHS 5s=237.58,GHS av=234.20,miner_count=3,frequency=400,fan_num=1,fan1=4200,fan2=0,temp_num=3,temp1=0,temp2=53,temp3=54,temp4=53,temp2_1=0,temp2_2=59,temp2_3=60,temp2_4=60,temp_max=56,Device Hardware%=0.0000,no_matching_work=0,chain_acn1=0,chain_acn2=60,chain_acn3=60,chain_acn4=60,chain_acs1=,chain_acs2= oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooo,chain_acs3= oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooo,chain_acs4= oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooo,chain_hw1=0,chain_hw2=0,chain_hw3=0,chain_hw4=0,chain_rate1=,chain_rate2=80.01,chain_rate3=81.49,chain_rate4=76.08|")}, []string{"1.3.0.10"}, false},
		{"S4?", args{[]byte("STATUS=S,When=1525262920,Code=70,Msg=CGMiner stats,Description=cgminer 4.9.0|CGMiner=4.9.0,Miner=1.3.0.10,CompileTime=Wed Apr 11 23:42:20 CST 2018,Type=Antminer X3|STATS=0,ID=X30,Elapsed=616151,Calls=0,Wait=0.000000,Max=0.000000,Min=99999999.000000,GHS 5s=237.58,GHS av=234.20,miner_count=3,frequency=400,hwv1=1.1,hwv2=1.2,fan_num=1,fan1=4200,fan2=0,temp_num=3,temp1=0,temp2=53,temp3=54,temp4=53,temp2_1=0,temp2_2=59,temp2_3=60,temp2_4=60,temp_max=56,Device Hardware%=0.0000,no_matching_work=0,chain_acn1=0,chain_acn2=60,chain_acn3=60,chain_acn4=60,chain_acs1=,chain_acs2= oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooo,chain_acs3= oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooo,chain_acs4= oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooo,chain_hw1=0,chain_hw2=0,chain_hw3=0,chain_hw4=0,chain_rate1=,chain_rate2=80.01,chain_rate3=81.49,chain_rate4=76.08|")}, []string{"1.1", "1.2"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseHWVersionsFromCGMinerStats(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseHWVersionsFromCGMinerStats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseHWVersionsFromCGMinerStats() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseSummaryFromCGMinerSummary(t *testing.T) {
	type args struct {
		in []byte
	}
	tests := []struct {
		name    string
		args    args
		want    MinerStatsSummary
		wantErr bool
	}{
		{"X3", args{[]byte(`STATUS=S,When=1525268733,Code=11,Msg=Summary,Description=cgminer 4.9.0|SUMMARY,Elapsed=621963,GHS 5s=226.03,GHS av=234.20,Found Blocks=10,Getworks=4473,Accepted=2462.0000,Rejected=2462.0000,Hardware Errors=0,Utility=6.23,Discarded=313941,Stale=7,Get Failures=5,Local Work=331834,Remote Failures=2,Network Blocks=4469,Total MH=145661881365.0000,Work Utility=13555343.29,Difficulty Accepted=135337607136.48925781,Difficulty Rejected=5163188222.79785156,Difficulty Stale=14680063.99658203,Best Share=37,Device Hardware%=0.0000,Device Rejected%=3.6745,Pool Rejected%=3.6745,Pool Stale%=0.0104,Last getwork=1525268732|`)},
			MinerStatsSummary{
				Elapsed:            "621963",
				GHS5s:              "226.03",
				GHSAvarage:         "234.20",
				Foundblocks:        "10",
				Getworks:           "4473",
				Accepted:           "2462.0000",
				Rejected:           "2462.0000",
				HardwareErrors:     "0",
				Utility:            "6.23",
				Discarded:          "313941",
				Stale:              "7",
				Localwork:          "331834",
				WorkUtility:        "13555343.29",
				DifficultyAccepted: "135337607136.48925781",
				DifficultyRejected: "5163188222.79785156",
				DifficultyStale:    "14680063.99658203",
				Bestshare:          "37",
			}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSummaryFromCGMinerSummary(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSummaryFromCGMinerSummary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseSummaryFromCGMinerSummary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parsePoolsFromCGMinerPools(t *testing.T) {
	type args struct {
		in []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []MinerStatsPool
		wantErr bool
	}{
		{"X3", args{[]byte(`STATUS=S,When=1525271759,Code=7,Msg=3 Pool(s),Description=cgminer 4.9.0|POOL=0,URL=stratum+tcp://stratum-xmc.antpool.com:5555,Status=Alive,Priority=0,Quota=1,Long Poll=N,Getworks=4497,Accepted=64835,Rejected=2473,Discarded=315495,Stale=7,Get Failures=5,Remote Failures=2,User=aminer.2,Last Share Time=0:00:14,Diff=2.097M,Diff1 Shares=0,Proxy Type=,Proxy=,Difficulty Accepted=135979335648.33984375,Difficulty Rejected=5186256894.79248047,Difficulty Stale=14680063.99658203,Last Share Difficulty=2097151.99951172,Has Stratum=true,Stratum Active=true,Stratum URL=stratum-xmc.antpool.com,Has GBT=false,Best Share=  0.0000,Pool Rejected%=3.6735,Pool Stale%=0.0104|POOL=1,URL=stratum+tcp://stratum-xmc.antpool.com:443%09,Status=Dead,Priority=1,Quota=1,Long Poll=N,Getworks=0,Accepted=0,Rejected=0,Discarded=0,Stale=0,Get Failures=0,Remote Failures=0,User=aminer.2,Last Share Time=0,Diff=0.000,Diff1 Shares=0,Proxy Type=,Proxy=,Difficulty Accepted=0.00000000,Difficulty Rejected=0.00000000,Difficulty Stale=0.00000000,Last Share Difficulty=0.00000000,Has Stratum=true,Stratum Active=false,Stratum URL=,Has GBT=false,Best Share=  0.0000,Pool Rejected%=0.0000,Pool Stale%=0.0000|POOL=2,URL=stratum+tcp://stratum-xmc.antpool.com:25%09,Status=Dead,Priority=2,Quota=1,Long Poll=N,Getworks=0,Accepted=0,Rejected=0,Discarded=0,Stale=0,Get Failures=0,Remote Failures=0,User=aminer.2,Last Share Time=0,Diff=0.000,Diff1 Shares=0,Proxy Type=,Proxy=,Difficulty Accepted=0.00000000,Difficulty Rejected=0.00000000,Difficulty Stale=0.00000000,Last Share Difficulty=0.00000000,Has Stratum=true,Stratum Active=false,Stratum URL=,Has GBT=false,Best Share=  0.0000,Pool Rejected%=0.0000,Pool Stale%=0.0000|`)},
			[]MinerStatsPool{
				{
					Index:               "0",
					URL:                 "stratum+tcp://stratum-xmc.antpool.com:5555",
					User:                "aminer.2",
					Status:              "Alive",
					StratumActive:       "true",
					Priority:            "0",
					Getworks:            "4497",
					Accepted:            "64835",
					Rejected:            "2473",
					Discarded:           "315495",
					Stale:               "7",
					Diff:                "2.097M",
					Diff1Shares:         "0",
					DifficultyAccepted:  "135979335648.33984375",
					DifficultyRejected:  "5186256894.79248047",
					DifficultyStale:     "14680063.99658203",
					LastShareDifficulty: "2097151.99951172",
					LastShareTime:       "0:00:14",
				},
				{
					Index:               "1",
					URL:                 "stratum+tcp://stratum-xmc.antpool.com:443%09",
					User:                "aminer.2",
					Status:              "Dead",
					StratumActive:       "false",
					Priority:            "1",
					Getworks:            "0",
					Accepted:            "0",
					Rejected:            "0",
					Discarded:           "0",
					Stale:               "0",
					Diff:                "0.000",
					Diff1Shares:         "0",
					DifficultyAccepted:  "0.00000000",
					DifficultyRejected:  "0.00000000",
					DifficultyStale:     "0.00000000",
					LastShareDifficulty: "0.00000000",
					LastShareTime:       "0",
				},
				{
					Index:               "2",
					URL:                 "stratum+tcp://stratum-xmc.antpool.com:25%09",
					User:                "aminer.2",
					Status:              "Dead",
					StratumActive:       "false",
					Priority:            "2",
					Getworks:            "0",
					Accepted:            "0",
					Rejected:            "0",
					Discarded:           "0",
					Stale:               "0",
					Diff:                "0.000",
					Diff1Shares:         "0",
					DifficultyAccepted:  "0.00000000",
					DifficultyRejected:  "0.00000000",
					DifficultyStale:     "0.00000000",
					LastShareDifficulty: "0.00000000",
					LastShareTime:       "0",
				},
			}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePoolsFromCGMinerPools(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePoolsFromCGMinerPools() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parsePoolsFromCGMinerPools() = \n%v\n%v", got, tt.want)
			}
		})
	}
}

func Test_parseDevsFromCGMinerStats(t *testing.T) {
	type args struct {
		in []byte
	}
	tests := []struct {
		name    string
		args    args
		want    MinerStatsDevs
		wantErr bool
	}{
		{"X3", args{[]byte(`STATUS=S,When=1525269011,Code=70,Msg=CGMiner stats,Description=cgminer 4.9.0|CGMiner=4.9.0,Miner=1.3.0.10,CompileTime=Wed Apr 11 23:42:20 CST 2018,Type=Antminer X3|STATS=0,ID=X30,Elapsed=622242,Calls=0,Wait=0.000000,Max=0.000000,Min=99999999.000000,GHS 5s=228.65,GHS av=234.20,miner_count=3,frequency=400,fan_num=1,fan1=4170,fan2=0,temp_num=3,temp1=0,temp2=52,temp3=53,temp4=52,temp2_1=0,temp2_2=58,temp2_3=59,temp2_4=59,temp_max=55,Device Hardware%=0.0000,no_matching_work=0,chain_acn1=0,chain_acn2=60,chain_acn3=60,chain_acn4=60,chain_acs1=,chain_acs2= oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooo,chain_acs3= oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooo,chain_acs4= oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooo,chain_hw1=0,chain_hw2=0,chain_hw3=0,chain_hw4=0,chain_rate1=,chain_rate2=75.83,chain_rate3=80.75,chain_rate4=72.07|`)},
			MinerStatsDevs{
				Fans: []string{"4170", "0"},
				Chains: []MinerStatsChain{
					{
						Index:    "1",
						TempPCB:  "0",
						TempChip: "0",
						Acn:      "0",
						Freq:     "400",
						Rate:     "",
						Hw:       "0",
						Status:   "",
					},
					{
						Index:    "2",
						TempPCB:  "52",
						TempChip: "58",
						Acn:      "60",
						Freq:     "400",
						Rate:     "75.83",
						Hw:       "0",
						Status:   " oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooo",
					},
					{
						Index:    "3",
						TempPCB:  "53",
						TempChip: "59",
						Acn:      "60",
						Freq:     "400",
						Rate:     "80.75",
						Hw:       "0",
						Status:   " oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooo",
					},
					{
						Index:    "4",
						TempPCB:  "52",
						TempChip: "59",
						Acn:      "60",
						Freq:     "400",
						Rate:     "72.07",
						Hw:       "0",
						Status:   " oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooo",
					},
				},
			}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDevsFromCGMinerStats(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDevsFromCGMinerStats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseDevsFromCGMinerStats() = %v, want %v", got, tt.want)
			}
		})
	}
}
