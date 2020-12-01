// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// +build darwin linux dragonfly freebsd netbsd openbsd

package sem

import (
	"syscall"
	"unsafe"
)

const (
	// IPC_CREAT creates if key is nonexistent
	IPC_CREAT = 01000

	// IPC_EXCL fails if key exists.
	IPC_EXCL = 02000

	// IPC_NOWAIT returns error no wait.
	IPC_NOWAIT = 04000

	// IPC_PRIVATE is private key
	IPC_PRIVATE = 00000

	// SEM_UNDO sets up adjust on exit entry
	SEM_UNDO = 010000

	// IPC_RMID removes identifier
	IPC_RMID = 0
	// IPC_SET sets ipc_perm options.
	IPC_SET = 1
	// IPC_STAT gets ipc_perm options.
	IPC_STAT = 2
)

// Sembuf represents an operation.
type Sembuf struct {
	SemNum uint16
	SemOp  int16
	SemFlg int16
}

// Get calls the semget system call.
//
// The semget() system call returns the System V semaphore set identifier
// associated with the argument key.
//
// A new set of nsems semaphores is created if key has the value
// IPC_PRIVATE or if no existing semaphore set is associated with key
// and IPC_CREAT is specified in semflg.
//
// If semflg specifies both IPC_CREAT and IPC_EXCL and a semaphore set
// already exists for key, then semget() fails with errno set to EEXIST.
//
// The argument nsems can be 0 (a don't care) when a semaphore set is
// not being created.  Otherwise, nsems must be greater than 0 and less
// than or equal to the maximum number of semaphores per semaphore set.
//
// If successful, the return value will be the semaphore set identifier,
// otherwise, -1 is returned, with errno indicating the error.
func Get(key int, nsems int, semflg int) (int, error) {
	r1, _, err := syscall.Syscall(syscall.SYS_SEMGET, uintptr(key), uintptr(nsems), uintptr(semflg))
	semid := int(r1)
	if semid < 0 {
		return semid, err
	}
	return semid, nil
}

// SetValue calls the semctl SETVAL system call.
func SetValue(semid int, semnum int, semun int) (bool, error) {
	r1, _, err := syscall.Syscall6(syscall.SYS_SEMCTL, uintptr(semid), uintptr(semnum), SETVAL, uintptr(semun), 0, 0)
	if int(r1) < 0 {
		return false, err
	}
	return true, nil
}

// GetValue calls the semctl GETVAL system call.
func GetValue(semid int, semnum int) (int, error) {
	r1, _, err := syscall.Syscall(syscall.SYS_SEMCTL, uintptr(semid), uintptr(semnum), GETVAL)
	count := int(r1)
	if count < 0 {
		return count, err
	}
	return count, nil
}

// P calls the semop P system call.
// Flags recognized in semflg are IPC_NOWAIT and SEM_UNDO.
// If an operation specifies SEM_UNDO, it will be automatically undone when the
// process terminates.
func P(semid int, semnum int, semflg int) (bool, error) {
	return op(semid, uint16(semnum), -1, int16(semflg))
}

// V calls the semop V system call.
// Flags recognized in semflg are IPC_NOWAIT and SEM_UNDO.
// If an operation specifies SEM_UNDO, it will be automatically undone when the
// process terminates.
func V(semid int, semnum int, semflg int) (bool, error) {
	return op(semid, uint16(semnum), 1, int16(semflg))
}

func op(semid int, semnum uint16, semop, semflg int16) (bool, error) {
	if semflg == 0 {
		semflg = SEM_UNDO
	}
	var sops [1]Sembuf
	sops[0] = Sembuf{SemNum: semnum, SemOp: semop, SemFlg: semflg}
	return Operate(semid, sops[:])
}

// Op calls the semop system call.
//
// semop() performs operations on selected semaphores in the set indi‐
// cated by semid.  Each of the nsops elements in the array pointed to
// by sops is a structure that specifies an operation to be performed on
// a single semaphore.  The elements of this structure are of type
// struct sembuf, containing the following members:
//
// unsigned short sem_num;  /* semaphore number */
// short          sem_op;   /* semaphore operation */
// short          sem_flg;  /* operation flags */
//
// Flags recognized in sem_flg are IPC_NOWAIT and SEM_UNDO.
// If an operation specifies SEM_UNDO, it will be automatically undone when the
// process terminates.
//
// The set of operations contained in sops is performed in array order,
// and atomically, that is, the operations are performed either as a
// complete unit, or not at all.  The behavior of the system call if not
// all operations can be performed immediately depends on the presence
// of the IPC_NOWAIT flag in the individual sem_flg fields, as noted be‐
// low.
func Op(semid int, sops uintptr, nsops int) (bool, error) {
	r1, _, err := syscall.Syscall(syscall.SYS_SEMOP, uintptr(semid), sops, uintptr(nsops))
	var ok = true
	if int(r1) < 0 {
		ok = false
	}
	if err != 0 && err != syscall.EAGAIN {
		return ok, err
	}
	return ok, nil
}

// Operate calls the semop system call.
func Operate(semid int, sops []Sembuf) (bool, error) {
	return Op(semid, uintptr(unsafe.Pointer(&sops[0])), len(sops))
}

// Remove removes the semaphore set with the given id.
func Remove(semid int) error {
	r1, _, errno := syscall.Syscall(syscall.SYS_SEMCTL, uintptr(semid), IPC_RMID, 0)
	if int(r1) < 0 {
		return syscall.Errno(errno)
	}
	return nil
}
