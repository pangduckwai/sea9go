package match

// Usage:
// $ cd .../match
// $ go test -bench . -benchmem

import (
	"fmt"
	"testing"

	"github.com/pangduckwai/sea9go/pkg/errors"
)

var aLGORITHMS = []string{
	"AES-128-GCM",
	"AES-192-GCM",
	"AES-256-GCM",
	"AES-256-CBC",
	"ChaCha20-Poly1305",
	"RSA-2048-OAEP-SHA256",
	"RSA-2048-OAEP-SHA512",
	"RSA-4096-OAEP-SHA512",
	"RSA-2048-PKCS1v15",
	"ECIES-SECP256K1-DECRED",
	"ECIES-SECP256K1-ECIESGO",
}

var tESTS = []struct {
	s string
	x int
}{
	{"AES-128-GCM", 1},
	{"AES-192-GCM", 1},
	{"AES-256-GCM", 1}, {"a256gcm", 1}, {"AES256-GCM", 1},
	{"AES-256-CBC", 1}, {"a256cbc", 1}, {"AES-256-CBC-HS512", 1},
	{"AES-256", 2},
	{"CHACHA20-POLY1305", 1}, {"chapoly", 1},
	{"RSA-2048-OAEP-SHA256", 1}, {"RSA-OAEP-256", 1}, {"rsa256", 1},
	{"RSA-2048-OAEP-SHA512", 1},
	{"RSA-4096-OAEP-SHA512", 1},
	{"rsa-oaep", 3}, {"rsa512", 2},
	{"RSA-2048-PKCS1V15", 1}, {"RSA-PKCS1v15", 1},
	{"ECIES-SECP256K1-DECRED", 1}, {"decred", 1},
	{"ECIES-SECP256K1-ECIESGO", 1}, {"iesgo", 1},
	{"SECP256K1", 2}, {"ecies", 2},
	{"AES-192-CBC-HS512", 0},
	{"A128CBC-HS256", 0},
	{"3DES-64-GCM", 0},
	{"abcde-def", 0},
	{"abc3de-def", 0},
}

func TestMatch(t *testing.T) {
	var err error
	for i, inp := range tESTS {
		indices, mth, typ := Best(inp.s, aLGORITHMS, true)
		switch len(indices) {
		case 0:
			fmt.Printf("TestMatch() %2v x false %-23v (-) : match not found\n", i, inp.s)
			if inp.x != 0 {
				err = errors.Appendf(err, "expects %v match(es) for '%v', found none", inp.x, inp.s)
			}
		case 1:
			fmt.Printf("TestMatch() %2v v true  %-23v (%v) -> %v\n", i, inp.s, typ, mth)
			if inp.x != 1 {
				err = errors.Appendf(err, "expects %v match(es) for '%v', found 1", inp.x, inp.s)
			}
		default:
			mths := make([]string, 0)
			for _, idx := range indices {
				mths = append(mths, aLGORITHMS[idx])
			}
			fmt.Printf("TestMatch() %2v - true  %-23v (%v) -> %v (%v)\n", i, inp.s, typ, mths, mth)
			if inp.x != len(indices) {
				err = errors.Appendf(err, "expects %v match(es) for '%v', found %v", inp.x, inp.s, len(indices))
			}
		}
	}
	if errors.Count(err) > 0 {
		t.Fatal(err)
	}
}

func BenchmarkMatch(b *testing.B) {
	b.ResetTimer()
	i := 0
	for b.Loop() {
		Best(tESTS[i%31].s, aLGORITHMS, true)
		i++
	}
}
