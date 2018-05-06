package asic

import (
	"context"

	"github.com/ka2n/masminer/inspect"

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
	RigInfo(context.Context) (inspect.RigInfo, error)
}

// StatReader : ASIC status
type StatReader interface {
	RigStat(context.Context) (inspect.RigStat, error)
}

// SettingReader : get setting
type SettingReader interface {
	MinerSetting(context.Context) (inspect.MinerSetting, error)
}

// SettingWriter updates setting
type SettingWriter interface {
	SetMinerSetting(context.Context, inspect.MinerSetting) error
}
