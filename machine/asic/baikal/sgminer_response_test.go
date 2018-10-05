package baikal

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/ka2n/masminer/minerapi"
)

func TestParseMinerResponse(t *testing.T) {
	tf, err := os.Open("../testdata/baikal/x10_summary+version+stats+pools.json")
	if err != nil {
		t.Fatal(err)
	}
	defer tf.Close()

	var resp struct {
		minerapi.MultipleResponse
		Summary []SGSummaryResponse `json:"summary"`
		Version []SGVersionResponse `json:"version"`
		Stats   []SGStatsResponse   `json:"stats"`
		Pools   []SGPoolsResponse   `json:"pools"`
	}

	decoder := json.NewDecoder(tf)
	if err := decoder.Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if len(resp.Summary) != 1 {
		t.Errorf("expect Summary have 1 result, got: %v", resp.Summary)
	}
	if len(resp.Summary[0].Summary) != 1 {
		t.Errorf("expect Summary[0].Summary have 1 result, got: %v", len(resp.Summary[0].Summary))
	}
	if len(resp.Version) != 1 {
		t.Errorf("expect Summary have 1 result, got: %v", resp.Version)
	}
	if len(resp.Version[0].Version) != 1 {
		t.Errorf("expect Version[0].Version have 1 result, got: %v", len(resp.Version[0].Version))
	}
	if len(resp.Stats) != 1 {
		t.Errorf("expect Summary have 1 result, got: %v", resp.Stats)
	}
	if len(resp.Stats[0].Stats) != 3 {
		t.Errorf("expect Stats[0].Stats have 3 result, got: %v", len(resp.Stats[0].Stats))
	}
	if len(resp.Pools) != 1 {
		t.Errorf("expect Pools have 1 result, got: %v", resp.Pools)
	}
	if len(resp.Pools[0].Pools) != 15 {
		t.Errorf("expect Pools[0].Pools have 15 result, got: %v", len(resp.Pools[0].Pools))
	}

	t.Log(resp)
}
