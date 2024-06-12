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

// Request represent a request tracked by the RequestStore
type Request struct {
	pdu.PDU
	TimeSent time.Time
}

// Response represents a response from a Request in the RequestStore
type Response struct {
	pdu.PDU
	OriginalRequest Request
}

// RequestStore interface used for WindowedRequestTracking
type RequestStore interface {
	Set(ctx context.Context, request Request) error
	Get(ctx context.Context, sequenceNumber int32) (Request, bool)
	List(ctx context.Context) []Request
	Delete(ctx context.Context, sequenceNumber int32) error
	Clear(ctx context.Context) error
	Length(ctx context.Context) (int, error)
}

type DefaultStore struct {
	store cmap.ConcurrentMap[string, Request]
}

func NewDefaultStore() DefaultStore {
	return DefaultStore{
		store: cmap.New[Request](),
	}
}

func (s DefaultStore) Set(ctx context.Context, request Request) error {
	select {
	case <-ctx.Done():
		fmt.Println("Task cancelled")
		return ctx.Err()
	default:
		s.store.Set(strconv.Itoa(int(request.PDU.GetSequenceNumber())), request)
		return nil
	}
}

func (s DefaultStore) Get(ctx context.Context, sequenceNumber int32) (Request, bool) {
	select {
	case <-ctx.Done():
		fmt.Println("Task cancelled")
		return Request{}, false
	default:
		return s.store.Get(strconv.Itoa(int(sequenceNumber)))
	}
}

func (s DefaultStore) List(ctx context.Context) []Request {
	select {
	case <-ctx.Done():
		return []Request{}
	default:
		return maps.Values(s.store.Items())
	}
}

func (s DefaultStore) Delete(ctx context.Context, sequenceNumber int32) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		s.store.Remove(strconv.Itoa(int(sequenceNumber)))
		return nil
	}
}

func (s DefaultStore) Clear(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		s.store.Clear()
		return nil
	}
}

func (s DefaultStore) Length(ctx context.Context) (int, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
		return s.store.Count(), nil
	}
}
