package main

import (
	"fmt"
	"log"

	"github.com/gbdubs/avibase_downloader"
)

func main() {
	input := avibase_downloader.AvibaseDownloaderInput{
		RegionCode: "USfl",
	}
	output, err := input.Execute()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", output)
}
