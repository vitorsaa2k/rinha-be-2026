package ivf

import (
	"encoding/binary"
	"fmt"
	"gin-test/models"
	"io"
	"math"
	"os"
	"runtime"
	"sync"
	"unsafe"
)

const (
	magic      = "IVF1"
	version    = uint32(1)
	dims       = 14
	magicLen   = 4
	headerSize = 16 // magic(4) + version(4) + K(4) + N(4)
)

// IVFIndex is the serializable representation of the IVF index.
type IVFIndex struct {
	K uint32 // number of clusters
	N uint32 // total vectors

	Centroids []float32 // K * dims flat array
	Offsets   []uint32  // K+1 offsets into Vectors/Labels arrays
	Vectors   []float32 // N * dims flat array, grouped by cluster
	Labels    []uint8   // N bytes, 0=legit 1=fraud
}

// Serialize writes the index to w in binary format.
func (idx *IVFIndex) Serialize(w io.Writer) error {
	header := make([]byte, headerSize)
	copy(header[0:magicLen], magic)
	binary.LittleEndian.PutUint32(header[4:8], version)
	binary.LittleEndian.PutUint32(header[8:12], idx.K)
	binary.LittleEndian.PutUint32(header[12:16], idx.N)
	if _, err := w.Write(header); err != nil {
		return err
	}
	if err := writeFloat32Slice(w, idx.Centroids); err != nil {
		return err
	}
	if err := writeUint32Slice(w, idx.Offsets); err != nil {
		return err
	}
	if err := writeFloat32Slice(w, idx.Vectors); err != nil {
		return err
	}
	if _, err := w.Write(idx.Labels); err != nil {
		return err
	}
	return nil
}

// SerializeToFile is a convenience wrapper for Serialize.
func (idx *IVFIndex) SerializeToFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return idx.Serialize(f)
}

// Load reads an index from the binary format into heap allocations.
func Load(r io.Reader) (*IVFIndex, error) {
	header := make([]byte, headerSize)
	if _, err := io.ReadFull(r, header); err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}
	if string(header[0:magicLen]) != magic {
		return nil, fmt.Errorf("bad magic: %q", header[0:magicLen])
	}
	if v := binary.LittleEndian.Uint32(header[4:8]); v != version {
		return nil, fmt.Errorf("bad version: %d (want %d)", v, version)
	}
	idx := &IVFIndex{
		K: binary.LittleEndian.Uint32(header[8:12]),
		N: binary.LittleEndian.Uint32(header[12:16]),
	}
	centroidsLen := int(idx.K) * dims
	idx.Centroids = make([]float32, centroidsLen)
	if err := readFloat32Slice(r, idx.Centroids); err != nil {
		return nil, fmt.Errorf("read centroids: %w", err)
	}
	idx.Offsets = make([]uint32, int(idx.K)+1)
	if err := readUint32Slice(r, idx.Offsets); err != nil {
		return nil, fmt.Errorf("read offsets: %w", err)
	}
	vectorsLen := int(idx.N) * dims
	idx.Vectors = make([]float32, vectorsLen)
	if err := readFloat32Slice(r, idx.Vectors); err != nil {
		return nil, fmt.Errorf("read vectors: %w", err)
	}
	idx.Labels = make([]uint8, idx.N)
	if _, err := io.ReadFull(r, idx.Labels); err != nil {
		return nil, fmt.Errorf("read labels: %w", err)
	}
	return idx, nil
}

// LoadFromFile is a convenience wrapper for Load.
func LoadFromFile(path string) (*IVFIndex, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Load(f)
}

// ToClusters converts the serialized index back to the runtime []Cluster format.
func (idx *IVFIndex) ToClusters() []Cluster {
	clusters := make([]Cluster, idx.K)
	for c := range clusters {
		centroidStart := c * dims
		clusters[c].Centroid = make([]float32, dims)
		copy(clusters[c].Centroid, idx.Centroids[centroidStart:centroidStart+dims])
		start := idx.Offsets[c]
		end := idx.Offsets[c+1]
		count := int(end - start)
		clusters[c].Lists = make([]models.DatasetStruct, count)
		for i := 0; i < count; i++ {
			globalIdx := int(start) + i
			vecStart := globalIdx * dims
			vec := make([]float32, dims)
			copy(vec, idx.Vectors[vecStart:vecStart+dims])
			label := "legit"
			if idx.Labels[globalIdx] == 1 {
				label = "fraud"
			}
			clusters[c].Lists[i] = models.DatasetStruct{
				Vector: vec,
				Label:  label,
			}
		}
	}
	return clusters
}

// BuildIndex constructs an IVFIndex from trained centroids and the full dataset.
func BuildIndex(data []models.DatasetStruct, centroids []Cluster) *IVFIndex {
	k := uint32(len(centroids))
	n := uint32(len(data))
	assign := make([]uint32, n)
	if totalCpus := runtime.NumCPU(); totalCpus > 4 {
		computeVectorsToCentroidsMultiThread(&assign, data, centroids, 4)
	} else {
		for i, v := range data {
			best := 0
			bestD := float64(-1)
			for c, cl := range centroids {
				d := SquaredDistance(cl.Centroid, v.Vector)
				if bestD < 0 || d < bestD {
					bestD = d
					best = c
				}
			}
			assign[i] = uint32(best)
		}
	}

	idx := &IVFIndex{
		K: k,
		N: n,
	}

	buckets := make([][]uint32, k)
	for c := range buckets {
		var approxSize uint32 = n/k + 8
		buckets[c] = make([]uint32, 0, approxSize)
	}
	for i := uint32(0); i < n; i++ {
		c := assign[i]
		buckets[c] = append(buckets[c], i)
	}

	idx.Offsets = make([]uint32, k+1)
	var offset uint32
	for c := 0; c < int(k); c++ {
		idx.Offsets[c] = offset
		offset += uint32(len(buckets[c]))
	}
	idx.Offsets[k] = offset

	idx.Centroids = make([]float32, int(k)*dims)
	for c := 0; c < int(k); c++ {
		start := c * dims
		copy(idx.Centroids[start:start+dims], centroids[c].Centroid)
	}

	idx.Vectors = make([]float32, int(n)*dims)
	idx.Labels = make([]uint8, n)
	for c := 0; c < int(k); c++ {
		base := idx.Offsets[c]
		for i, id := range buckets[c] {
			globalIdx := int(base) + i
			dstStart := globalIdx * dims
			copy(idx.Vectors[dstStart:dstStart+dims], data[id].Vector)
			if data[id].Label == "fraud" {
				idx.Labels[globalIdx] = 1
			}
		}
	}

	return idx
}

// SizeOnDisk returns the approximate file size of the serialized index.
func (idx *IVFIndex) SizeOnDisk() int {
	return headerSize +
		len(idx.Centroids)*4 +
		len(idx.Offsets)*4 +
		len(idx.Vectors)*4 +
		len(idx.Labels)
}

func writeFloat32Slice(w io.Writer, s []float32) error {
	b := make([]byte, 4*len(s))
	for i, v := range s {
		binary.LittleEndian.PutUint32(b[4*i:4*i+4], math.Float32bits(v))
	}
	if len(b) == 0 {
		return nil
	}
	_, err := w.Write(b)
	return err
}

func writeUint32Slice(w io.Writer, s []uint32) error {
	b := make([]byte, 4*len(s))
	for i, v := range s {
		binary.LittleEndian.PutUint32(b[4*i:4*i+4], v)
	}
	if len(b) == 0 {
		return nil
	}
	_, err := w.Write(b)
	return err
}

func readFloat32Slice(r io.Reader, s []float32) error {
	b := make([]byte, 4*len(s))
	if len(b) == 0 {
		return nil
	}
	if _, err := io.ReadFull(r, b); err != nil {
		return err
	}
	for i := range s {
		s[i] = math.Float32frombits(binary.LittleEndian.Uint32(b[4*i : 4*i+4]))
	}
	return nil
}

func readUint32Slice(r io.Reader, s []uint32) error {
	b := make([]byte, 4*len(s))
	if len(b) == 0 {
		return nil
	}
	if _, err := io.ReadFull(r, b); err != nil {
		return err
	}
	for i := range s {
		s[i] = binary.LittleEndian.Uint32(b[4*i : 4*i+4])
	}
	return nil
}

func computeVectorsToCentroidsMultiThread(assign *[]uint32, data []models.DatasetStruct, centroids []Cluster, threads int) {
	n := len(data)
	var wg sync.WaitGroup
	numWorkers := threads
	chunkSize := (int(n) + numWorkers - 1) / numWorkers

	for w := range numWorkers {
		start := w * chunkSize
		end := start + chunkSize
		if start >= int(n) {
			break
		}
		if end > int(n) {
			end = int(n)
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i < end; i++ {
				best := 0
				bestD := float64(-1)
				for centIdx, centroid := range centroids {
					d := SquaredDistance(centroid.Centroid, data[i].Vector[:])
					if bestD < 0 || d < bestD {
						bestD = d
						best = centIdx
					}
				}
				(*assign)[i] = uint32(best)
			}
		}(start, end)
	}
	wg.Wait()
}

// ---------------------------------------------------------------------------
// Quantized index (int16, fixed scale QuantScale = 10000)
// ---------------------------------------------------------------------------

const (
	quantizedMagic      = "IVFQ"
	quantizedVersion    = uint32(2)
	quantizedHeaderSize = 16 // magic(4) + version(4) + K(4) + N(4)
)

type IVFQuantizedIndex struct {
	K uint32
	N uint32

	Centroids []int16  // K * dims
	Offsets   []uint32 // K + 1
	Vectors   []int16  // N * dims
	Labels    []uint8  // N
}

func QuantizeFloat32(v float32) int16 {
	s := v * QuantScale
	if s >= QuantScale {
		return QuantScale
	}
	if s <= -QuantScale {
		return -QuantScale
	}
	if s >= 0 {
		return int16(s + 0.5)
	}
	return int16(s - 0.5)
}

func QuantizeIVFIndex(idx *IVFIndex) *IVFQuantizedIndex {
	q := &IVFQuantizedIndex{
		K:         idx.K,
		N:         idx.N,
		Centroids: make([]int16, len(idx.Centroids)),
		Offsets:   make([]uint32, len(idx.Offsets)),
		Vectors:   make([]int16, len(idx.Vectors)),
		Labels:    make([]uint8, len(idx.Labels)),
	}
	copy(q.Offsets, idx.Offsets)
	copy(q.Labels, idx.Labels)

	for i, v := range idx.Centroids {
		q.Centroids[i] = QuantizeFloat32(v)
	}
	for i, v := range idx.Vectors {
		q.Vectors[i] = QuantizeFloat32(v)
	}
	return q
}

func (idx *IVFQuantizedIndex) Serialize(w io.Writer) error {
	header := make([]byte, quantizedHeaderSize)
	copy(header[0:magicLen], quantizedMagic)
	binary.LittleEndian.PutUint32(header[4:8], quantizedVersion)
	binary.LittleEndian.PutUint32(header[8:12], idx.K)
	binary.LittleEndian.PutUint32(header[12:16], idx.N)
	if _, err := w.Write(header); err != nil {
		return err
	}
	if err := writeInt16Slice(w, idx.Centroids); err != nil {
		return err
	}
	if err := writeUint32Slice(w, idx.Offsets); err != nil {
		return err
	}
	if err := writeInt16Slice(w, idx.Vectors); err != nil {
		return err
	}
	if _, err := w.Write(idx.Labels); err != nil {
		return err
	}
	return nil
}

func (idx *IVFQuantizedIndex) SerializeToFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return idx.Serialize(f)
}

func LoadQuantized(r io.Reader) (*IVFQuantizedIndex, error) {
	header := make([]byte, quantizedHeaderSize)
	if _, err := io.ReadFull(r, header); err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}
	if string(header[0:magicLen]) != quantizedMagic {
		return nil, fmt.Errorf("bad magic: %q", header[0:magicLen])
	}
	if v := binary.LittleEndian.Uint32(header[4:8]); v != quantizedVersion {
		return nil, fmt.Errorf("bad version: %d (want %d)", v, quantizedVersion)
	}
	idx := &IVFQuantizedIndex{
		K: binary.LittleEndian.Uint32(header[8:12]),
		N: binary.LittleEndian.Uint32(header[12:16]),
	}
	centroidsLen := int(idx.K) * dims
	idx.Centroids = make([]int16, centroidsLen)
	if err := readInt16Slice(r, idx.Centroids); err != nil {
		return nil, fmt.Errorf("read centroids: %w", err)
	}
	idx.Offsets = make([]uint32, int(idx.K)+1)
	if err := readUint32Slice(r, idx.Offsets); err != nil {
		return nil, fmt.Errorf("read offsets: %w", err)
	}
	vectorsLen := int(idx.N) * dims
	idx.Vectors = make([]int16, vectorsLen)
	if err := readInt16Slice(r, idx.Vectors); err != nil {
		return nil, fmt.Errorf("read vectors: %w", err)
	}
	idx.Labels = make([]uint8, idx.N)
	if _, err := io.ReadFull(r, idx.Labels); err != nil {
		return nil, fmt.Errorf("read labels: %w", err)
	}
	return idx, nil
}

func LoadQuantizedFromFile(path string) (*IVFQuantizedIndex, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return LoadQuantized(f)
}

func (idx *IVFQuantizedIndex) SizeOnDisk() int {
	return quantizedHeaderSize +
		len(idx.Centroids)*2 +
		len(idx.Offsets)*4 +
		len(idx.Vectors)*2 +
		len(idx.Labels)
}

func quantizedFileOffsets(K, N uint32) (offC, offO, offV, offL, totalSize int64) {
	offC = quantizedHeaderSize
	offO = offC + int64(K)*int64(dims)*2
	offV = offO + int64(K+1)*4
	offL = offV + int64(N)*int64(dims)*2
	totalSize = offL + int64(N)
	return
}

func LoadQuantizedMMap(path string) (*QuantizedPartition, error) {
	data, h, err := mmapFile(path)
	if err != nil {
		return nil, err
	}

	if string(data[:magicLen]) != quantizedMagic {
		munmapFile(data, h)
		return nil, fmt.Errorf("bad magic: %q", data[:magicLen])
	}
	if v := binary.LittleEndian.Uint32(data[4:8]); v != quantizedVersion {
		munmapFile(data, h)
		return nil, fmt.Errorf("bad version: %d (want %d)", v, quantizedVersion)
	}

	K := binary.LittleEndian.Uint32(data[8:12])
	N := binary.LittleEndian.Uint32(data[12:16])
	offC, offO, offV, offL, expectedSize := quantizedFileOffsets(K, N)

	if int64(len(data)) != expectedSize {
		munmapFile(data, h)
		return nil, fmt.Errorf("file size mismatch: %d (expected %d)", len(data), expectedSize)
	}

	idx := &IVFQuantizedIndex{
		K: K,
		N: N,
		Centroids: unsafe.Slice((*int16)(unsafe.Pointer(&data[offC])), int(K)*dims),
		Offsets:   unsafe.Slice((*uint32)(unsafe.Pointer(&data[offO])), int(K)+1),
		Vectors:   unsafe.Slice((*int16)(unsafe.Pointer(&data[offV])), int(N)*dims),
		Labels:    data[offL : offL+int64(N)],
	}
	return &QuantizedPartition{IVFQuantizedIndex: *idx, mmapData: data, mmapFile: h}, nil
}

func writeInt16Slice(w io.Writer, s []int16) error {
	b := make([]byte, 2*len(s))
	for i, v := range s {
		binary.LittleEndian.PutUint16(b[2*i:2*i+2], uint16(v))
	}
	if len(b) == 0 {
		return nil
	}
	_, err := w.Write(b)
	return err
}

func readInt16Slice(r io.Reader, s []int16) error {
	b := make([]byte, 2*len(s))
	if len(b) == 0 {
		return nil
	}
	if _, err := io.ReadFull(r, b); err != nil {
		return err
	}
	for i := range s {
		s[i] = int16(binary.LittleEndian.Uint16(b[2*i : 2*i+2]))
	}
	return nil
}
