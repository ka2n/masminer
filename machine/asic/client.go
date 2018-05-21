package asic

import (
	"context"

	"github.com/ka2n/masminer/machine"

	"golang.org/x/crypto/ssh"
)

// Client interface
type Client interface {
	Connector
	Controller
	InfoReader
	StatReader
	SettingReader
	SettingWriter
	Setup() error
}

// Connector handles *ssh.Client
type Connector interface {
	SetSSH(ssh *ssh.Client)
	Close() error
}

// Controller controlls ASIC
type Controller interface {
	MineStop(context.Context) error
	MineStart(context.Context) error
	Restart(context.Context) error
	Reboot(context.Context) error
}

// InfoReader : ASIC basic info
type InfoReader interface {
	RigInfo(context.Context) (machine.RigInfo, error)
}

// StatReader : ASIC status
type StatReader interface {
	RigStat(context.Context) (machine.RigStat, error)
}

// SettingReader : get setting
type SettingReader interface {
	MinerSetting(context.Context) (machine.MinerSetting, error)
}

// SettingWriter updates setting
type SettingWriter interface {
	SetMinerSetting(context.Context, machine.MinerSetting) error
}
