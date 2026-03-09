package io

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

type Encoder interface {
	Encode(io.Reader, io.Writer) error
}

type Decoder interface {
	Decode(io.Reader, io.Writer) error
}

func pipedEncode(in io.Reader, out io.Writer, encoders ...Encoder) (err error) {
	lgth := len(encoders)
	cs := make([]chan error, 0)
	rs := make([]io.Reader, 0)
	ws := make([]io.Writer, 0)
	for i := 1; i < lgth; i++ {
		r, w := io.Pipe()
		cs = append(cs, make(chan error))
		rs = append(rs, r)
		ws = append(ws, w)
	}
	last := len(cs) - 1

	for i := range encoders[1 : lgth-1] {
		go func() {
			var err error
			err = encoders[i+1].Encode(rs[i], ws[i+1])
			if err != nil {
				cs[i] <- err
			}
			cs[i] <- ws[i+1].(*io.PipeWriter).CloseWithError(nil)
		}()
	}
	go func() {
		cs[last] <- encoders[lgth-1].Encode(rs[last], out)
	}()

	for j, c := range cs[:last] {
		go func() {
			for e := range c {
				if e != nil {
					cs[last] <- fmt.Errorf("[PIPE][ENCODE] %v: %v", j, e)
				}
			}
		}()
	}

	err = encoders[0].Encode(in, ws[0])
	if err != nil {
		err = fmt.Errorf("[PIPE][ENCODE] %v", err)
		return
	}
	err = ws[0].(*io.PipeWriter).CloseWithError(nil)
	if err != nil {
		err = fmt.Errorf("[PIPE][ENCODE][CLOSE] %v", err)
		return
	}

	err = <-cs[last]
	return
}

func pipedDecode(in io.Reader, out io.Writer, decoders ...Decoder) (err error) {
	lgth := len(decoders)
	cs := make([]chan error, 0)
	rs := make([]io.Reader, 0)
	ws := make([]io.Writer, 0)
	for i := 1; i < lgth; i++ {
		r, w := io.Pipe()
		cs = append(cs, make(chan error))
		rs = append(rs, r)
		ws = append(ws, w)
	}
	last := len(cs) - 1

	for i := range decoders[1 : lgth-1] {
		go func() {
			var err error
			err = decoders[i+1].Decode(rs[i], ws[i+1])
			if err != nil {
				cs[i] <- err
			}
			cs[i] <- ws[i+1].(*io.PipeWriter).CloseWithError(nil)
		}()
	}
	go func() {
		cs[last] <- decoders[lgth-1].Decode(rs[last], out)
	}()

	for j, c := range cs[:last] {
		go func() {
			for e := range c {
				if e != nil {
					cs[last] <- fmt.Errorf("[PIPE][DECODE] %v: %v", j, e)
				}
			}
		}()
	}

	err = decoders[0].Decode(in, ws[0])
	if err != nil {
		err = fmt.Errorf("[PIPE][DECODE] %v", err)
		return
	}
	err = ws[0].(*io.PipeWriter).CloseWithError(nil)
	if err != nil {
		err = fmt.Errorf("[PIPE][DECODE][CLOSE] %v", err)
		return
	}

	err = <-cs[last]
	return
}

// BufferedRead read from the given reader the specified num of bytes at a time. After
// each read operation, the bytes read is passed to the given function to process.
// To support interactive input from stdin, an entire line with a single period ('.')
// character is used to signify the end of input.
func BufferedRead(
	rdr *bufio.Reader,
	size int,
	action func(int, []byte) error,
) (
	err error,
) {
	buf := make([]byte, 0, size)
	cnt := 0

	for err == nil {
		// As described in the doc, handle read data first if n > 0 before handling error,
		// it is because the returned error could have been EOF
		cnt, err = rdr.Read(buf[:cap(buf)])

		// If getting input from stdin interactively, pressing <enter> would signify the end of an input line.
		// An entire line with a signle period ('.') means the end of input.
		if (cnt == 2 && buf[:cnt][1] == 10) || (cnt == 3 && (buf[:cnt][1] == 13 && buf[:cnt][2] == 10 || buf[:cnt][1] == 10 && buf[:cnt][2] == 13)) {
			// ASCII code 10: line feed (LF)
			// ASCII code 13: carriage return (CR)
			// ASCII code 46: period ('.')
			if buf[:cnt][0] == 46 {
				cnt = 0
				err = io.EOF
			}
		}

		if cnt > 0 {
			err = action(cnt, buf[:cnt])
			if err != nil {
				break
			}
		}
	}

	if err == io.EOF {
		err = nil
	}
	return
}

// Read read from the file of the given path. An empty path means to read from stdin. A chain of
// decoders can be applied for encoded inputs.
func Read(
	in io.Reader,
	bufferSize int,
	decoders ...Decoder,
) (
	data []byte,
	err error,
) {
	if in == nil {
		err = fmt.Errorf("reader not ready")
		return
	}
	rdr := bufio.NewReaderSize(in, bufferSize)

	dec := make([]Decoder, 0)
	for _, n := range decoders { // filter out nil decoder inputs
		if n != nil {
			dec = append(dec, n)
		}
	}
	if len(dec) <= 0 {
		data = make([]byte, 0, bufferSize*2)
		err = BufferedRead(rdr, bufferSize, func(cnt int, buf []byte) error {
			data = append(data, buf...)
			return nil
		})
	} else {
		var buf bytes.Buffer
		wtr := bufio.NewWriter(&buf)
		if len(dec) <= 1 {
			err = dec[0].Decode(rdr, wtr)
			if err != nil {
				return
			}
		} else {
			err = pipedDecode(rdr, wtr, dec...)
			if err != nil {
				return
			}
		}
		err = wtr.Flush()
		if err != nil {
			return
		}
		data = buf.Bytes()
	}
	return
}

// Write write to the file of the given path. An empty path means to write to stdout. A chain of
// encoders can be applied to encode the data before writing to file.
func Write(
	out io.Writer,
	data []byte,
	encoders ...Encoder,
) (err error) {
	if out == nil {
		err = fmt.Errorf("writer not ready")
		return
	}
	wtr := bufio.NewWriter(out)

	enc := make([]Encoder, 0)
	for _, n := range encoders { // filter out nil encoder inputs
		if n != nil {
			enc = append(enc, n)
		}
	}
	if len(enc) <= 0 {
		fmt.Fprintf(wtr, "%s", data)
		wtr.Flush()
	} else {
		rdr := bytes.NewReader(data)
		if wtr == nil {
			wtr = bufio.NewWriter(os.Stdout)
		}
		if len(enc) <= 1 {
			err = enc[0].Encode(rdr, wtr)
			if err != nil {
				return
			}
		} else {
			err = pipedEncode(rdr, wtr, enc...)
			if err != nil {
				return
			}
		}
		err = wtr.Flush()
	}
	return
}

// Prompt get a single line input interactively from the prompt.
func Prompt(header, prompt string) (str string, err error) {
	rdr := bufio.NewReader(os.Stdin)
	if header != "" {
		fmt.Printf("%v:\n", header)
	}
	fmt.Printf("%v: ", prompt)
	str, err = rdr.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			err = fmt.Errorf("[IACTS] stdin already ended, cannot read input")
		} else {
			err = fmt.Errorf("[IACTS] %v", err)
		}
	}

	buf := []byte(str)
	lgh := len(buf)
	if lgh > 1 && (buf[lgh-2] == 13 && buf[lgh-1] == 10 || buf[lgh-2] == 10 && buf[lgh-1] == 13) {
		str = str[:lgh-2]
	} else if lgh > 0 && buf[lgh-1] == 10 {
		str = str[:lgh-1]
	}

	return
}
