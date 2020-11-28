// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// +build linux

package sem

const (
	// GETNCNT returns the value of semncnt {READ}.
	GETNCNT = 14
	// GETPID returns the value of sempid {READ}
	GETPID = 11
	// GETVAL returns the value of semval {READ}
	GETVAL = 12
	// GETALL returns semvals into arg.array {READ}
	GETALL = 13
	// GETZCNT returns the value of semzcnt {READ}
	GETZCNT = 15
	// SETVAL sets the value of semval to arg.val {ALTER}
	SETVAL = 16
	// SETALL sets semvals from arg.array {ALTER}
	SETALL = 17
)
