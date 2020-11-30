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
	semnum := 0
	nsems := 1
	go func() {
		key, err := ftok.Ftok("/tmp", 0x22)
		if err != nil {
			panic(err)
		}
		semid, err := Get(key, nsems, 0666)
		if err != nil {
			semid, err = Get(key, nsems, IPC_CREAT|IPC_EXCL|0666)
			if err != nil {
				panic(err)
			}
			defer Remove(semid)
			for semnum := 0; semnum < nsems; semnum++ {
				_, err := SetValue(semid, semnum, 1)
				if err != nil {
					panic(err)
				}
			}
		}
		ok, err := P(semid, semnum, IPC_NOWAIT|SEM_UNDO)
		if err != nil {
			t.Error(err)
		}
		if !ok {
			t.Error("P failed!\n")
		}
		time.Sleep(time.Millisecond * 200)
		ok, err = V(semid, semnum, IPC_NOWAIT|SEM_UNDO)
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
	semid, err := Get(key, nsems, 0666)
	if err != nil {
		t.Error()
	}
	if r, err := GetValue(semid, semnum); err != nil {
		panic(err)
	} else if r > 0 {
		t.Error()
	}
	ok, err := P(semid, semnum, IPC_NOWAIT|SEM_UNDO)
	if err != nil {
		t.Error(err)
	}
	if ok {
		t.Error()
	}
	<-done
}
