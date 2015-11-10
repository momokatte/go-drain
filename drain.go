// Copyright 2015 Michael O'Rourke. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package drain

import (
	"bufio"
	"io"
	"runtime"
)

const (
	NewLineString     = string(0x0A)
	DefaultBufferSize = 64 * 1024
)

// Scans all tokens from an io.Reader into a channel until EOF
//
func ReaderToChan(r io.Reader, split bufio.SplitFunc, dest chan<- string) error {
	s := bufio.NewScanner(r)
	s.Split(split)
	return ScannerToChan(s, dest)
}

// Scans all tokens into a channel until EOF
//
func ScannerToChan(s *bufio.Scanner, dest chan<- string) error {
	for s.Scan() {
		dest <- s.Text()
	}
	return s.Err()
}

// Writes all strings from a channel to an io.Writer until the channel is closed. When the channel is empty,
// the execution thread is yielded to another goroutine.
//
func ChanToWriter(source <-chan string, separator string, w io.Writer) (err error) {
	for s := range source {
		for _, err = io.WriteString(w, s); err != nil; {
			return
		}
		if len(separator) > 0 {
			for _, err = io.WriteString(w, separator); err != nil; {
				return
			}
		}
		if len(source) == 0 {
			// we don't know when the next string will come in, so yield
			runtime.Gosched()
		}
	}
	return
}

// Writes all strings from a channel to a bufio.Writer until the channel is closed. When the channel is empty, the
// buffer is flushed to the underlying io.Writer and the execution thread is yielded to another goroutine.
//
func ChanToBufioWriter(source <-chan string, separator string, bw *bufio.Writer) (err error) {
	for s := range source {
		for _, err = bw.WriteString(s); err != nil; {
			return
		}
		if len(separator) > 0 {
			for _, err = bw.WriteString(separator); err != nil; {
				return
			}
		}
		if len(source) == 0 {
			// we don't know when the next string will come in, so flush and yield
			for err = bw.Flush(); err != nil; {
				return
			}
			runtime.Gosched()
		}
	}
	err = bw.Flush()
	return
}
