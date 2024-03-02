package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/allegro/bigcache/v3"
	"github.com/linxGnu/gosmpp/pdu"
	"strconv"
	"time"

	"github.com/linxGnu/gosmpp"
)

// This is a just an example how to implement a custom store.
//
// Your implementation must be concurrency safe
//
// In this example we use bigcache https://github.com/allegro/bigcache
// Warning:
//  - This is just an example and should be tested before using in production
//	- We are serializing with gob, some field cannot be serialized for simplicity
//  - We recommend you implement your own serialization/deserialization if you choose to use bigcache

type CustomStore struct {
	store *bigcache.BigCache
}

func NewCustomStore() CustomStore {
	cache, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(30*time.Second))
	return CustomStore{
		store: cache,
	}
}

func (s CustomStore) Set(ctx context.Context, request gosmpp.Request) {
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

func (s CustomStore) Get(ctx context.Context, sequenceNumber int32) (gosmpp.Request, bool) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Task cancelled")
			return gosmpp.Request{}, false
		default:
			bRequest, err := s.store.Get(strconv.Itoa(int(sequenceNumber)))
			if err != nil {
				return gosmpp.Request{}, false
			}
			request, err := deserialize(bRequest)
			if err != nil {
				return gosmpp.Request{}, false
			}
			return request, true
		}
	}
}

func (s CustomStore) List(ctx context.Context) []gosmpp.Request {
	var requests []gosmpp.Request
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

func (s CustomStore) Delete(ctx context.Context, sequenceNumber int32) {
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

func (s CustomStore) Clear(ctx context.Context) {
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

func (s CustomStore) Length(ctx context.Context) int {
	for {
		select {
		case <-ctx.Done():
			return -1

		default:
			return s.store.Len()
		}
	}
}

func serialize(request gosmpp.Request) ([]byte, error) {
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

func deserialize(bRequest []byte) (request gosmpp.Request, err error) {
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
		return gosmpp.Request{}, err
	}
	return gosmpp.Request{
		PDU:      p,
		TimeSent: r.TimeSent,
	}, nil
}

type requestGob struct {
	Pdu      []byte
	TimeSent time.Time
}
