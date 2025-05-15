package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type TaskParams struct {
	Args           []byte
	IdempotencyKey pgtype.Text
}

// DataGenerator is a structure that emits data continuously
type DataGenerator struct {
	taskChan chan TaskParams
}

// NewDataGenerator creates a new data generator
func NewDataGenerator(bufferSize int) *DataGenerator {
	return &DataGenerator{
		taskChan: make(chan TaskParams, bufferSize),
	}
}

// Start begins the data generation process
func (g *DataGenerator) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(g.taskChan)
				return
			default:
				payload := generateJSONPayload()
				task := TaskParams{
					Args: payload,
					IdempotencyKey: pgtype.Text{
						String: uuid.NewString(),
						Valid:  true,
					},
				}
				g.taskChan <- task
			}
		}
	}()
}

// Tasks returns the channel for tasks
func (g *DataGenerator) Tasks() <-chan TaskParams {
	return g.taskChan
}
