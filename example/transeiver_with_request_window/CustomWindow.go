package main

import (
	"context"
	"fmt"
	cmap "github.com/orcaman/concurrent-map/v2"
	"strconv"

	"github.com/linxGnu/gosmpp"
	"golang.org/x/exp/maps"
)

type CustomWindow struct {
	store cmap.ConcurrentMap[string, gosmpp.Request]
}

func (w CustomWindow) Set(ctx context.Context, request gosmpp.Request) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Task cancelled")
			return
		default:
			w.store.Set(strconv.Itoa(int(request.PDU.GetSequenceNumber())), request)
		}
	}
}

func (w CustomWindow) Get(ctx context.Context, sequenceNumber int32) (gosmpp.Request, bool) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Task cancelled")
			return gosmpp.Request{}, false
		default:
			return w.store.Get(strconv.Itoa(int(sequenceNumber)))
		}
	}
}

func (w CustomWindow) List(ctx context.Context) []gosmpp.Request {
	for {
		select {
		case <-ctx.Done():
			return []gosmpp.Request{}
		default:
			return maps.Values(w.store.Items())
		}
	}
}

func (w CustomWindow) Delete(ctx context.Context, sequenceNumber int32) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			w.store.Remove(strconv.Itoa(int(sequenceNumber)))
		}
	}
}

func (w CustomWindow) Clear(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			w.store.Clear()
		}
	}
}

func (w CustomWindow) Length(ctx context.Context) int {
	for {
		select {
		case <-ctx.Done():
			return -1
		default:
			return w.store.Count()
		}
	}
}

func NewCustomWindow() gosmpp.RequestWindowStore {
	return &CustomWindow{
		store: cmap.New[gosmpp.Request](),
	}
}
