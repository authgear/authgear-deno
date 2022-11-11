package deno

import (
	"regexp"
)

var accessToRegexp = regexp.MustCompile(`⚠️  ┌ Deno requests (.+) access to "(.+)"\.`)
var allAccessRegexp = regexp.MustCompile(`⚠️  ┌ Deno requests (.+) access\.`)
var hrtimeRegexp = regexp.MustCompile(`⚠️  ┌ Deno requests access to high precision time\.`)

func LineToPermissionDescriptor(line string) (*PermissionDescriptor, bool) {
	if matches := hrtimeRegexp.FindStringSubmatch(line); len(matches) == 1 {
		return &PermissionDescriptor{
			Name: PermissionNameHrtime,
		}, true
	}

	if matches := accessToRegexp.FindStringSubmatch(line); len(matches) == 3 {
		rawName := matches[1]
		target := matches[2]
		return ParsePermissionDescriptor(rawName, target)
	}

	if matches := allAccessRegexp.FindStringSubmatch(line); len(matches) == 2 {
		rawName := matches[1]
		return ParsePermissionDescriptor(rawName, "")
	}

	return nil, false
}
