package app

import (
	"fmt"

	"github.com/lmorg/murex/utils/semver"
)

// Name is the name of the $TERM
const Name = "mxtty"

// Version number of $TERM
// Format of version string should be "(major).(minor).(revision) DESCRIPTION"
const (
	version  = "%d.%d.%d"
	Major    = 0
	Minor    = 3
	Revision = 3000
)

const Title = "mxtty - Multimedia Terminal Emulator"

// Copyright is the copyright owner string
const Copyright = "© 2023-2024 Laurence Morgan"

// License is the projects software license
const License = "License GPL v2"

func init() {
	v = fmt.Sprintf(version, Major, Minor, Revision)
	sv, _ = semver.Parse(v)
}

var v string

func Version() string {
	return v
}

var sv *semver.Version

func Semver() *semver.Version {
	return sv
}

const ProjectSourcePath = "github.com/lmorg/mxtty/"
