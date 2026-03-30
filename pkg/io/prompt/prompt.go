package prompt

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

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
