package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/alitto/pond"
	"github.com/gammazero/workerpool"
	"github.com/panjf2000/ants/v2"
)

// var nUnder int = 10000
var sum uint64
var maxWorkers = 100
var loop int = 100

func benchCpuBound() {
	for i := 0; i < loop; i++ {
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
	for i := 0; i < loop; i++ {
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

// func benchIOBound() {
// 	time.Sleep(1 * time.Second)
// }
// func benchIoCpuComplex() {
// 	benchCpuBound()
// 	benchIOBound()
// }

func TestBenchCpu(t *testing.T) {
	benchCpuBound()
	fmt.Println(sum)
}
func TestBenchCpuM(t *testing.T) {
	benchCpuBoundMultiGoroutine()
	fmt.Println(sum)
}

// func TestBenchIO(t *testing.T) {
// 	benchIOBound()
// }
// func TestBenchIOCpu(t *testing.T) {
// 	benchIoCpuComplex()
// }

func BenchmarkGanmaPoolCpu(b *testing.B) {
	var wp = workerpool.New(maxWorkers)
	for i := 0; i < b.N; i++ {
		wp.Submit(benchCpuBound)
	}
	wp.StopWait()
}
func BenchmarkPondPoolCpuBound(b *testing.B) {
	var wp = pond.New(maxWorkers, 0, pond.Strategy(pond.Eager()))
	defer wp.StopAndWait()
	for i := 0; i < b.N; i++ {
		wp.Submit(benchCpuBound)
	}
}
func BenchmarkAntsPoolCpu(b *testing.B) {
	var wg sync.WaitGroup
	ff := func() {
		benchCpuBound()
		defer wg.Done()
	}
	antsPool, _ := ants.NewPool(maxWorkers)
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		antsPool.Submit(ff)
	}
	wg.Wait()
}
func BenchmarkGoroutineCpu(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			benchCpuBound()
		}()
	}
	wg.Wait()
}
func BenchmarkGanmaPoolCpuMulriGoroutine(b *testing.B) {
	var wp = workerpool.New(maxWorkers)
	for i := 0; i < b.N; i++ {
		wp.Submit(benchCpuBoundMultiGoroutine)
	}
	wp.StopWait()
}
func BenchmarkPondPoolCpuMultiGoroutine(b *testing.B) {
	var wp = pond.New(maxWorkers, 0, pond.Strategy(pond.Eager()))
	defer wp.StopAndWait()
	for i := 0; i < b.N; i++ {
		wp.Submit(benchCpuBoundMultiGoroutine)
	}
}
func BenchmarkAntsPoolCpuMultiGoroutine(b *testing.B) {
	var wg sync.WaitGroup
	ff := func() {
		benchCpuBoundMultiGoroutine()
		defer wg.Done()
	}
	antsPool, _ := ants.NewPool(maxWorkers)
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		antsPool.Submit(ff)
	}
	wg.Wait()
}
func BenchmarkGoroutineCpuMultiGoroutine(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			benchCpuBoundMultiGoroutine()
		}()
	}
	wg.Wait()
}

// func BenchmarkGanmaPoolIO(b *testing.B) {
// 	var wp = workerpool.New(maxWorkers)
// 	for i := 0; i < b.N; i++ {
// 		wp.Submit(benchIOBound)
// 	}
// 	wp.StopWait()
// }
// func BenchmarkPondPoolIO(b *testing.B) {
// 	var wp = pond.New(maxWorkers, 0, pond.Strategy(pond.Eager()))
// 	defer wp.StopAndWait()
// 	for i := 0; i < b.N; i++ {
// 		wp.Submit(benchIOBound)
// 	}
// }
// func BenchmarkAntsPoolIO(b *testing.B) {
// 	var wg sync.WaitGroup
// 	ff := func() {
// 		benchIOBound()
// 		defer wg.Done()
// 	}
// 	antsPool, _ := ants.NewPool(maxWorkers)
// 	for i := 0; i < b.N; i++ {
// 		wg.Add(1)
// 		antsPool.Submit(ff)
// 	}
// 	wg.Wait()
// }
// func BenchmarkGoroutineIO(b *testing.B) {
// 	var wg sync.WaitGroup
// 	for i := 0; i < b.N; i++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			benchIOBound()
// 		}()
// 	}
// 	wg.Wait()
// }
// func BenchmarkGanmaPoolIoCpuComplex(b *testing.B) {
// 	var wp = workerpool.New(maxWorkers)
// 	for i := 0; i < b.N; i++ {
// 		wp.Submit(benchIoCpuComplex)
// 	}
// 	wp.StopWait()
// }
// func BenchmarkPondPoolCpuIOComplex(b *testing.B) {
// 	var wp = pond.New(maxWorkers, 0, pond.Strategy(pond.Eager()))
// 	defer wp.StopAndWait()
// 	for i := 0; i < b.N; i++ {
// 		wp.Submit(benchIoCpuComplex)
// 	}
// }
// func BenchmarkAntsPoolCpuIOComplex(b *testing.B) {
// 	var wg sync.WaitGroup
// 	ff := func() {
// 		benchIoCpuComplex()
// 		defer wg.Done()
// 	}
// 	mp, _ := ants.NewPool(maxWorkers)
// 	for i := 0; i < b.N; i++ {
// 		wg.Add(1)
// 		mp.Submit(ff)
// 	}
// 	wg.Wait()
// }
// func BenchmarkGoroutineCpuIOComplex(b *testing.B) {
// 	var wg sync.WaitGroup
// 	for i := 0; i < b.N; i++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			benchIoCpuComplex()
// 		}()
// 	}
// 	wg.Wait()
// }
