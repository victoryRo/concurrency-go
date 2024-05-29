package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/victoryRo/concurrency-go/chapter8/streams/filter"
	"github.com/victoryRo/concurrency-go/chapter8/streams/store"
)

type Message struct {
	At    time.Time `json:"at"`
	Value float64   `json:"value"`
	Error string    `json:"err"`
}

func initDB(dbName string) {
	db, err := sql.Open("sqlite", dbName)
	panicErr(err)

	// Create the timeseries table
	tx, err := db.Begin()
	panicErr(err)

	_, err = tx.Exec(`create table if not exists measurements(at integer, value double)`)
	panicErr(err)

	// Fill it with some data
	result, err := tx.Query(`select count(*) from measurements`)
	panicErr(err)

	result.Next()
	var nItems int
	err = result.Scan(&nItems)
	panicErr(err)

	_ = tx.Commit()
	if nItems < 10000 {
		tx, err := db.Begin()
		panicErr(err)

		fmt.Printf("nRows: %d, inserting data...\n", nItems)
		tm := time.Now().UnixMilli()
		stmt, err := tx.Prepare(`insert into measurements(at,value) values(?,?)`)
		panicErr(err)

		for i := 0; i < 10000; i++ {
			_, err := stmt.Exec(tm, rand.Float64())
			panicErr(err)

			tm -= 100
		}
		_ = tx.Commit()
	}
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func DecodeToChan[T any](decode func(*T) error) (<-chan T, <-chan error) {
	ret := make(chan T)
	errch := make(chan error, 1)
	go func() {
		defer close(ret)
		defer close(errch)
		var entry T
		for {
			if err := decode(&entry); err != nil {
				if !errors.Is(err, io.EOF) {
					errch <- err
				}
				return
			}
			ret <- entry
		}
	}()
	return ret, errch
}

func EncodeFromChan[T any](input <-chan T, encode func(T) ([]byte, error), out io.Writer) <-chan error {
	ret := make(chan error, 1)
	go func() {
		defer close(ret)
		for entry := range input {
			data, err := encode(entry)
			if err != nil {
				ret <- err
				return
			}
			if _, err := out.Write(data); err != nil {
				if !errors.Is(err, io.EOF) {
					ret <- err
				}
				return
			}
		}
	}()
	return ret
}

func httpMain() {
	initDB("test.db")
	// Start database
	db, err := sql.Open("sqlite", "test.db")
	if err != nil {
		panic(err)
	}

	// Create the store
	st := store.Store{DB: db}

	// Create an HTTP server
	http.HandleFunc("/db", func(w http.ResponseWriter, req *http.Request) {
		data, err := st.Stream(req.Context(), store.Request{})
		if err != nil {
			fmt.Println("Store error", err)
		}

		errCh := EncodeFromChan(data, func(entry store.Entry) ([]byte, error) {
			msg := Message{
				At:    entry.At,
				Value: entry.Value,
			}
			if entry.Error != nil {
				msg.Error = entry.Error.Error()
			}
			return json.Marshal(msg)
		}, w)
		err = <-errCh
		if err != nil {
			fmt.Println("Encode error", err)
		}
	})
	go func() {
		fmt.Println("Server started at :10001")
		if err := http.ListenAndServe(":10001", nil); err != nil {
			panic(err)
		}
	}()

	// Create a client
	resp, err := http.Get("http://localhost:10001/db")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	entries, rcvErr := DecodeToChan[store.Entry](func(entry *store.Entry) error {
		var msg Message
		if err := decoder.Decode(&msg); err != nil {
			return err
		}
		entry.At = msg.At
		entry.Value = msg.Value
		if msg.Error != "" {
			entry.Error = fmt.Errorf(msg.Error)
		}

		return nil
	})

	filteredEntries := filter.MinFilter(0.001, entries)
	entryCh, errCh := filter.ErrFilter(filteredEntries)
	resultCh := filter.MovingAvg(0.5, 5, entryCh)

	go func() {
		for err := range errCh {
			fmt.Println("Stream error", err)
		}
	}()
	for entry := range resultCh {
		fmt.Printf("%+v\n", entry)
	}
	err = <-rcvErr
	if err != nil {
		fmt.Println("Receive error", err)
	}
}
