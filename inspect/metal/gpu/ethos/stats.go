package ethos

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/ka2n/masminer/inspect/metal/gpu/gpustat"
)

// Stat : gpuの情報をethOSのAPIで取得します

func Stat() ([]gpustat.GPUStat, error) {
	return nil, nil
}

func GPUs() (string, error) {
	c, err := readFile("/var/run/ethos/gpucount.file")
	if err != nil {
		return "", err
	}
	return c, nil
}

func Driver() (string, error) {
	c, err := runScript(`/opt/ethos/sbin/ethos-readconf driver`)
	if err != nil {
		return "", err
	}
	return c, nil
}

func Miner() (string, error) {
	c, err := runScript(`/opt/ethos/sbin/ethos-readconf miner`)
	if err != nil {
		return "", err
	}
	return c, nil
}

func Defunct() (int64, error) {
	c, err := readFile("/var/run/ethos/defunct.file")
	if err != nil {
		return -1, err
	}
	return strconv.ParseInt(c, 10, 64)
}

func Off() (string, error) {
	c, err := runScript(`/opt/ethos/sbin/ethos-readconf off`)
	if err != nil {
		return "", err
	}
	return c, nil
}

func Allowed() (int64, error) {
	c, err := readFile("/opt/ethos/etc/allow.file")
	if err != nil {
		return -1, err
	}
	return strconv.ParseInt(c, 10, 64)
}

func Overheat() (int64, error) {
	c, err := readFile("/var/run/ethos/overheat.file")
	if err != nil {
		return -1, err
	}
	return strconv.ParseInt(c, 10, 64)
}

func PoolInfo() (string, error) {
	c, err := runScript(`cat /home/ethos/local.conf | grep -v '^#' | egrep -i 'pool|wallet|proxy'`)
	if err != nil {
		return "", err
	}
	return c, nil
}

func Pool() (string, error) {
	c, err := runScript(`/opt/ethos/sbin/ethos-readconf proxypool1`)
	if err != nil {
		return "", err
	}
	return c, nil
}

func MinerVersion() (string, error) {
	c, err := runScript(`cat /var/run/ethos/miner.versions | grep '$miner ' | cut -d" " -f2 | head -1`)
	if err != nil {
		return "", err
	}
	return c, nil
}

func BootMode() (string, error) {
	c, err := runScript(`/opt/ethos/sbin/ethos-readdata bootmode`)
	if err != nil {
		return "", err
	}
	return c, nil
}

func RackLoc() (string, error) {
	c, err := runScript(`/opt/ethos/sbin/ethos-readconf loc`)
	if err != nil {
		return "", err
	}
	return c, nil
}

func Motherboard() (string, error) {
	c, err := readFile("/var/run/ethos/motherboard.file")
	if err != nil {
		return "", err
	}
	return c, nil
}

func Rofs() (time.Duration, error) {
	c, err := readFile("/opt/ethos/etc/check-ro.file")
	if err != nil {
		return -1, err
	}
	i, err := strconv.ParseInt(c, 10, 64)
	if err != nil {
		return -1, err
	}
	t := time.Unix(i, 0)
	return time.Now().Sub(t), nil
}

func DriveName() (string, error) {
	c, err := runScript(`/opt/ethos/sbin/ethos-readdata driveinfo`)
	if err != nil {
		return "", err
	}
	return c, nil
}

func Temp() (string, error) {
	c, err := runScript(`/opt/ethos/sbin/ethos-readdata temps`)
	if err != nil {
		return "", err
	}
	return c, nil
}

func Version() (string, error) {
	c, err := readFile("/opt/ethos/etc/version")
	if err != nil {
		return "", err
	}
	return c, nil
}

func MinerSecs() (int64, error) {
	miner, err := Miner()
	if err != nil {
		return -1, err
	}
	c, err := runScript(`ps -eo pid,etime,command | grep ` + miner + ` | grep -v grep | head -1 | awk '{print $2}' |  /opt/ethos/bin/convert_time.awk`)
	if err != nil {
		return -1, err
	}
	return strconv.ParseInt(c, 10, 64)
}

func AdlError() (string, error) {
	c, err := readFile("/var/run/ethos/adl_error.file")
	if err != nil {
		return "", err
	}
	return c, nil
}

func ProxyProblem() (string, error) {
	c, err := readFile("/var/run/ethos/proxy_error.file")
	if err != nil {
		return "", err
	}
	return c, nil
}

func Updating() (string, error) {
	c, err := readFile("/var/run/ethos/updating.file")
	if err != nil {
		return "", err
	}
	return c, nil
}

func ConnectedDisplays() (string, error) {
	c, err := runScript(`/opt/ethos/sbin/ethos-readdata connecteddisplays`)
	if err != nil {
		return "", err
	}
	return c, nil
}

func Resolution() (string, error) {
	c, err := runScript(`/opt/ethos/sbin/ethos-readdata resolution`)
	if err != nil {
		return "", err
	}
	return c, nil
}

func Gethelp() (string, error) {
	c, err := runScript(`tail -1 /var/log/gethelp.log`)
	if err != nil {
		return "", err
	}
	return c, nil
}

func ConfigError() (string, error) {
	c, err := runScript(`cat /var/run/ethos/config_mode.file`)
	if err != nil {
		return "", err
	}
	return c, nil
}

func SendRemote() (string, error) {
	c, err := runScript(`cat /var/run/ethos/send_remote.file`)
	if err != nil {
		return "", err
	}
	return c, nil
}

func Autorebooted() (string, error) {
	c, err := runScript(`cat /opt/ethos/etc/autorebooted.file`)
	if err != nil {
		return "", err
	}
	return c, nil
}

func Status() (string, error) {
	c, err := runScript(`cat /var/run/ethos/status.file`)
	if err != nil {
		return "", err
	}
	return c, nil
}

func SelectedGPUs() (string, error) {
	c, err := runScript(`/opt/ethos/sbin/ethos-readconf selectedgpus`)
	if err != nil {
		return "", err
	}
	return c, nil
}

func FanRPM() (string, error) {
	c, err := runScript(`/opt/ethos/sbin/ethos-readdata fanrpm | xargs | tr -s ' '`)
	if err != nil {
		return "", err
	}
	return c, nil
}

func FanPercent() (string, error) {
	c, err := runScript(`/opt/ethos/sbin/ethos-readdata fan | xargs | tr -s ' '`)
	if err != nil {
		return "", err
	}
	return c, nil
}

func Hash() (float64, error) {
	c, err := runScript(`tail -10 /var/run/ethos/miner_hashes.file | sort -V | tail -1 | tr ' ' '\n' | awk '{sum += $1} END {print sum}'`)
	if err != nil {
		return -1, err
	}
	return strconv.ParseFloat(c, 64)
}

func MinerHashes() (string, error) {
	c, err := runScript(`tail -10 /var/run/ethos/miner_hashes.file | sort -V | tail -1`)
	if err != nil {
		return "", err
	}
	return c, nil
}

func GPUModels() (string, error) {
	c, err := readFile("/var/run/ethos/gpulist.file")
	if err != nil {
		return "", err
	}
	return c, nil
}

func Bioses() (string, error) {
	c, err := runScript(`/opt/ethos/sbin/ethos-readdata bios | xargs | tr -s ' '`)
	if err != nil {
		return "", err
	}
	return c, err
}

func DefaultCore() (string, error) {
	c, err := readFile("/var/run/ethos/defaultcore.file")
	if err != nil {
		return "", err
	}
	return c, nil
}

func DefaultMem() (string, error) {
	c, err := readFile("/var/run/ethos/defaultmem.file")
	if err != nil {
		return "", err
	}
	return c, nil
}

func VRAMSize() (string, error) {
	c, err := readFile("/var/run/ethos/vrams.file")
	if err != nil {
		return "", err
	}
	return c, nil
}

func Core() (string, error) {
	c, err := runScript(`/opt/ethos/sbin/ethos-readdata core | xargs | tr -s ' '`)
	if err != nil {
		return "", err
	}
	return c, nil
}

func Mem() (string, error) {
	c, err := runScript(`/opt/ethos/sbin/ethos-readdata mem | xargs | tr -s ' '`)
	if err != nil {
		return "", err
	}
	return c, nil
}

func MemStates() (string, error) {
	c, err := runScript(`/opt/ethos/sbin/ethos-readdata memstate | xargs | tr -s ' '`)
	if err != nil {
		return "", err
	}
	return c, nil
}

func MemInfo() (string, error) {
	c, err := readFile("/var/run/ethos/meminfo.file")
	if err != nil {
		return "", err
	}
	return c, nil
}

func Voltage() (string, error) {
	c, err := runScript(`/opt/ethos/sbin/ethos-readdata voltage | xargs | tr -s ' '`)
	if err != nil {
		return "", err
	}
	return c, nil
}

func OverheatedGPU() (string, error) {
	c, err := readFile("/var/run/ethos/overheatedgpu.file")
	if err != nil {
		return "", err
	}
	return c, nil
}

func Throttled() (string, error) {
	c, err := readFile("/var/run/ethos/throttled.file")
	if err != nil {
		return "", err
	}
	return c, nil
}

func Powertune() (string, error) {
	c, err := runScript(`/opt/ethos/sbin/ethos-readdata powertune | xargs | tr -s ' '`)
	if err != nil {
		return "", err
	}
	return c, nil
}

func runScript(cmdStr string) (string, error) {
	out, err := exec.Command("bash", "-c", cmdStr).Output()
	return strings.TrimSpace(string(out)), err
}

func readFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}
