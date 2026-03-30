package line

import (
	"bufio"
	"fmt"
	"io"
	"slices"
)

// ReadLine read line from file
func ReadLine(
	in io.Reader,
	bufferSize int,
	action func(string) error,
	log func(int, int, int, []byte),
) (
	rcnt int,
	lcnt int,
	err error,
) {
	if in == nil {
		err = fmt.Errorf("reader not ready")
		return
	}

	linefeed := func(c byte) bool { return c == 10 } // ascii 10 is "\n"

	rdr := bufio.NewReaderSize(in, bufferSize)
	n := 0
	buf, line := make([]byte, 0, bufferSize), make([]byte, 0, bufferSize)
	for rcnt = 0; ; rcnt++ {
		n, err = rdr.Read(buf[:bufferSize])
		if n > 0 {
			i0, i1 := 0, slices.IndexFunc(buf[:n], linefeed)
			for i1 >= 0 {
				if log != nil {
					log(n, i0, i1, buf[0:n])
				}
				line = append(line, buf[i0:i0+i1]...)
				errr := action(string(line))
				if errr != nil {
					if errr == io.EOF {
						break
					}
					return 0, 0, errr
				}
				lcnt++
				line = line[:0]
				i0 = i0 + i1 + 1
				if i0 < n {
					i1 = slices.IndexFunc(buf[i0:n], linefeed)
				} else {
					i1 = -1
				}
			}
			if i0 < n {
				if log != nil {
					log(n, i0, n, buf[0:n])
				}
				line = append(line, buf[i0:n]...)
			}
		}
		if err != nil {
			if err == io.EOF {
				err = nil // done
				if len(line) > 0 {
					err = action(string(line))
					if err != nil {
						return
					}
					lcnt++
				}
			}
			break
		}
	}
	return
}
