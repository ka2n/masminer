package antminer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ka2n/masminer/machine"
	"golang.org/x/crypto/ssh"
)

// GetMinerSetting returns MinerSetting or create with default setting
func (c *Client) GetMinerSetting() (setting MinerSetting, err error) {
	return c.GetMinerSettingContext(nil)
}

// GetMinerSettingContext returns MinerSetting or create with default setting
func (c *Client) GetMinerSettingContext(ctx context.Context) (setting MinerSetting, err error) {
	info, err := c.GetSystemInfo()
	if err != nil {
		return setting, err
	}
	return getMinerSetting(ctx, c.ssh, info.Model)
}

func (c *Client) WriteCGMinerSettingContext(ctx context.Context, setting MinerSetting) error {
	buf, err := json.Marshal(setting)
	if err != nil {
		return err
	}
	err = runRemoteShell(ctx, c.ssh, fmt.Sprintf(`
set -ex
CONFIG_PATH=%s
cat <<'EOF' > $CONFIG_PATH
%s
EOF
`, minerConfigPath, string(buf)))
	if err != nil {
		return err
	}

	_, err = outputMinerRPC(ctx, c.ssh, "restart", "")
	return err
}

func (c *Client) WriteCGMinerSetting(setting MinerSetting) error {
	return c.WriteCGMinerSettingContext(nil, setting)
}

func getMinerSetting(ctx context.Context, client *ssh.Client, minerType machine.Model) (setting MinerSetting, err error) {
	defaultConf := defaultSetting(minerType)
	dconfb, err := json.Marshal(defaultConf)
	if err != nil {
		panic(err)
	}

	output, err := outputRemoteShell(ctx, client, fmt.Sprintf(`
set -ex
CONFIG_PATH=%s
if [ ! -f $CONFIG_PATH ] ; then
cat <<'EOF' > $CONFIG_PATH
%s
EOF
fi
cat $CONFIG_PATH
`, minerConfigPath, string(dconfb)))
	if err != nil {
		return setting, err
	}

	if err := json.Unmarshal(output, &setting); err != nil {
		return setting, err
	}
	return
}

func defaultSetting(mt machine.Model) (setting MinerSetting) {
	opt := make(map[string]string)
	opt["api-allow"] = defaultAPIAllow
	opt["api-groups"] = defaultAPIGroups
	opt["api-listen"] = defaultAPIListen
	opt["api-network"] = defaultAPINetwork
	opt["bitmain-use-vil"] = defaultBitmainUseVil
	setting.Options = opt

	switch mt {
	case ModelX3:
		setting.Pools = []PoolSetting{
			{"stratum+tcp://stratum-xmc.antpool.com:5555", "aminerr.1", "x"},
			{"stratum+tcp://stratum-xmc.antpool.com:443", "aminerr.1", "x"},
			{"stratum+tcp://stratum-xmc.antpool.com:25", "aminerr.1", "x"},
		}
		return
	case ModelL3P:
		setting.Pools = []PoolSetting{
			{"stratum+tcp://scrypt.jp.nicehash.com:3333", "383qTNfyvT3cVaiVWZPLsyehR9dW7fRpPK", "x"},
			{"stratum+tcp://scrypt.eu.nicehash.com:3333", "383qTNfyvT3cVaiVWZPLsyehR9dW7fRpPK", "x"},
			{"stratum+tcp://scrypt.usa.nicehash.com:3333", "383qTNfyvT3cVaiVWZPLsyehR9dW7fRpPK", "x"},
		}
		return
	}
	return
}
