package pipeline

import (
	"encoding/csv"
	"fmt"
	"io"
)

func synchronousPipeline(input *csv.Reader) {
	fmt.Println("---- Synchronous Pipeline -----")

	_, _ = input.Read() // ignore first row
	for {
		rec, err := input.Read()
		if err == io.EOF {
			return
		}
		if err != nil {
			panic(err)
		}
		out := encode(convert(parse(rec)))
		fmt.Println(string(out))
	}
}
