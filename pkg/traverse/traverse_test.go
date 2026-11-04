package traverse

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func getInput() []byte {
	return []byte(
		"header:\n" +
			"  version: 1.1.3\n" +
			"  type: runner\n" +
			"content:\n" +
			"  logLevel: info\n" +
			"  logFile: default\n" +
			"  concurrent: 10\n" +
			"  relogin: true\n" +
			"  retry: 8\n" +
			"  backoff: 2\n" +
			"  wait: 5\n" +
			"  list:\n" +
			"  - id: 1\n" +
			"    name: Search\n" +
			"    path: dosrch.yaml\n" +
			"    body:\n" +
			"      query:\n" +
			"        status:\n" +
			"          $eq: Valid\n" +
			"  - id: 2\n" +
			"    name: Another test\n" +
			"    path: here.json\n" +
			"  - id: 3\n" +
			"    name: Last test\n" +
			"    path: there.json\n" +
			"    body:\n" +
			"      query:\n" +
			"        value:\n" +
			"        - How are you today\n" +
			"        - I'm fine thank you\n" +
			"        - Good to know bye\n" +
			"  envVar:\n" +
			"    LOG_LEVEL: logLevel\n" +
			"    THREAD_NUM: concurrent\n" +
			"    TLS_SERVER: serverCert\n" +
			"  envList:\n" +
			"  - retries:\n" +
			"    - retry1\n" +
			"    - retry2\n" +
			"    - retry3\n" +
			"  - relogin\n" +
			"  - - tester1\n" +
			"    - tester2\n" +
			"    - tester3\n" +
			"  - verifier",
	)
}

var verify []string = []string{
	"1.1.3", "runner", "info", "default", "10", "true", "8", "2", "5", "1", "Search", "dosrch.yaml", "Valid", "2", "Another test", "here.json", "3", "Last test", "there.json", "How are you today", "I'm fine thank you",
	"Good to know bye", "logLevel", "concurrent", "serverCert", "retry1", "retry2", "retry3", "relogin", "tester1", "tester2", "tester3", "verifier",
}

func TestTraverse(t *testing.T) {
	inp := make(map[string]interface{})
	err := yaml.Unmarshal(getInput(), &inp)
	if err != nil {
		t.Fatal(err)
	}

	cnt := 0
	_, err = Traverse(
		inp,
		func(typ NodeType, keys []string, in any) (out any, err error) {
			var sb strings.Builder
			sb.WriteString(keys[0])
			for _, k := range keys[1:] {
				if k[0] != '[' {
					sb.WriteString(".")
				}
				sb.WriteString(k)
			}

			if typ != LEAF {
				t := "..OBJT.."
				if typ == LIST {
					t = "..LIST.."
				}
				fmt.Printf("TestTraverse() %-19v %v\n", t, sb.String())
				return
			}

			out = fmt.Sprintf("%v", in)
			if !slices.Contains(verify, out.(string)) {
				t.Fatalf("TestTraverse() %v %v not found", cnt, out)
			}
			cnt++

			fmt.Printf("TestTraverse() %19v %v\n", out, sb.String())
			return
		},
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestError(t *testing.T) {
	err1 := ShouldStop
	err2 := errors.New("file does not exist")

	if err1 != ShouldStop {
		t.Fatal("should equal...")
	}

	if err2 == ShouldStop {
		t.Fatal("should not equal...")
	}

	fmt.Println("TestError() test passed")
}

func ctrlParseKey(key string) int {
	ln := len(key)
	if ln < 3 {
		return -1
	}
	if key[0] != '[' || key[ln-1] != ']' || key[1] == '-' {
		return -1
	}
	idx, err := strconv.Atoi(key[1 : ln-1])
	if err != nil {
		return -1
	}
	return idx
}

var kEYS = []string{"hello", "[0]", "[9]", "137", "[-46]", "[]", "[90]", "[there]", "[2048]", "[3]", "[7]", "41"}

func TestParseKey(t *testing.T) {
	for i := 0; i < 11; i++ {
		idx := ParseKey(kEYS[i])
		switch i {
		case 0:
			fallthrough
		case 3:
			fallthrough
		case 4:
			fallthrough
		case 11:
			if idx != -2 {
				t.Fatalf("'%v' is not a valid index with value %v", kEYS[i], idx)
			}
		case 5:
			if idx != -3 {
				t.Fatalf("'%v' is not a valid index with value %v", kEYS[i], idx)
			}
		case 7:
			if idx != -1 {
				t.Fatalf("'%v' is not a valid index with value %v", kEYS[i], idx)
			}
		case 1:
			if idx != 0 {
				t.Fatalf("'%v' should be a valid index instead of %v", kEYS[i], idx)
			}
		case 2:
			if idx != 9 {
				t.Fatalf("'%v' should be a valid index instead of %v", kEYS[i], idx)
			}
		case 6:
			if idx != 90 {
				t.Fatalf("'%v' should be a valid index instead of %v", kEYS[i], idx)
			}
		case 8:
			if idx != 2048 {
				t.Fatalf("'%v' should be a valid index instead of %v", kEYS[i], idx)
			}
		case 9:
			if idx != 3 {
				t.Fatalf("'%v' should be a valid index instead of %v", kEYS[i], idx)
			}
		case 10:
			if idx != 7 {
				t.Fatalf("'%v' should be a valid index instead of %v", kEYS[i], idx)
			}
		}
	}
	fmt.Println("TestParseKey() test passed")
}

func BenchmarkParseKey(b *testing.B) {
	b.ResetTimer()
	i := 0
	for b.Loop() {
		ParseKey(kEYS[i%11])
		i++
	}
}

func BenchmarkCtrlParseKey(b *testing.B) {
	b.ResetTimer()
	i := 0
	for b.Loop() {
		ctrlParseKey(kEYS[i%11])
		i++
	}
}
