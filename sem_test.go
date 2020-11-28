// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package sem

import (
	"github.com/hslam/ftok"
	"testing"
	"time"
)

func TestSem(t *testing.T) {
	done := make(chan struct{})
	go func() {
		key, err := ftok.Ftok("/tmp", 0x22)
		if err != nil {
			panic(err)
		}
		semid, err := Get(key)
		if err != nil {
			panic(err)
		}
		defer Remove(semid)
		if r, err := GetValue(semid); err != nil {
			panic(err)
		} else if r == 0 {
			t.Error("wait\n")
		}
		ok, err := P(semid, IPC_NOWAIT|SEM_UNDO)
		if err != nil {
			t.Error(err)
		}
		if !ok {
			t.Error("P failed!\n")
		}
		time.Sleep(time.Millisecond * 200)
		ok, err = V(semid, IPC_NOWAIT|SEM_UNDO)
		if err != nil {
			t.Error(err)
		}
		if !ok {
			t.Error("V failed!\n")
		}
		time.Sleep(time.Millisecond * 100)
		close(done)
	}()
	time.Sleep(time.Millisecond * 100)
	key, err := ftok.Ftok("/tmp", 0x22)
	if err != nil {
		panic(err)
	}
	semid, err := Get(key)
	if err != nil {
		panic(err)
	}
	if r, err := GetValue(semid); err != nil {
		panic(err)
	} else if r > 0 {
		t.Error()
	}
	ok, err := P(semid, IPC_NOWAIT|SEM_UNDO)
	if err != nil {
		t.Error(err)
	}
	if ok {
		t.Error()
	}
	<-done
}
