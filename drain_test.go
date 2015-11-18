package drain

import (
	"bufio"
	"bytes"
	"strconv"
	"strings"
	"testing"
)

func TestReaderLinesToChan(t *testing.T) {
	r := strings.NewReader("a b c d e\nf g h i j\n")
	c := make(chan string, 2)
	ReaderToChan(r, bufio.ScanLines, c)
	if line := <-c; line != "a b c d e" {
		t.FailNow()
	}
	if line := <-c; line != "f g h i j" {
		t.FailNow()
	}
	if len(c) != 0 {
		t.FailNow()
	}
}

func TestReaderWordsToChan(t *testing.T) {
	r := strings.NewReader("a b c\nd e f \n g h  i j")
	c := make(chan string, 10)
	ReaderToChan(r, bufio.ScanWords, c)
	expectedWords := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	for i := 0; i < len(expectedWords); i++ {
		if word := <-c; word != expectedWords[i] {
			t.FailNow()
		}
	}
	if len(c) != 0 {
		t.FailNow()
	}
}

func TestChanToWriter(t *testing.T) {
	c := make(chan string, 10)
	go func() {
		for _, s := range []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"} {
			c <- s
		}
		close(c)
	}()
	w := new(bytes.Buffer)
	ChanToWriter(c, ".", w)
	if w.String() != "a.b.c.d.e.f.g.h.i.j." {
		t.FailNow()
	}
}

func TestChanToBufioWriter(t *testing.T) {
	c := make(chan string, 10)
	go func() {
		for _, s := range []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"} {
			c <- s
		}
		close(c)
	}()
	w := new(bytes.Buffer)
	bw := bufio.NewWriterSize(w, 10)
	ChanToBufioWriter(c, ".", bw)
	if w.String() != "a.b.c.d.e.f.g.h.i.j." {
		t.FailNow()
	}
}

func BenchmarkChanToWriter(b *testing.B) {
	c := make(chan string, 1024)
	go func() {
		for i := 0; i < b.N*1024; i++ {
			c <- strconv.Itoa(i)
		}
		close(c)
	}()
	w := new(bytes.Buffer)
	b.StartTimer()
	ChanToWriter(c, ".", w)
	b.StopTimer()
}

func BenchmarkChanToBufioWriter(b *testing.B) {
	c := make(chan string, 1024)
	go func() {
		for i := 0; i < b.N*1024; i++ {
			c <- strconv.Itoa(i)
		}
		close(c)
	}()
	w := new(bytes.Buffer)
	bw := bufio.NewWriterSize(w, 64*1024)
	b.StartTimer()
	ChanToBufioWriter(c, ".", bw)
	b.StopTimer()
}
