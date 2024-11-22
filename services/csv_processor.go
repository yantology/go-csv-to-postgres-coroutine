package services

import (
	"context"
	"encoding/csv"
	"io"
	"log"
	"mime/multipart"
	"sync"

	"github.com/yantology/go-csv-to-postgres-coroutine/config"
)

const (
	batchSize   = 1000
	workerCount = 100
)

type CSVProcessor struct {
	db       *config.Database
	wg       sync.WaitGroup
	jobsChan chan []interface{}
}

func NewCSVProcessor(db *config.Database) *CSVProcessor {
	return &CSVProcessor{
		db:       db,
		jobsChan: make(chan []interface{}, batchSize),
	}
}

func ProcessCSVFile(db *config.Database, file *multipart.FileHeader) error {
	processor := NewCSVProcessor(db)
	return processor.ProcessFile(file)
}

func (p *CSVProcessor) ProcessFile(file *multipart.FileHeader) error {
	f, err := file.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	reader := csv.NewReader(f)

	// Start workers
	for i := 0; i < workerCount; i++ {
		go p.worker(i)
	}

	// Read and process CSV
	isHeader := true
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if isHeader {
			isHeader = false
			continue
		}

		data := make([]interface{}, len(record))
		for i, v := range record {
			data[i] = v
		}

		p.wg.Add(1)
		p.jobsChan <- data
	}

	close(p.jobsChan)
	p.wg.Wait()

	return nil
}

func (p *CSVProcessor) worker(id int) {
	counter := 0
	for data := range p.jobsChan {
		err := p.insertRecord(data)
		if err != nil {
			log.Printf("Worker %d error: %v", id, err)
		}
		log.Printf("Worker %d processed %d records succses", id, counter)
		counter++
		p.wg.Done()
	}
}

func (p *CSVProcessor) insertRecord(values []interface{}) error {
	query := `
		INSERT INTO domain (
			GlobalRank, TldRank, Domain, TLD, 
			RefSubNets, RefIPs, IDN_Domain, IDN_TLD,
			PrevGlobalRank, PrevTldRank, PrevRefSubNets, PrevRefIPs
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	ctx := context.Background()
	conn, err := p.db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(ctx, query, values...)
	return err
}
