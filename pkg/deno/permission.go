package deno

type PermissionName string

const (
	PermissionNameRun    PermissionName = "run"
	PermissionNameRead   PermissionName = "read"
	PermissionNameWrite  PermissionName = "write"
	PermissionNameNet    PermissionName = "net"
	PermissionNameEnv    PermissionName = "env"
	PermissionNameSys    PermissionName = "sys"
	PermissionNameFfi    PermissionName = "ffi"
	PermissionNameHrtime PermissionName = "hrtime"
)

type SysKind string

const (
	SysKindAll               SysKind = ""
	SysKindLoadavg           SysKind = "loadavg"
	SysKindHostname          SysKind = "hostname"
	SysKindSystemMemoryInfo  SysKind = "systemMemoryInfo"
	SysKindNetworkInterfaces SysKind = "networkInterfaces"
	SysKindosRelease         SysKind = "osRelease"
	SysKindosUid             SysKind = "uid"
	SysKindosGid             SysKind = "gid"
)

func ParseSysKind(kind string) (SysKind, bool) {
	switch SysKind(kind) {
	case SysKindAll:
		return SysKindAll, true
	case SysKindLoadavg:
		return SysKindLoadavg, true
	case SysKindHostname:
		return SysKindHostname, true
	case SysKindSystemMemoryInfo:
		return SysKindSystemMemoryInfo, true
	case SysKindNetworkInterfaces:
		return SysKindNetworkInterfaces, true
	case SysKindosRelease:
		return SysKindosRelease, true
	case SysKindosUid:
		return SysKindosUid, true
	case SysKindosGid:
		return SysKindosGid, true
	default:
		return "", false
	}
}

type PermissionDescriptor struct {
	Name PermissionName `json:"name"`
	// run
	Command string `json:"command,omitempty"`
	// read, write, ffi
	Path string `json:"path,omitempty"`
	// net
	Host *HostPort `json:"host,omitempty"`
	// env
	Variable string `json:"variable,omitempty"`
	// sys
	Kind SysKind `json:"kind,omitempty"`
}

func ParsePermissionDescriptor(name string, target string) (*PermissionDescriptor, bool) {
	switch name {
	case string(PermissionNameRun):
		return &PermissionDescriptor{
			Name:    PermissionNameRun,
			Command: target,
		}, true
	case string(PermissionNameRead):
		return &PermissionDescriptor{
			Name: PermissionNameRead,
			Path: target,
		}, true
	case string(PermissionNameWrite):
		return &PermissionDescriptor{
			Name: PermissionNameWrite,
			Path: target,
		}, true
	case string(PermissionNameNet), "network":
		hostport, err := ParseHostPort(target)
		if err != nil {
			return nil, false
		}
		return &PermissionDescriptor{
			Name: PermissionNameNet,
			Host: hostport,
		}, true
	case string(PermissionNameEnv):
		return &PermissionDescriptor{
			Name:     PermissionNameEnv,
			Variable: target,
		}, true
	case string(PermissionNameSys):
		kind, ok := ParseSysKind(target)
		if !ok {
			return nil, false
		}
		return &PermissionDescriptor{
			Name: PermissionNameSys,
			Kind: kind,
		}, true
	case string(PermissionNameFfi):
		return &PermissionDescriptor{
			Name: PermissionNameFfi,
			Path: target,
		}, true
	case string(PermissionNameHrtime):
		return &PermissionDescriptor{
			Name: PermissionNameHrtime,
		}, true
	default:
		return nil, false
	}
}
