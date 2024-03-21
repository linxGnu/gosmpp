package gosmpp

import (
	"context"
	"fmt"
	"github.com/linxGnu/gosmpp/pdu"
	cmap "github.com/orcaman/concurrent-map/v2"
	"golang.org/x/exp/maps"
	"strconv"
	"time"
)

type Request struct {
	pdu.PDU
	TimeSent time.Time
}

type Response struct {
	pdu.PDU
	OriginalRequest Request
}

type RequestStore interface {
	Set(ctx context.Context, request Request)
	Get(ctx context.Context, sequenceNumber int32) (Request, bool)
	List(ctx context.Context) []Request
	Delete(ctx context.Context, sequenceNumber int32)
	Clear(ctx context.Context)
	Length(ctx context.Context) int
}

type DefaultStore struct {
	store cmap.ConcurrentMap[string, Request]
}

func NewDefaultStore() DefaultStore {
	return DefaultStore{
		store: cmap.New[Request](),
	}
}

func (s DefaultStore) Set(ctx context.Context, request Request) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Task cancelled")
			return
		default:
			s.store.Set(strconv.Itoa(int(request.PDU.GetSequenceNumber())), request)
			return
		}
	}
}

func (s DefaultStore) Get(ctx context.Context, sequenceNumber int32) (Request, bool) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Task cancelled")
			return Request{}, false
		default:
			return s.store.Get(strconv.Itoa(int(sequenceNumber)))
		}
	}
}

func (s DefaultStore) List(ctx context.Context) []Request {
	for {
		select {
		case <-ctx.Done():
			return []Request{}
		default:
			return maps.Values(s.store.Items())
		}
	}
}

func (s DefaultStore) Delete(ctx context.Context, sequenceNumber int32) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			s.store.Remove(strconv.Itoa(int(sequenceNumber)))
			return
		}
	}
}

func (s DefaultStore) Clear(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			s.store.Clear()
			return
		}
	}
}

func (s DefaultStore) Length(ctx context.Context) int {
	for {
		select {
		case <-ctx.Done():
			return -1
		default:
			return s.store.Count()
		}
	}
}
