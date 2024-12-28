package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// AOF = append only file
type Aof struct {
	file *os.File
	rd   *bufio.Reader
	mu sync.Mutex
}

var (
	aofInstance *Aof

	// reason for 2 locks
	// obj lock is for accessing and closing the AOF singleton
	// RW lock (inside the instance) is for controlling the I/O of the AOF singleton
	aofObjLock sync.Mutex
)

var filePath = "database.aof"

func GetAof() (*Aof, error) {

	aofObjLock.Lock()
	defer aofObjLock.Unlock()

	if aofInstance != nil {
		return aofInstance, nil
	}

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	aofInstance = &Aof{
		file: f,
		rd:   bufio.NewReader(f),
	}

	go aofInstance.syncAOF()

	return aofInstance, nil
}

func CloseAOF() {

	if aofInstance == nil {
		fmt.Println("aofInstance is nil, exiting close func")
		return
	}

	aofObjLock.Lock()
	defer aofObjLock.Unlock()

	if aofInstance != nil {

		// close I/O
		err := aofInstance.file.Close()
		if err != nil {
			fmt.Printf("aof close err: %v\n", err)
		}

		// set singleton to nil
		aofInstance = nil
	}

	fmt.Println("aof closed")

}

// start a goroutine to sync AOF to disk every 1 second
func (aof *Aof) syncAOF() {
	for {
		aof.mu.Lock()

		aofInstance.file.Sync()

		aof.mu.Unlock()

		time.Sleep(time.Second)
	}
}

func (aof *Aof) Write(value Value) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	_, err := aof.file.Write(value.Marshal())
	if err != nil {
		return err
	}

	return nil
}

func (aof *Aof) Read(callback func(value Value)) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	resp := NewResp(aof.file)

	for {
		value, err := resp.Read()
		if err != nil {
			callback(value)
		}

		if err == io.EOF {
			break
		}

		return err
	}

	return nil
}
