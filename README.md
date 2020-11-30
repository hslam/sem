# sem
[![PkgGoDev](https://pkg.go.dev/badge/github.com/hslam/sem)](https://pkg.go.dev/github.com/hslam/sem)
[![Build Status](https://travis-ci.org/hslam/sem.svg?branch=master)](https://travis-ci.org/hslam/sem)
[![Go Report Card](https://goreportcard.com/badge/github.com/hslam/sem)](https://goreportcard.com/report/github.com/hslam/sem)
[![LICENSE](https://img.shields.io/github/license/hslam/sem.svg?style=flat-square)](https://github.com/hslam/sem/blob/master/LICENSE)

Package sem provides a way to use System V semaphores.

## Get started

### Install
```
go get github.com/hslam/sem
```
### Import
```
import "github.com/hslam/sem"
```
### Usage
####  Example
```go
package main

import (
	"fmt"
	"github.com/hslam/ftok"
	"github.com/hslam/sem"
	"time"
)

func main() {
	key, err := ftok.Ftok("/tmp", 0x22)
	if err != nil {
		panic(err)
	}
	nsems := 1
	semid, err := sem.Get(key, nsems, 0666)
	if err != nil {
		semid, err = sem.Get(key, nsems, sem.IPC_CREAT|sem.IPC_EXCL|0666)
		if err != nil {
			panic(err)
		}
		defer sem.Remove(semid)
		for semnum := 0; semnum < nsems; semnum++ {
			_, err := sem.SetValue(semid, semnum, 1)
			if err != nil {
				panic(err)
			}
		}
	}
	semnum := 0
	if count, err := sem.GetValue(semid, semnum); err != nil {
		panic(err)
	} else if count == 0 {
		fmt.Printf("%s semnum %d wait\n", time.Now().Format("15:04:05"), semnum)
	}
	ok, err := sem.P(semid, semnum, sem.SEM_UNDO)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s semnum %d P %t\n", time.Now().Format("15:04:05"), semnum, ok)
	time.Sleep(time.Second * 10)
	ok, err = sem.V(semid, semnum, sem.SEM_UNDO)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s semnum %d V %t\n", time.Now().Format("15:04:05"), semnum, ok)
	time.Sleep(time.Second * 20)
}
```

#### Output

```sh
$ go run main.go
12:35:21 semnum 0 P true
12:35:31 semnum 0 V true
```
In another terminal.
```sh
$ go run main.go
12:35:25 semnum 0 wait
12:35:31 semnum 0 P true
12:35:41 semnum 0 V true
```

### License
This package is licensed under a MIT license (Copyright (c) 2020 Meng Huang)


### Author
sem was written by Meng Huang.


