package gosmpp

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/linxGnu/gosmpp/pdu"
	"strconv"
	"time"

	"github.com/allegro/bigcache/v3"
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
	store *bigcache.BigCache
}

func NewRequestStore() RequestStore {
	cache, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(30*time.Second))
	return &DefaultStore{
		store: cache,
	}
}

func (s DefaultStore) Set(ctx context.Context, request Request) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Task cancelled")
			return
		default:
			b, _ := serialize(request)
			_ = s.store.Set(strconv.Itoa(int(request.PDU.GetSequenceNumber())), b)
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
			bRequest, err := s.store.Get(strconv.Itoa(int(sequenceNumber)))
			if err != nil {
				return Request{}, false
			}
			request, err := deserialize(bRequest)
			if err != nil {
				return Request{}, false
			}
			return request, true
		}
	}
}

func (s DefaultStore) List(ctx context.Context) []Request {
	var requests []Request
	for {
		select {
		case <-ctx.Done():
			return requests
		default:
			it := s.store.Iterator()
			for it.SetNext() {
				value, err := it.Value()
				if err != nil {
					return requests
				}
				request, _ := deserialize(value.Value())
				requests = append(requests, request)
			}
			return requests
		}
	}
}

func (s DefaultStore) Delete(ctx context.Context, sequenceNumber int32) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_ = s.store.Delete(strconv.Itoa(int(sequenceNumber)))
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
			_ = s.store.Reset()
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
			return s.store.Len()
		}
	}
}

func serialize(request Request) ([]byte, error) {
	buf := pdu.NewBuffer(make([]byte, 0, 64))
	request.PDU.Marshal(buf)
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(requestGob{
		Pdu:      buf.Bytes(),
		TimeSent: time.Time{},
	})
	if err != nil {
		return b.Bytes()[:], errors.New("serialization failed")
	}
	return b.Bytes(), nil
}

func deserialize(bRequest []byte) (request Request, err error) {
	r := requestGob{}
	b := bytes.Buffer{}
	_, err = b.Write(bRequest)
	if err != nil {
		return request, errors.New("deserialization failed")
	}
	d := gob.NewDecoder(&b)
	err = d.Decode(&r)
	if err != nil {
		return request, errors.New("deserialization failed")
	}
	p, err := pdu.Parse(bytes.NewReader(r.Pdu))
	if err != nil {
		return Request{}, err
	}
	return Request{
		PDU:      p,
		TimeSent: r.TimeSent,
	}, nil
}

type requestGob struct {
	Pdu      []byte
	TimeSent time.Time
}
