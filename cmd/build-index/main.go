package main

import (
	"fmt"
	"gin-test/internal/ivf"
	"log"
	"os"
	"runtime"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: build-index <references.json.gz> <out.bin>")
		os.Exit(2)
	}
	in := os.Args[1]
	t0 := time.Now()
	ivf.WritePartitions(in, 4)

	/* t0 := time.Now()
	store.LoadReferencesGzipedPath(in, 600000)
	data := store.References
	log.Printf("load: %d vectors in %s", len(data), time.Since(t0))

	t1 := time.Now()
	centroids := ivf.TrainKMeans(data, 4096)
	log.Printf("train: k=%d sample=%d iters=%d in %s",
		4096, ivf.SampleSize, ivf.MaxIters, time.Since(t1))

	t2 := time.Now()
	idx := ivf.BuildIndex(data, centroids)
	log.Printf("build: %d vectors in %s", idx.N, time.Since(t2))

	t3 := time.Now()
	if err := idx.SerializeToFile(out); err != nil {
		log.Fatalf("save: %v", err)
	}
	info, _ := os.Stat(out)
	log.Printf("saved %d bytes (%.1f MB) to %s in %s",
		info.Size(), float64(info.Size())/(1<<20), out, time.Since(t3))

	t4 := time.Now()
	quantIdx := ivf.QuantizeIVFIndex(idx)
	var quantOut string
	if strings.HasSuffix(out, ".bin") {
		quantOut = out[:len(out)-4] + "_quantized.bin"
	} else {
		quantOut = out + "_quantized.bin"
	}
	if err := quantIdx.SerializeToFile(quantOut); err != nil {
		log.Fatalf("save quantized: %v", err)
	}
	quantInfo, _ := os.Stat(quantOut)
	log.Printf("saved %d bytes (%.1f MB) to %s in %s",
		quantInfo.Size(), float64(quantInfo.Size())/(1<<20), quantOut, time.Since(t4)) */

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log.Printf("heap inuse %.1f MB, alloc total %.1f MB",
		float64(m.HeapInuse)/(1<<20), float64(m.TotalAlloc)/(1<<20))
	log.Printf("total wall time: %s", time.Since(t0))
}
