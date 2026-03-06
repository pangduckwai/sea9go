package traverse

import (
	"fmt"
	"slices"
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
		func(keys []string, in any) (out any, err error) {
			out = fmt.Sprintf("%v", in)
			if !slices.Contains(verify, out.(string)) {
				t.Fatalf("TestTraverse() %v %v not found", cnt, out)
			}
			cnt++

			var sb strings.Builder
			sb.WriteString(keys[0])
			for _, k := range keys[1:] {
				if k[0] != '[' {
					sb.WriteString(".")
				}
				sb.WriteString(k)
			}
			fmt.Printf("TestTraverse() %19v %v\n", out, sb.String())
			return
		},
	)
	if err != nil {
		t.Fatal(err)
	}
}
