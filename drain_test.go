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

	expectedLines := []string{"a b c d e", "f g h i j"}
	for i, expected := range expectedLines {
		if actual := <-c; expected != actual {
			t.Errorf("Expected '%s' at index %d, got '%s'", expected, i, actual)
		}
	}
	if l := len(c); l != 0 {
		t.Errorf("Expected empty channel, got length %d", l)
	}
}

func TestReaderWordsToChan(t *testing.T) {
	r := strings.NewReader("a b c\nd e f \n g h  i j")
	c := make(chan string, 10)

	ReaderToChan(r, bufio.ScanWords, c)

	expectedWords := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	for i, expected := range expectedWords {
		if actual := <-c; expected != actual {
			t.Errorf("Expected '%s' at index %d, got '%s'", expected, i, actual)
		}
	}
	if l := len(c); l != 0 {
		t.Errorf("Expected empty channel, got length %d", l)
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

	expected := "a.b.c.d.e.f.g.h.i.j."
	if actual := w.String(); expected != actual {
		t.Errorf("Expected '%s', got '%s'", expected, actual)
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

	expected := "a.b.c.d.e.f.g.h.i.j."
	if actual := w.String(); expected != actual {
		t.Errorf("Expected '%s', got '%s'", expected, actual)
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
