package avibase_downloader

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (input *Input) memoizedFileName() string {
	suffix := ""
	if input.IncludeRare {
		suffix = "_rare"
	}
	return fmt.Sprintf("/memo/avibase_downloader/%v%s.xml", input.RegionCodes, suffix)
}

func (input *Input) readMemoized() (*Output, error) {
	output := &Output{}
	asBytes, err := ioutil.ReadFile(input.memoizedFileName())
	if err != nil {
		return output, err
	}
	err = xml.Unmarshal(asBytes, output)
	return output, err
}

func (input *Input) writeMemoized(output *Output) error {
	err := os.MkdirAll(filepath.Dir(input.memoizedFileName()), 0777)
	if err != nil {
		return err
	}
	asBytes, err := xml.MarshalIndent(*output, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(input.memoizedFileName(), asBytes, 0777)
}
