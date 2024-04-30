package deno

import (
	"regexp"
)

// On deno < 1.31.0
// The permission prompt looks like
//
// ⚠️  ┌ Deno requests net access to "0.0.0.0:8080".
//    ├ Requested by `Deno.listen()` API
//    ├ Run again with --allow-net to bypass this prompt.
//    └ Allow? [y/n] (y = yes, allow; n = no, deny) >
//
// On deno >= 1.31.0
// The permission prompt looks like
//
// ┌ ⚠️  Deno requests net access to "0.0.0.0:8080".
// ├ Requested by `Deno.listen()` API
// ├ Run again with --allow-net to bypass this prompt.

var accessToRegexp = regexp.MustCompile(`Deno requests (.+) access to "(.+)"\.`)
var allAccessRegexp = regexp.MustCompile(`Deno requests (.+) access\.`)
var hrtimeRegexp = regexp.MustCompile(`Deno requests access to high precision time\.`)

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
