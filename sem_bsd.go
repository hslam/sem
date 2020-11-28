// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd netbsd openbsd

package sem

const (
	// GETNCNT returns the value of semncnt {READ}.
	GETNCNT = 3
	// GETPID returns the value of sempid {READ}
	GETPID = 4
	// GETVAL returns the value of semval {READ}
	GETVAL = 5
	// GETALL returns semvals into arg.array {READ}
	GETALL = 6
	// GETZCNT returns the value of semzcnt {READ}
	GETZCNT = 7
	// SETVAL sets the value of semval to arg.val {ALTER}
	SETVAL = 8
	// SETALL sets semvals from arg.array {ALTER}
	SETALL = 9
)
