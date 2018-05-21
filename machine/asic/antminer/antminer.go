package antminer

const (
	manufactureName = "antminer"
	minerConfigPath = "/config/cgminer.conf"
	metadataPath    = "/usr/bin/compile_time"

	minerAPISummaryCMD = "cgminer-api -o"
	minerAPIPoolsCMD   = "cgminer-api -o pools"
	minerAPIStatsCMD   = "cgminer-api -o stats"
	minerAPIVersionCMD = "cgminer-api -o version"

	minerInitdCMD = "/etc/init.d/cgminer.sh %s >/dev/null 2>&1"

	minerBMMinerConfigPath    = "/config/bmminer.conf"
	minerBMMinerAPISummaryCMD = "bmminer-api -o"
	minerBMMinerAPIPoolsCMD   = "bmminer-api -o pools"
	minerBMMinerAPIStatsCMD   = "bmminer-api -o stats"
	minerBMMinerAPIVersionCMD = "bmminer-api -o version"
	minerBMMinerInitdCMD      = "/etc/init.d/bmminer.sh %s >/dev/null 2>&1"
)

const (
	defaultAPIListen     = "true"
	defaultAPINetwork    = "true"
	defaultAPIAllow      = "A:0/0,W:*"
	defaultAPIGroups     = "A:stats:pools:devs:summary:version"
	defaultBitmainUseVil = "true"
)
