package svach_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/nilsocket/svach"
)

var (
	blnsDir = "blns"

	assetDir = "assets"

	blnsFile = filepath.Join(assetDir, "blns.json")

	fileExistsCount = 0
)

var (
	input, output []string
)

var s1 = svach.DefaultSvach
var s2, _ = svach.WithOpts("-", 235)

func init() {
	input = fileData(blnsFile)
}

func TestAll(t *testing.T) {
	os.RemoveAll(blnsDir)

	testSvach(t, s1, "s1")
	testSvach(t, s2, "s2")
}

func testSvach(t *testing.T, s *svach.Svach, prefix string) {

	outputDir := filepath.Join(blnsDir, prefix+"name")
	expectedOutputFile := filepath.Join(assetDir, prefix+"ExpectedName.json")

	test(t, s.Name, outputDir, expectedOutputFile)

	outputDir = filepath.Join(blnsDir, prefix+"clean")
	expectedOutputFile = filepath.Join(assetDir, prefix+"ExpectedClean.json")

	test(t, s.Clean, outputDir, expectedOutputFile)
}

func test(t *testing.T, fn func(string) string, outputDir, expectedOutputFile string) {
	t.Log("\n\n" + expectedOutputFile + "\n")
	output = testCommon(t, outputDir, fn)
	// writeToFile(expectedOutputFile, output)
	testValid(t, expectedOutputFile, output)
}

func testValid(t *testing.T, expectedOutputFile string, output []string) {
	equal(t, fileData(expectedOutputFile), output)
}

func fileData(fileName string) []string {
	data, _ := ioutil.ReadFile(fileName)
	dataSlice := make([]string, 0, 550)
	json.Unmarshal(data, &dataSlice)
	return dataSlice
}

func equal(t *testing.T, expected, result []string) {

	if (expected == nil) != (result == nil) {
		t.Error("expected:", expected, "result:", result)
		return
	}

	if len(expected) != len(result) {
		t.Error("length doesn't match")
		return
	}

	for i := range expected {
		if expected[i] != result[i] {
			t.Error(i+1, "expected:", expected[i], "got:", result[i])
			t.Errorf("\n%+q\n%+q\n%+q\n", []byte(input[i]), []byte(expected[i]), []byte(result[i]))
			return
		}
	}

}

func writeToFile(fileName string, data []string) {
	encData, err := json.MarshalIndent(data, "", "")
	if err != nil {
		log.Println(err)
	}
	ioutil.WriteFile(fileName, encData, os.ModePerm)
}

func testCommon(t *testing.T, dir string, fn func(string) string) []string {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatalln(err)
	}

	output := make([]string, 0, 550)

	for i, istr := range input {
		res := fn(istr)
		output = append(output, res)
		createFile(t, i, dir, res)
	}
	return output
}

func createFile(t *testing.T, i int, dir, name string) {
	file := filepath.Join(dir, name)

	if fileExists(file) {
		t.Log(i, ":", file, "already exists")
		fileExistsCount++
		return
	}

	f, err := os.Create(file)
	if err != nil {
		t.Error(i+2, file, err)
	} else {
		f.Close()
	}
}

func fileExists(name string) bool {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		// if file doesn't exist return false
		return false
	}
	return true
}

// Examples

func ExampleClean() {
	res := svach.Clean(`.....Hello<>:/---\....W|orld?.!..`)
	fmt.Println(res)
	// Output: .Hello-.World.!
}

func ExampleWithOpts() {
	s, _ := svach.WithOpts(" ", 6)
	res := s.Name(`.....H<>e:l.!..`)
	fmt.Println(res)
	// Output: .H e l
}

func ExampleWithOpts_error() {
	_, err := svach.WithOpts("?", 6)
	fmt.Println(err)
	// Output: Invalid characters like `., <, >, :, ", /, \, |, ?, *` exist in replaceStr
}

func ExampleName() {
	// Incase of invalid filename, md5sum of filename is returned.
	res := svach.Name(`<>:"/\|?*`)
	fmt.Println(res)
	// Output: 3e4bde3cb1e4c9cfa2db74bbc536d5e2
}
