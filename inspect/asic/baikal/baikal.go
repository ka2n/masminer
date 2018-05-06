package baikal

const (
	manufactureName = "baikal"
	minerStopCMD    = "sudo /opt/scripta/startup/miner-stop.sh"
	minerStartCMD   = "sudo /opt/scripta/startup/miner-start.sh"

	minerConfPath    = "/opt/scripta/etc/miner.conf"
	minerOptionsPath = "/opt/scripta/etc/miner.options.json"
	minerPoolsPath   = "/opt/scripta/etc/miner.pools.json"
)

const (
	defaultAPIPort  = "4028"
	defaultAPIAllow = "W:127.0.0.1,W:192.168.0.0/16"
)
