package io

import (
	"fmt"
	"os"
	"testing"
)

// const BUFFER = 64 //1048576

// type encoding interface {
// 	Padding([]byte) []byte
// 	Multiple() (int, int)
// 	EncodeToString([]byte) string
// 	Encode(io.Reader, io.Writer) error
// 	DecodeString(string) ([]byte, error)
// 	Decode(io.Reader, io.Writer) error
// }

// func padding(inp []byte) (out []byte, err error) {
// 	ln := len(inp)
// 	out = make([]byte, 0)
// 	switch ln % 4 {
// 	case 2:
// 		out = append(inp, '=', '=')
// 	case 3:
// 		out = append(inp, '=')
// 	case 1:
// 		err = fmt.Errorf("invalid input \"%s\", %v %% 4 = 1", inp, len(inp))
// 		return
// 	default:
// 		out = inp
// 	}
// 	return
// }

// func encode(n encoding, r io.Reader, w io.Writer) (err error) {
// 	rdr, ok := r.(*bufio.Reader)
// 	if !ok {
// 		rdr = bufio.NewReaderSize(r, BUFFER)
// 	}

// 	wtr, ok := w.(*bufio.Writer)
// 	if !ok {
// 		wtr = bufio.NewWriter(w)
// 	}

// 	size := rdr.Size()
// 	lgh := 0
// 	dat := make([]byte, 0, size*2)
// 	enc, _ := n.Multiple()

// 	encode := func(inp []byte, ln int, flush bool) (err error) {
// 		if ln > 0 {
// 			encoded := n.EncodeToString(inp)
// 			_, err = fmt.Fprint(wtr, encoded)
// 			if err != nil {
// 				return
// 			}
// 		}
// 		if flush {
// 			err = wtr.Flush()
// 		}
// 		return
// 	}

// 	err = BufferedRead(rdr, size, func(cnt int, inp []byte) (err error) {
// 		if enc > 1 {
// 			lgh += cnt
// 			dat = append(dat, inp...)

// 			lgh -= lgh % enc // num of characters to encode each time is multiple of 3
// 			err = encode(dat[:lgh], lgh, false)

// 			if len(dat) > lgh {
// 				dat = dat[lgh:]
// 				lgh = len(dat)
// 			} else {
// 				dat = dat[:0]
// 				lgh = 0
// 			}
// 		} else {
// 			err = encode(inp[:cnt], cnt, false)
// 		}
// 		if err == nil {
// 			fmt.Printf("TEMP... %v (%v) bytes read\n", cnt, len(inp))
// 		}
// 		return err
// 	})
// 	if err != nil {
// 		return
// 	}

// 	err = encode(dat, lgh, true)
// 	return
// }

// func decode(n encoding, r io.Reader, w io.Writer) (err error) {
// 	rdr, ok := r.(*bufio.Reader)
// 	if !ok {
// 		rdr = bufio.NewReaderSize(r, BUFFER)
// 	}

// 	wtr, ok := w.(*bufio.Writer)
// 	if !ok {
// 		wtr = bufio.NewWriter(w)
// 	}

// 	size := rdr.Size()
// 	lgh := 0
// 	dat := make([]byte, 0, size*2)
// 	_, dec := n.Multiple()

// 	decode := func(inp []byte, ln int, flush bool) (err error) {
// 		if ln > 0 {
// 			inp = n.Padding(inp)
// 			var decoded []byte
// 			decoded, err = n.DecodeString(string(inp))
// 			if err != nil {
// 				return
// 			}
// 			_, err = wtr.Write(decoded)
// 			if err != nil {
// 				return
// 			}
// 		}
// 		if flush {
// 			err = wtr.Flush()
// 		}
// 		return
// 	}

// 	err = BufferedRead(rdr, size, func(cnt int, inp []byte) (err error) {
// 		if dec > 1 {
// 			lgh += cnt
// 			dat = append(dat, inp...)

// 			lgh -= lgh % dec // num of characters to decode each time is multiple of 4
// 			err = decode(dat[:lgh], lgh, false)

// 			if len(dat) > lgh {
// 				dat = dat[lgh:]
// 				lgh = len(dat)
// 			} else {
// 				dat = dat[:0]
// 				lgh = 0
// 			}
// 		} else {
// 			err = decode(inp[:cnt], cnt, false)
// 		}
// 		if err == nil {
// 			fmt.Printf("TEMP... %v (%v) bytes read\n", cnt, len(inp))
// 		}
// 		return err
// 	})
// 	if err != nil {
// 		return
// 	}

// 	err = decode(dat, lgh, true)
// 	return
// }

// type encdHex int

// func (n encdHex) Padding(inp []byte) []byte {
// 	return inp
// }
// func (n encdHex) Multiple() (int, int) {
// 	return 1, 1
// }
// func (n encdHex) EncodeToString(inp []byte) string {
// 	return hex.EncodeToString(inp)
// }
// func (n encdHex) Encode(in io.Reader, out io.Writer) error {
// 	return encode(n, in, out)
// }
// func (n encdHex) DecodeString(inp string) (out []byte, err error) {
// 	out, err = hex.DecodeString(inp)
// 	return
// }
// func (n encdHex) Decode(in io.Reader, out io.Writer) error {
// 	return decode(n, in, out)
// }

// type encdBase64 int

// func (n encdBase64) Padding(inp []byte) []byte {
// 	out, err := padding(inp)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return out
// }
// func (n encdBase64) Multiple() (int, int) {
// 	return 3, 4
// }
// func (n encdBase64) EncodeToString(inp []byte) string {
// 	return base64.StdEncoding.EncodeToString(inp)
// }
// func (n encdBase64) Encode(in io.Reader, out io.Writer) error {
// 	return encode(n, in, out)
// }
// func (n encdBase64) DecodeString(inp string) (out []byte, err error) {
// 	out, err = base64.StdEncoding.DecodeString(inp)
// 	return
// }
// func (n encdBase64) Decode(in io.Reader, out io.Writer) error {
// 	return decode(n, in, out)
// }

// func TestRead(t *testing.T) {
// 	raw := "NDg2NTZjNmM2ZjQ4NmY3NzQxNzI2NTU5NmY3NTNmNDkyNzZkNDY2OTZlNjU1NDY4NjE2ZTZiNTk2Zjc1NTY2NTcyNzk0ZDc1NjM2ODIxNDc2ZjZmNjQ1NDZmNDg2NTYxNzI1NDY4NjE3NDJjNGM2NTc0Mjc3MzQ3NmY0NDcyNjk2ZTZiNTQ2NTYxNTM2ZjZkNjU0NDYxNzkyMTQzNjU3Mjc0NjE2OTZlNmM3OTJjNGM2Zjc2NjU1NDZmNDQ3MjY5NmU2YjU0NjU2MTU2NjU3Mjc5NGQ3NTYzNjgyZTUzNjU2NTU5NmY3NTQxNzI2Zjc1NmU2NDIxNTk3NTcwNTM2NTY1NTk2MTIx"
// 	ite, _ := base64.StdEncoding.DecodeString(raw)
// 	expected, _ := hex.DecodeString(string(ite))

// 	var e0, e1, e2, e3 encoding = nil, encdBase64(2), nil, encdHex(1)
// 	in := bytes.NewReader([]byte(raw))
// 	result, err := Read(in, BUFFER, e0, e1, e2, e3)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if string(result) != string(expected) {
// 		t.Fatalf("TestRead() result '%s' does not match the expectation", result)
// 	}
// 	fmt.Printf("TestRead() result matches '%s'\n", expected)
// }

// func TestWrite(t *testing.T) {
// 	raw := []byte("HelloHowAreYou?I'mFineThankYouVeryMuch!GoodToHearThat,Let'sGoDrinkTeaSomeDay!Certainly,LoveToDrinkTeaVeryMuch.SeeYouAround!YupSeeYa!")
// 	ite := hex.EncodeToString(raw)
// 	expected := base64.StdEncoding.EncodeToString([]byte(ite))

// 	var e0, e1, e2, e3 encoding = nil, encdHex(1), nil, encdBase64(2)
// 	var buf bytes.Buffer
// 	wtr := bufio.NewWriter(&buf)
// 	err := Write(wtr, raw, e0, e1, e2, e3)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	result := buf.Bytes()

// 	if string(result) != expected {
// 		t.Fatalf("TestWrite() result '%s' does not match the expectation", result)
// 	}
// 	fmt.Printf("TestWrite() result matches '%v'\n", expected)
// }

func TestLine(t *testing.T) {
	f, err := os.Open("../../README.md")
	if err != nil {
		return
	}
	defer f.Close()

	rcnt, lcnt, err := ReadLine(
		f, 64,
		func(s string) error {
			// fmt.Printf("TestLine() - READ:'%v'\n", s)
			return nil
		},
		func(i1, i2 int, b ...[]byte) {
			fmt.Printf("TestLine() - LOG :'%s'\n", b)
		},
	)
	fmt.Printf("TestLine() - R:%v L:%v\n", rcnt, lcnt)
}
