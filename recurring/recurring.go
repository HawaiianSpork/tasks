// Copyright 2013 Travis Keep. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or
// at http://opensource.org/licenses/BSD-3-Clause.

// Package recurring contains implementations of the the Recurring interface
// in the tasks package.
package recurring

import (
  "github.com/keep94/gofunctional3/functional"
  "time"
)

// R represents a recurring time such as each Monday at 7:00.
type R interface {

  // ForTime returns a Stream of time.Time starting at t. The times that
  // returned Stream emits shall be after t and be in ascending order.
  ForTime(t time.Time) functional.Stream
}

// RFunc converts an ordinary function to an R instance
type RFunc func(t time.Time) functional.Stream

func (f RFunc) ForTime(t time.Time) functional.Stream {
  return f(t)
}

// Combine combines multiple R instances together and returns them
// as a single one.
func Combine(rs ...R) R {
  return RFunc(func(t time.Time) functional.Stream {
    streams := make([]functional.Stream, len(rs))
    for i := range rs {
      streams[i] = rs[i].ForTime(
    }
    return combineStreams(streams)
  }
}

// Modify returns a new R instance that uses f to modify the time.Time
// Streams that r creates.
func Modify(r R, f func(s functional.Stream) functional.Stream) R {
  return RFunc(func(t time.Time) functional.Stream {
    return f(r.ForTime(t))
  })
}

// Filter returns a new R instance that filters the time.Time Streams
// that r creates
func Filter(r R, f functional.Filterer) R {
  return Modify(
      r,
      func(s functional.Stream) functional.Stream {
        return functional.Filter(f, s)
      })
}

// FirstN returns a new R instance that generates only the first N
// times that r generates.
func FirstN(r R, n int) R {
  return Modify(
      r,
      func(s functional.Stream) functional.Stream {
        return functional.Slice(s, 0, n)
      })
}

// AtInterval returns a new R instance that represents repeating
// at d intervals.
func AtInterval(d time.Duration) R {
  return RFunc(func(t time.Time) functional.Stream {
    return &intervalStream{t: t.Add(d), d: d}
  })
}
  
// AtTime returns a new R instance that represents repeating at a
// certain time of day.
func AtTime(hour24, minute int) R {
  return RFunc(func(t time.Time) functional.Stream {
    firstT := time.Date(t.Year(), t.Month(), t.Day(), hour24, minute, 0, 0, t.Location())
    if !firstT.After(t) {
      firstT = firstT.AddDate(0, 0, 1)
    }
    return &dateStream{t: firstT}
  })
}

type dateStream struct {
  t time.Time
  closeDoesNothing
}

func (s *dateStream) Next(ptr interface{}) error {
  p := ptr.(*time.Time)
  *p = s.t
  s.t = s.t.AddDate(0, 0, 1)
  return nil
}

type intervalStream struct {
  t time.Time
  d time.Duration
  closeDoesNothing
}

func (s *intervalStream) Next(ptr interface{}) error {
  p := ptr.(*time.Time)
  *p = s.t
  s.t = s.t.Add(s.d)
  return nil
}

type closeDoesNothing struct {
}

func (n closeDoesNothing) Close() error {
  return nil
}

func combineStreams(streams []functional.Stream) functional.Stream {
  h := make(streamHeap, len(streams))
  for i := range streams {
    h[i] = &item{stream: streams[i]}
    h[i].pop()
  }
  heap.Init(&h)
  return &mergeStream{streams: streams, sh: h}
}

type item struct {
  stream functional.Stream
  t time.Time
  e error
}

func (i *item) pop() {
  i.e = i.stream.Next(&i.t)
}

type streamHeap []*item

func (sh streamHeap) Len() int {
  return len(sh)
}

func (sh streamHeap) Less(i, j int) bool {
  if sh[i].e != nil {
    return sh[j].e == nil
  }
  if sh[j].e != nil {
    return false
  }
  return sh[i].t.Before(sh[j].t)
}

func (sh streamHeap) Swap(i, j int) {
  return sh[i], sh[j] = sh[j], sh[i]
}

func (sh streamHeap) Push(x interface{}) {
  k

  