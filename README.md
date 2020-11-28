# sem
[![PkgGoDev](https://pkg.go.dev/badge/github.com/hslam/sem)](https://pkg.go.dev/github.com/hslam/sem)
[![Build Status](https://travis-ci.org/hslam/sem.svg?branch=master)](https://travis-ci.org/hslam/sem)
[![Go Report Card](https://goreportcard.com/badge/github.com/hslam/sem)](https://goreportcard.com/report/github.com/hslam/sem)
[![LICENSE](https://img.shields.io/github/license/hslam/sem.svg?style=flat-square)](https://github.com/hslam/sem/blob/master/LICENSE)

Package sem provides a way to use System V semaphore.

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
	semid, err := sem.Get(key)
	if err != nil {
		panic(err)
	}
	defer sem.Remove(semid)
	if r, err := sem.GetValue(semid); err != nil {
		panic(err)
	} else if r == 0 {
		fmt.Printf("%s wait\n", time.Now().Format("15:04:05"))
	}
	ok, err := sem.P(semid, sem.SEM_UNDO)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s P %t\n", time.Now().Format("15:04:05"), ok)
	time.Sleep(time.Second * 10)
	ok, err = sem.V(semid, sem.SEM_UNDO)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s V %t\n", time.Now().Format("15:04:05"), ok)
	time.Sleep(time.Second * 20)
}
```

#### Output

```sh
$ go run main.go
12:35:21 P true
12:35:31 V true
```
In another terminal.
```sh
$ go run main.go
12:35:25 wait
12:35:31 P true
12:35:41 V true
```

### License
This package is licensed under a MIT license (Copyright (c) 2020 Meng Huang)


### Author
shm was written by Meng Huang.


