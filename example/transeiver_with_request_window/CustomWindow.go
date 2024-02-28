package main

import (
	"context"
	"github.com/linxGnu/gosmpp"
	"golang.org/x/exp/maps"
	"strconv"

	"github.com/orcaman/concurrent-map/v2"
)

type CustomWindow struct {
	store cmap.ConcurrentMap[string, gosmpp.Request]
}

func (w CustomWindow) Set(ctx context.Context, request gosmpp.Request) {
	w.store.Set(strconv.Itoa(int(request.PDU.GetSequenceNumber())), request)
}

func (w CustomWindow) Get(ctx context.Context, sequenceNumber int32) (gosmpp.Request, bool) {
	return w.store.Get(strconv.Itoa(int(sequenceNumber)))
}

func (w CustomWindow) List(ctx context.Context) []gosmpp.Request {
	return maps.Values(w.store.Items())
	//return maps.Values(w.store)
}

func (w CustomWindow) Delete(ctx context.Context, sequenceNumber int32) {
	w.store.Remove(strconv.Itoa(int(sequenceNumber)))
}

func (w CustomWindow) Clear(ctx context.Context) {
	w.store.Clear()
}

func (w CustomWindow) Length(ctx context.Context) int {
	return w.store.Count()
}

func NewCustomWindow() gosmpp.RequestWindowStore {
	return &CustomWindow{
		store: cmap.New[gosmpp.Request](),
	}
}
