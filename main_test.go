package main

import (
	"sync"
	"testing"

	"github.com/alitto/pond"
	"github.com/gammazero/workerpool"
	"github.com/panjf2000/ants/v2"
)

// poolのworker数
var maxWorkers = 1000

// taskの数
var tasks int = 100000

func BenchmarkGanmaPoolCpu(b *testing.B) {
	var wp = workerpool.New(maxWorkers)
	for i := 0; i < b.N; i++ {
		for i := 0; i < tasks; i++ {
			wp.Submit(benchCpuBound)
		}
		wp.StopWait()
	}
}
func BenchmarkPondPoolCpuBound(b *testing.B) {
	var wp = pond.New(maxWorkers, 0)
	defer wp.StopAndWait()
	for i := 0; i < b.N; i++ {
		for i := 0; i < tasks; i++ {
			wp.Submit(benchCpuBound)
		}
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
		for i := 0; i < tasks; i++ {
			wg.Add(1)
			antsPool.Submit(ff)
		}
		wg.Wait()
	}
}
func BenchmarkGoroutineCpu(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		for i := 0; i < tasks; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				benchCpuBound()
			}()
		}
		wg.Wait()
	}
}
func BenchmarkGanmaPoolCpuMulriGoroutine(b *testing.B) {
	var wp = workerpool.New(maxWorkers)
	for i := 0; i < b.N; i++ {
		for i := 0; i < tasks; i++ {
			wp.Submit(benchCpuBoundMultiGoroutine)
		}
		wp.StopWait()
	}
}
func BenchmarkPondPoolCpuMultiGoroutine(b *testing.B) {
	var wp = pond.New(maxWorkers, 0)
	defer wp.StopAndWait()
	for i := 0; i < b.N; i++ {
		for i := 0; i < tasks; i++ {
			wp.Submit(benchCpuBoundMultiGoroutine)
		}
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
		for i := 0; i < tasks; i++ {
			wg.Add(1)
			antsPool.Submit(ff)
		}
		wg.Wait()
	}
}
func BenchmarkGoroutineCpuMultiGoroutine(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		for i := 0; i < tasks; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				benchCpuBoundMultiGoroutine()
			}()
		}
		wg.Wait()
	}
}
