# iorate

A Go package that provides a reader and a writer that wrap `io.Reader`
or `io.Writer` objects limiting the byte rate for the read or write
operations.


## Examples

To limit an outbound data rate for a network connection to 33.6 Kbps:

	// Suppose we have a connection:
	ln, err := net.Dial("tcp", "example.net:31415")
	if err != nil {
		...
	}

	// Obtain a writer that will feed the data to the connection:
	w := iorate.NewWriter(ln, 33.6 * iorate.Kbps)

	...

To limit reading from a file to 5 MBps:

	// Suppose we have a file:
	f, err := os.Open("example.bin")
	if err != nil {
		...
	}

	// Obtain a rate-limited reader:
	r := iorate.NewReader(f, 5 * iorate.MBps)

	...

