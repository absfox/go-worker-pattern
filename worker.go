package main

import (
	"context"
	"log"
	"time"
)

type Worker interface {
	Start()
	Stop()
}

type worker struct {
	// buffered channel ???
	incoming chan interface{}

	//
	maxWorkers int

	pollInterval time.Duration

	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}
}

func newWorker() *worker {
	return &worker{
		done:         make(chan struct{}),
		incoming:     make(chan interface{}, 128),
		maxWorkers:   5,
		pollInterval: 10 * time.Second,
	}
}

func (w *worker) Start() {
	w.ctx, w.cancel = context.WithCancel(context.Background())
	go w.startReceivingMessages()
	go w.startChildWorkers()
}

// fetch messages from sqs/sns and write to channel consumption channel
func (w *worker) startReceivingMessages() {
	ticker := time.NewTicker(w.pollInterval)

	for {
		select {

		case <-ticker.C:
			msgs, err := w.fetchFromBackend()
			if err != nil {
				//
				// TODO:
				//
				// log.....
				continue
			}

			// Write to incoming channel for consumption
			for _, msg := range msgs {
				select {
				case w.incoming <- msg:
				case <-w.ctx.Done():
					return
				}
			}

		case <-w.ctx.Done():
			ticker.Stop()
			close(w.done)
			return

		}
	}
}

func (w *worker) fetchFromBackend() ([]interface{}, error) {
	//
	// TODO: fetch from sns/sqs
	//
	// resp, err:=client.SQS.ReceiveMessage(&params)
	// if err !=nil {
	//
	//}
	// return resp.Messages, nil
	log.Println("TODO fetch from sns/sqs")

	return []interface{}{}, nil
}

func (w *worker) startChildWorkers() {
	for i := 0; i < w.maxWorkers; i++ {
		go w.startWorker()
	}
}

func (w *worker) startWorker() {
	for {
		select {
		case msg := <-w.incoming:
			//
			// TODO: process the message
			//
			log.Println("TODO: process the message", msg)

		case <-w.ctx.Done():
			return

		}
	}
}

func (w *worker) Stop() {
	w.cancel()
	<-w.done
}
