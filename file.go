// Copyright 2015 Michael O'Rourke. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package drain

import (
	"bufio"
	"os"
)

// Scans all lines from a File into a channel until EOF
//
func FileLinesToChan(fileName string, dest chan<- string) error {
	return FileToChan(fileName, bufio.ScanLines, dest)
}

// Scans all tokens from a File into a channel until EOF
//
func FileToChan(fileName string, split bufio.SplitFunc, dest chan<- string) (err error) {
	f, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer f.Close()
	err = ReaderToChan(f, split, dest)
	return
}

// Writes all strings from a channel to a File, with trailing newline characters, until the channel is closed.
// Uses the default buffer size.
//
func ChanToFile(source <-chan string, fileName string) error {
	return ChanToFileSize(source, DefaultBufferSize, fileName)
}

// Writes all strings from a channel to a File, with trailing newline characters, until the channel is closed.
// The file will be overwritten if it already exists, or created if it does not.
//
func ChanToFileSize(source <-chan string, bufSize int, fileName string) (err error) {
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return
	}
	defer f.Close()
	bw := bufio.NewWriterSize(f, bufSize)
	err = ChanToBufioWriter(source, NewLineString, bw)
	return
}
