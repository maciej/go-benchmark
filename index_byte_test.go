package benchmark

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
)

func loopIndex(s []byte, sep byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			return i
		}
	}
	return -1
}

func loopSplit(s []byte, sep byte) [][]byte {
	n := bytes.Count(s, []byte{'/'}) + 1

	a := make([][]byte, n)
	n--
	i := 0
	for i < n {
		m := loopIndex(s, sep)
		if m < 0 {
			break
		}
		a[i] = s[:m:m]
		s = s[m+1:]
		i++
	}
	a[i] = s
	return a[:i+1]
}

func indexByteSplit(s []byte, sep byte) [][]byte {
	n := bytes.Count(s, []byte{'/'}) + 1

	a := make([][]byte, n)
	n--
	i := 0
	for i < n {
		m := bytes.IndexByte(s, sep)
		if m < 0 {
			break
		}
		a[i] = s[:m:m]
		s = s[m+1:]
		i++
	}
	a[i] = s
	return a[:i+1]

}

func BenchmarkIndexByte(b *testing.B) {
	const maxSliceLength = 128
	rnd := rand.New(rand.NewSource(42))

	b.Run("last", func(b *testing.B) {
		for j := 2; j <= maxSliceLength; j *= 2 {
			s := make([]byte, j)
			for k := 0; k < j-2; k++ {
				s[k] = 'a' + byte(rnd.Int31n(26))
			}
			s[j-1] = '/'

			b.Run(fmt.Sprintf("stdlib %d", j), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					bytes.IndexByte(s, '/')
				}
			})

			b.Run(fmt.Sprintf("loop %d", j), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					loopIndex(s, '/')
				}
			})
		}
	})
}

func BenchmarkSplit(b *testing.B) {
	const sliceLength = 4096
	rnd := rand.New(rand.NewSource(42))

	b.Run("regular", func(b *testing.B) {
		s := make([]byte, sliceLength)

		for k := 1; k < sliceLength; k *= 2 {
			n := 0
			// initialize slice with characters
			for j := 0; j < sliceLength; j++ {
				if j%k == 0 {
					n++
					s[j] = '/'
				} else {
					s[j] = 'a' + byte(rnd.Int31n(26))
				}
			}

			b.Run(fmt.Sprintf("stdlib  %d/%d (every %d bytes)", n, sliceLength, k), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					bytes.Split(s, []byte{'/'})
				}
			})

			b.Run(fmt.Sprintf("loop %d/%d (every %d bytes)", n, sliceLength, k), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					loopSplit(s, '/')
				}
			})

			b.Run(fmt.Sprintf("index_byte %d/%d (every %d bytes)", n, sliceLength, k), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					indexByteSplit(s, '/')
				}
			})
		}

	})
}
