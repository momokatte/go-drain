package drain

func ExampleChanToStdout() {
	c := make(chan string, 2)
	go func() {
		for _, s := range []string{"abcde", "fghij", "klmno", "pqrst"} {
			c <- s
		}
		close(c)
	}()

	ChanToStdout(c)

	// Output:
	// abcde
	// fghij
	// klmno
	// pqrst
}
