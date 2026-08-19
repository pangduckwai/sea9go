package maplogger

import (
	"fmt"
	"testing"

	"gopkg.in/yaml.v2"
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

func TestLogging(t *testing.T) {
	buf := getInput()
	// fmt.Printf("%s\n", buf)

	inp := make(map[string]interface{})
	err := yaml.Unmarshal(buf, &inp)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%v\n", inp)
	str, err := Format(inp)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(str)
}
