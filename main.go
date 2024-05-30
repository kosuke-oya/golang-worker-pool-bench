package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
)

// atomicで足し合わせるための変数
var sum uint64

func main() {
	fmt.Println("Hello, World!")
}

func benchCpuBound() {
	for i := 0; i < 100; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(100))
		if err != nil {
			panic(err)
		}
		// []byteに変換する
		b := n.Bytes()
		s := sha256.Sum256(b)
		// []byteを2進数に変換して足し合わせる
		atomic.AddUint64(&sum, binary.BigEndian.Uint64(s[:]))
	}
}
func benchCpuBoundMultiGoroutine() {
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			n, err := rand.Int(rand.Reader, big.NewInt(100))
			if err != nil {
				panic(err)
			}
			// []byteに変換する
			b := n.Bytes()
			s := sha256.Sum256(b)
			// []byteを2進数に変換して足し合わせる
			atomic.AddUint64(&sum, binary.BigEndian.Uint64(s[:]))
		}()
	}
	wg.Wait()
}
