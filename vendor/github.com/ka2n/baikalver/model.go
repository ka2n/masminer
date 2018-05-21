package baikalver

import (
	"strconv"
)

// Model of ASIC
type Model uint8

// List of Model
const (
	Unknown  Model = 0
	Mini     Model = 0x11
	Giant    Model = 0x12
	Cube     Model = 0x20
	CubeRev2 Model = 0x21
	GiantP   Model = 0x22
	GX10     Model = 0x71
	GN20     Model = 0x51
	GN40     Model = 0x52
	GB       Model = 0x91
)

func (i Model) String() string {
	switch i {
	case Mini:
		return "Mini"
	case Giant:
		return "Giant"
	case Cube:
		return "Cube"
	case CubeRev2:
		return "Cube"
	case GiantP:
		return "Giant+"
	case GX10:
		return "GX10"
	case GN20:
		return "GN20"
	case GN40:
		return "GN40"
	case GB:
		return "GB"
	}
	return ""
}

func (i Model) Equal(v Model) bool {
	if i == Cube && v == CubeRev2 {
		return true
	}
	return i == v
}

// ModelFromHWV returns Model from HWV string
func ModelFromHWV(s string) (Model, error) {
	hwv, err := strconv.Atoi(s)
	if err != nil {
		return Unknown, err
	}
	return Model(uint8(hwv)), nil
}
