// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package sem

import (
	"syscall"
	"unsafe"
)

const (
	// IPC_CREAT creates if key is nonexistent
	IPC_CREAT = 00001000

	// IPC_EXCL fails if key exists.
	IPC_EXCL = 00002000

	// IPC_NOWAIT returns error no wait.
	IPC_NOWAIT = 04000

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

	// SEM_UNDO sets up adjust on exit entry
	SEM_UNDO = 010000

	// IPC_RMID removes identifier
	IPC_RMID = 0
	// IPC_SET sets ipc_perm options.
	IPC_SET = 1
	// IPC_STAT gets ipc_perm options.
	IPC_STAT = 2
)

type sembuf struct {
	num uint16
	op  int16
	flg int16
}

// Get calls the semget system call.
func Get(key int) (uintptr, error) {
	semid, _, _ := syscall.Syscall(syscall.SYS_SEMGET, uintptr(key), 1, 0666)
	if int(semid) < 0 {
		semid, _, err := syscall.Syscall(syscall.SYS_SEMGET, uintptr(key), 1, IPC_CREAT|IPC_EXCL|00600)
		if int(semid) < 0 {
			return 0, err
		}
		var semun int = 1
		r1, _, err := syscall.Syscall6(syscall.SYS_SEMCTL, semid, 0, SETVAL, uintptr(semun), 0, 0)
		if int(r1) < 0 {
			return 0, err
		}
		return semid, nil
	}
	return semid, nil
}

// P calls the semop P system call.
func P(semid uintptr, flg int16) (bool, error) {
	if flg == 0 {
		flg = SEM_UNDO
	}
	buf := sembuf{num: 0, op: -1, flg: flg}
	r1, _, err := syscall.Syscall(syscall.SYS_SEMOP, semid, uintptr(unsafe.Pointer(&buf)), 1)
	var ok bool
	if r1 == 0 {
		ok = true
	}
	if err != 0 && err != syscall.EAGAIN {
		return ok, err
	}
	return ok, nil
}

// V calls the semop V system call.
func V(semid uintptr, flg int16) (bool, error) {
	if flg == 0 {
		flg = SEM_UNDO
	}
	buf := sembuf{num: 0, op: +1, flg: flg}
	r1, _, err := syscall.Syscall(syscall.SYS_SEMOP, semid, uintptr(unsafe.Pointer(&buf)), 1)
	var ok bool
	if r1 == 0 {
		ok = true
	}
	if err != 0 && err != syscall.EAGAIN {
		return ok, err
	}
	return ok, nil
}

// GetValue calls the semctl GETVAL system call.
func GetValue(semid uintptr) (int, error) {
	r1, _, err := syscall.Syscall(syscall.SYS_SEMCTL, semid, 0, GETVAL)
	if int(r1) < 0 {
		return 0, err
	}
	return int(r1), nil
}

// Remove removes the semaphore with the given id.
func Remove(semid uintptr) error {
	r1, _, errno := syscall.Syscall(syscall.SYS_SEMCTL, semid, IPC_RMID, 0)
	if int(r1) < 0 {
		return syscall.Errno(errno)
	}
	return nil
}
