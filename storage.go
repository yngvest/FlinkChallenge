package main

import (
	"sync"
)

type Coordinates struct {
	Lat float32 `json:"lat"`
	Lng float32 `json:"lng"`
}

type History struct {
	coords []Coordinates
	sync.RWMutex
}

type OrdersHistory struct {
	orders map[string]*History
	sync.RWMutex
}

func NewOrdersHistory() *OrdersHistory {
	return &OrdersHistory{
		orders: make(map[string]*History),
	}
}

func (o *OrdersHistory) Add(order string, coordinates Coordinates) {
	history := o.getOrCreate(order)
	history.add(coordinates)
}

func (o *OrdersHistory) Get(order string, n int) []Coordinates {
	history := o.getOrder(order)
	if history == nil {
		return nil
	}
	return history.getLast(n)
}

func (o *OrdersHistory) Delete(order string) {
	o.Lock()
	defer o.Unlock()
	delete(o.orders, order)
}

func (o *OrdersHistory) getOrder(order string) *History {
	o.RLock()
	defer o.RUnlock()
	return o.orders[order]
}

func (o *OrdersHistory) getOrCreate(order string) *History {
	o.RLock()
	history, ok := o.orders[order]
	o.RUnlock()
	if !ok {
		h := &History{}
		o.Lock()
		history, ok = o.orders[order]
		if !ok {
			o.orders[order] = h
			history = h
		}
		o.Unlock()
	}
	return history
}

func (h *History) getLast(n int) []Coordinates {
	var from int
	h.RLock()
	l := len(h.coords)
	cnt := l
	if l > n && n != 0 {
		from = l - n
		cnt = n
	}
	res := make([]Coordinates, cnt)
	copy(res, h.coords[from:])
	h.RUnlock()
	return res
}

func (h *History) add(c Coordinates) {
	h.Lock()
	defer h.Unlock()
	h.coords = append(h.coords, c)
}
