package main

import (
	"context"
	"fmt"
	cmap "github.com/orcaman/concurrent-map/v2"
	"strconv"

	"github.com/linxGnu/gosmpp"
	"golang.org/x/exp/maps"
)

type CustomStore struct {
	store cmap.ConcurrentMap[string, gosmpp.Request]
}

func NewCustomStore() gosmpp.RequestStore {
	return &CustomStore{
		store: cmap.New[gosmpp.Request](),
	}
}

func (s CustomStore) Set(ctx context.Context, request gosmpp.Request) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Task cancelled")
			return
		default:
			s.store.Set(strconv.Itoa(int(request.PDU.GetSequenceNumber())), request)
		}
	}
}

func (s CustomStore) Get(ctx context.Context, sequenceNumber int32) (gosmpp.Request, bool) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Task cancelled")
			return gosmpp.Request{}, false
		default:
			return s.store.Get(strconv.Itoa(int(sequenceNumber)))
		}
	}
}

func (s CustomStore) List(ctx context.Context) []gosmpp.Request {
	for {
		select {
		case <-ctx.Done():
			return []gosmpp.Request{}
		default:
			return maps.Values(s.store.Items())
		}
	}
}

func (s CustomStore) Delete(ctx context.Context, sequenceNumber int32) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			s.store.Remove(strconv.Itoa(int(sequenceNumber)))
		}
	}
}

func (s CustomStore) Clear(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			s.store.Clear()
		}
	}
}

func (s CustomStore) Length(ctx context.Context) int {
	for {
		select {
		case <-ctx.Done():
			return -1
		default:
			return s.store.Count()
		}
	}
}
