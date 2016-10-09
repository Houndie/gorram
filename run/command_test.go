package run

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// func Now() Time
// tests zero arg Function.
// Tests printing of value with ToString method.
func TestTimeNow(t *testing.T) {
	t.Parallel()
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(dir)
	stderr := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	env := Env{
		Stderr: stderr,
		Stdout: stdout,
	}
	err = Run(Command{Package: "time", Function: "Now", Cache: dir}, env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := stdout.String()
	expected := fmt.Sprint(time.Now()) + "\n"
	if !strings.HasPrefix(out, expected[:15]) {
		t.Fatalf("Expected ~%q but got %q", expected, out)
	}
	if !strings.HasSuffix(out, expected[len(expected)-9:]) {
		t.Fatalf("Expected ~%q but got %q", expected, out)
	}
}

// func Sqrt(x float64) float64
// Tests float parsing arguments and outputs.
func TestMathSqrt(t *testing.T) {
	t.Parallel()
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(dir)
	stderr := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	env := Env{
		Stderr: stderr,
		Stdout: stdout,
	}
	err = Run(Command{
		Package:  "math",
		Function: "Sqrt",
		Args:     []string{"25.4"},
		Cache:    dir},
		env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := stdout.String()
	expected := "5.039841267341661\n"
	if out != expected {
		t.Fatalf("Expected %q but got %q", expected, out)
	}
}

// func Indent(dst *bytes.Buffer, src []byte, prefix, indent string) error
// Tests stdin to []byte argument.
// Tests a dst *bytes.Buffer with a []byte src.
// Tests string arguments.
func TestJsonIndentStdin(t *testing.T) {
	t.Parallel()
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(dir)
	stderr := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	stdin := strings.NewReader(`{ "foo" : "bar" }`)
	env := Env{
		Stderr: stderr,
		Stdout: stdout,
		Stdin:  stdin,
	}
	err = Run(Command{
		Package:  "encoding/json",
		Function: "Indent",
		Args:     []string{"", "  "},
		Cache:    dir},
		env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := stdout.String()
	expected := `
{
  "foo": "bar"
}
`[1:]
	if out != expected {
		t.Fatalf("Expected %q but got %q", expected, out)
	}
}

// func Get(url string) (resp *Response, err error)
// Tests a single string argument.
// Tests val, err return value.
// Tests struct return value that contains an io.Reader.
func TestNetHTTPGet(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(dir)
	stderr := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	env := Env{
		Stderr: stderr,
		Stdout: stdout,
	}
	err = Run(Command{
		Package:  "net/http",
		Function: "Get",
		Args:     []string{ts.URL},
		Cache:    dir},
		env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := stdout.String()
	expected := "Hello, client\n\n"
	if out != expected {
		t.Fatalf("Expected %q but got %q", expected, out)
	}
}

// func (enc *Encoding) EncodeToString(src []byte) string
// Tests calling a method on a global variable.
// Tests passing a filename as a []byte argument.
func TestBase64EncodeToStringFromFilename(t *testing.T) {
	t.Parallel()
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(dir)
	filename := filepath.Join(dir, "out.txt")
	if err := ioutil.WriteFile(filename, []byte("12345"), 0600); err != nil {
		t.Fatal(err)
	}

	stderr := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	env := Env{
		Stderr: stderr,
		Stdout: stdout,
	}
	err = Run(Command{
		Package:   "encoding/base64",
		GlobalVar: "StdEncoding",
		Function:  "EncodeToString",
		Args:      []string{filename},
		Cache:     dir},
		env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := stdout.String()
	expected := "MTIzNDU=\n"
	if out != expected {
		t.Fatalf("Expected %q but got %q", expected, out)
	}
}
