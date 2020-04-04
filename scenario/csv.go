package scenario

import (
	"context"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"log"
	"os"
	"strconv"
	"time"
)

type CSV struct {
	file    *os.File
	encoder vegeta.Encoder
}

func EstablishCSV(name string, now time.Time) *CSV {
	file, err := os.Create("nft-results-" + name + "-" + strconv.FormatInt(now.Unix(), 10) + ".csv")
	if err != nil {
		log.Fatal(err)
	}
	csvEncoder := vegeta.NewCSVEncoder(file)
	return &CSV{
		file:    file,
		encoder: csvEncoder,
	}
}

func (C *CSV) SendResults(_ context.Context, r vegeta.Result) error {
	return C.encoder.Encode(&r)
}

func (C *CSV) Close(context.Context) error {
	return C.file.Close()
}
