// Copyright 2015 Michael O'Rourke. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package drain

import (
	"bufio"
	"os"
)

// Scans all lines from os.Stdin into a channel until EOF
//
func StdinLinesToChan(dest chan<- string) error {
	return StdinToChan(bufio.ScanLines, dest)
}

// Scans all words from os.Stdin into a channel until EOF
//
func StdinWordsToChan(dest chan<- string) error {
	return StdinToChan(bufio.ScanWords, dest)
}

// Scans all tokens from os.Stdin into a channel until EOF
//
func StdinToChan(split bufio.SplitFunc, dest chan<- string) error {
	return ReaderToChan(os.Stdin, split, dest)
}

// Writes all strings from a channel to os.Stdout, with trailing newline characters, until the channel is closed.
// Uses the default buffer size.
//
func ChanToStdout(source <-chan string) error {
	return ChanToStdoutSize(source, DefaultBufferSize)
}

// Writes all strings from a channel to os.Stdout, with trailing newline characters, until the channel is closed.
//
func ChanToStdoutSize(source <-chan string, bufSize int) error {
	bw := bufio.NewWriterSize(os.Stdout, bufSize)
	return ChanToBufioWriter(source, NewLineString, bw)
}
