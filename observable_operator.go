package rxgo

import (
	"container/ring"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
)

// All determine whether all items emitted by an Observable meet some criteria.
func (o *observable) All(predicate Predicate, opts ...Option) Single {
	all := true
	return newSingleFromOperator(o, func(item Item, dst chan<- Item, operator operatorOptions) {
		if !predicate(item.Value) {
			dst <- FromValue(false)
			all = false
			operator.stop()
		}
	}, defaultErrorFuncOperator, func(dst chan<- Item) {
		if all {
			dst <- FromValue(true)
		}
	}, opts...)
}

// AverageFloat32 calculates the average of numbers emitted by an Observable and emits the average float32.
func (o *observable) AverageFloat32(opts ...Option) Single {
	var sum float32
	var count float32

	return newSingleFromOperator(o, func(item Item, dst chan<- Item, operator operatorOptions) {
		if v, ok := item.Value.(float32); ok {
			sum += v
			count++
		} else {
			dst <- FromError(errors.Wrap(&IllegalInputError{},
				fmt.Sprintf("expected type: float32, got: %t", item)))
			operator.stop()
		}
	}, defaultErrorFuncOperator, func(dst chan<- Item) {
		if count == 0 {
			dst <- FromValue(0)
		} else {
			dst <- FromValue(sum / count)
		}
	}, opts...)
}

// AverageFloat64 calculates the average of numbers emitted by an Observable and emits the average float64.
func (o *observable) AverageFloat64(opts ...Option) Single {
	var sum float64
	var count float64

	return newSingleFromOperator(o, func(item Item, dst chan<- Item, operator operatorOptions) {
		if v, ok := item.Value.(float64); ok {
			sum += v
			count++
		} else {
			dst <- FromError(errors.Wrap(&IllegalInputError{},
				fmt.Sprintf("expected type: float64, got: %t", item)))
			operator.stop()
		}
	}, defaultErrorFuncOperator, func(dst chan<- Item) {
		if count == 0 {
			dst <- FromValue(0)
		} else {
			dst <- FromValue(sum / count)
		}
	}, opts...)
}

// AverageInt calculates the average of numbers emitted by an Observable and emits the average int.
func (o *observable) AverageInt(opts ...Option) Single {
	var sum int
	var count int

	return newSingleFromOperator(o, func(item Item, dst chan<- Item, operator operatorOptions) {
		if v, ok := item.Value.(int); ok {
			sum += v
			count++
		} else {
			dst <- FromError(errors.Wrap(&IllegalInputError{},
				fmt.Sprintf("expected type: int, got: %t", item)))
			operator.stop()
		}
	}, defaultErrorFuncOperator, func(dst chan<- Item) {
		if count == 0 {
			dst <- FromValue(0)
		} else {
			dst <- FromValue(sum / count)
		}
	}, opts...)
}

// AverageInt8 calculates the average of numbers emitted by an Observable and emits the≤ average int8.
func (o *observable) AverageInt8(opts ...Option) Single {
	var sum int8
	var count int8

	return newSingleFromOperator(o, func(item Item, dst chan<- Item, operator operatorOptions) {
		if v, ok := item.Value.(int8); ok {
			sum += v
			count++
		} else {
			dst <- FromError(errors.Wrap(&IllegalInputError{},
				fmt.Sprintf("expected type: int8, got: %t", item)))
			operator.stop()
		}
	}, defaultErrorFuncOperator, func(dst chan<- Item) {
		if count == 0 {
			dst <- FromValue(0)
		} else {
			dst <- FromValue(sum / count)
		}
	}, opts...)
}

// AverageInt16 calculates the average of numbers emitted by an Observable and emits the average int16.
func (o *observable) AverageInt16(opts ...Option) Single {
	var sum int16
	var count int16

	return newSingleFromOperator(o, func(item Item, dst chan<- Item, operator operatorOptions) {
		if v, ok := item.Value.(int16); ok {
			sum += v
			count++
		} else {
			dst <- FromError(errors.Wrap(&IllegalInputError{},
				fmt.Sprintf("expected type: int16, got: %t", item)))
			operator.stop()
		}
	}, defaultErrorFuncOperator, func(dst chan<- Item) {
		if count == 0 {
			dst <- FromValue(0)
		} else {
			dst <- FromValue(sum / count)
		}
	}, opts...)
}

// AverageInt32 calculates the average of numbers emitted by an Observable and emits the average int32.
func (o *observable) AverageInt32(opts ...Option) Single {
	var sum int32
	var count int32

	return newSingleFromOperator(o, func(item Item, dst chan<- Item, operator operatorOptions) {
		if v, ok := item.Value.(int32); ok {
			sum += v
			count++
		} else {
			dst <- FromError(errors.Wrap(&IllegalInputError{},
				fmt.Sprintf("expected type: int32, got: %t", item)))
			operator.stop()
		}
	}, defaultErrorFuncOperator, func(dst chan<- Item) {
		if count == 0 {
			dst <- FromValue(0)
		} else {
			dst <- FromValue(sum / count)
		}
	}, opts...)
}

// AverageInt64 calculates the average of numbers emitted by an Observable and emits this average int64.
func (o *observable) AverageInt64(opts ...Option) Single {
	var sum int64
	var count int64

	return newSingleFromOperator(o, func(item Item, dst chan<- Item, operator operatorOptions) {
		if v, ok := item.Value.(int64); ok {
			sum += v
			count++
		} else {
			dst <- FromError(errors.Wrap(&IllegalInputError{},
				fmt.Sprintf("expected type: int64, got: %t", item)))
			operator.stop()
		}
	}, defaultErrorFuncOperator, func(dst chan<- Item) {
		if count == 0 {
			dst <- FromValue(0)
		} else {
			dst <- FromValue(sum / count)
		}
	}, opts...)
}

// BufferWithCount returns an Observable that emits buffers of items it collects
// from the source Observable.
// The resulting Observable emits buffers every skip items, each containing a slice of count items.
// When the source Observable completes or encounters an error,
// the resulting Observable emits the current buffer and propagates
// the notification from the source Observable.
func (o *observable) BufferWithCount(count, skip int, opts ...Option) Observable {
	if count <= 0 {
		return newObservableFromError(errors.Wrap(&IllegalInputError{}, "count must be positive"))
	}
	if skip <= 0 {
		return newObservableFromError(errors.Wrap(&IllegalInputError{}, "skip must be positive"))
	}

	buffer := make([]interface{}, count)
	iCount := 0
	iSkip := 0

	return newObservableFromOperator(o, func(item Item, dst chan<- Item, operator operatorOptions) {
		if iCount >= count {
			// Skip
			iSkip++
		} else {
			// Add to buffer
			buffer[iCount] = item.Value
			iCount++
			iSkip++
		}
		if iSkip == skip {
			// Send current buffer
			dst <- FromValue(buffer)
			buffer = make([]interface{}, count)
			iCount = 0
			iSkip = 0
		}
	}, func(item Item, dst chan<- Item, operator operatorOptions) {
		if iCount != 0 {
			dst <- FromValue(buffer[:iCount])
		}
		dst <- item
		iCount = 0
		operator.stop()
	}, func(dst chan<- Item) {
		if iCount != 0 {
			dst <- FromValue(buffer[:iCount])
		}
	}, opts...)
}

// BufferWithTime returns an Observable that emits buffers of items it collects from the source
// Observable. The resulting Observable starts a new buffer periodically, as determined by the
// timeshift argument. It emits each buffer after a fixed timespan, specified by the timespan argument.
// When the source Observable completes or encounters an error, the resulting Observable emits
// the current buffer and propagates the notification from the source Observable.
func (o *observable) BufferWithTime(timespan, timeshift Duration, opts ...Option) Observable {
	if timespan == nil || timespan.duration() == 0 {
		return newObservableFromError(errors.Wrap(&IllegalInputError{}, "timespan must no be nil"))
	}
	if timeshift == nil {
		timeshift = WithDuration(0)
	}

	var mux sync.Mutex
	var listenMutex sync.Mutex
	buffer := make([]interface{}, 0)
	stop := false
	listen := true

	option := parseOptions(opts...)
	next := option.buildChannel()
	ctx := option.buildContext()

	stopped := false

	// First goroutine in charge to check the timespan
	go func() {
		for {
			time.Sleep(timespan.duration())
			mux.Lock()
			if !stop {
				next <- FromValue(buffer)
				buffer = make([]interface{}, 0)
				mux.Unlock()

				if timeshift.duration() != 0 {
					listenMutex.Lock()
					listen = false
					listenMutex.Unlock()
					time.Sleep(timeshift.duration())
					listenMutex.Lock()
					listen = true
					listenMutex.Unlock()
				}
			} else {
				mux.Unlock()
				return
			}
		}
	}()

	go func() {
		observe := o.Observe()
	loop:
		for !stopped {
			select {
			case <-ctx.Done():
				break loop
			case i, ok := <-observe:
				if !ok {
					break loop
				}
				if i.IsError() {
					mux.Lock()
					if len(buffer) > 0 {
						next <- FromValue(buffer)
					}
					next <- i
					close(next)
					stop = true
					mux.Unlock()
					return
				}
				listenMutex.Lock()
				l := listen
				listenMutex.Unlock()

				mux.Lock()
				if l {
					buffer = append(buffer, i.Value)
				}
				mux.Unlock()
			}
		}
		mux.Lock()
		if len(buffer) > 0 {
			next <- FromValue(buffer)
		}
		close(next)
		stop = true
		mux.Unlock()
	}()

	return &observable{
		iterable: newChannelIterable(next),
	}
}

func (o *observable) BufferWithTimeOrCount(timespan Duration, count int, opts ...Option) Observable {
	if timespan == nil || timespan.duration() == 0 {
		return newObservableFromError(errors.Wrap(&IllegalInputError{}, "timespan must no be nil"))
	}
	if count <= 0 {
		return newObservableFromError(errors.Wrap(&IllegalInputError{}, "count must be positive"))
	}

	sendCh := make(chan []interface{})
	errCh := make(chan error)
	buffer := make([]interface{}, 0)
	var bufferMutex sync.Mutex
	option := parseOptions(opts...)
	next := option.buildChannel()
	ctx := option.buildContext()

	// First sender goroutine
	go func() {
		for {
			select {
			case currentBuffer := <-sendCh:
				next <- FromValue(currentBuffer)
			case error := <-errCh:
				bufferMutex.Lock()
				if len(buffer) > 0 {
					next <- FromValue(buffer)
				}
				bufferMutex.Unlock()
				if error != nil {
					next <- FromError(error)
				}
				close(next)
				return
			case <-time.After(timespan.duration()):
				bufferMutex.Lock()
				b := make([]interface{}, len(buffer))
				copy(b, buffer)
				buffer = make([]interface{}, 0)
				bufferMutex.Unlock()
				next <- FromValue(b)
			}
		}
	}()

	go func() {
		observe := o.Observe()
	loop:
		for {
			select {
			case <-ctx.Done():
				break loop
			case i, ok := <-observe:
				if !ok {
					break loop
				}
				if i.IsError() {
					errCh <- i.Err
					break loop
				}
				// TODO Improve implementation without mutex (sending data over channel)
				bufferMutex.Lock()
				buffer = append(buffer, i.Value)
				if len(buffer) >= count {
					b := make([]interface{}, len(buffer))
					copy(b, buffer)
					buffer = make([]interface{}, 0)
					bufferMutex.Unlock()
					sendCh <- b
				} else {
					bufferMutex.Unlock()
				}
			}
		}
		errCh <- nil
		close(sendCh)
		close(errCh)
	}()

	return &observable{
		iterable: newChannelIterable(next),
	}
}

// Contains determines whether an Observable emits a particular item or not.
func (o *observable) Contains(equal Predicate, opts ...Option) Single {
	return newSingleFromOperator(o, func(item Item, dst chan<- Item, operator operatorOptions) {
		if equal(item.Value) {
			dst <- FromValue(true)
			operator.stop()
			return
		}
	}, defaultErrorFuncOperator, func(dst chan<- Item) {
		dst <- FromValue(false)
	}, opts...)
}

// Count counts the number of items emitted by the source Observable and emit only this value.
func (o *observable) Count(opts ...Option) Single {
	var count int64
	return newSingleFromOperator(o, func(_ Item, dst chan<- Item, _ operatorOptions) {
		count++
	}, func(_ Item, dst chan<- Item, operator operatorOptions) {
		count++
		dst <- FromValue(count)
		operator.stop()
	}, defaultEndFuncOperator, opts...)
}

// DefaultIfEmpty returns an Observable that emits the items emitted by the source
// Observable or a specified default item if the source Observable is empty.
func (o *observable) DefaultIfEmpty(defaultValue interface{}, opts ...Option) Observable {
	empty := true

	return newObservableFromOperator(o, func(item Item, dst chan<- Item, operator operatorOptions) {
		empty = false
		dst <- item
	}, defaultErrorFuncOperator, func(dst chan<- Item) {
		if empty {
			dst <- FromValue(defaultValue)
		}
	}, opts...)
}

// Distinct suppresses duplicate items in the original Observable and returns
// a new Observable.
func (o *observable) Distinct(apply Func, opts ...Option) Observable {
	keyset := make(map[interface{}]interface{})

	return newObservableFromOperator(o, func(item Item, dst chan<- Item, operator operatorOptions) {
		key, err := apply(item.Value)
		if err != nil {
			dst <- FromError(err)
			operator.stop()
			return
		}
		_, ok := keyset[key]
		if !ok {
			dst <- item
		}
		keyset[key] = nil
	}, defaultErrorFuncOperator, defaultEndFuncOperator, opts...)
}

// DistinctUntilChanged suppresses consecutive duplicate items in the original
// Observable and returns a new Observable.
func (o *observable) DistinctUntilChanged(apply Func, opts ...Option) Observable {
	var current interface{}

	return newObservableFromOperator(o, func(item Item, dst chan<- Item, operator operatorOptions) {
		key, err := apply(item.Value)
		if err != nil {
			dst <- FromError(err)
			operator.stop()
			return
		}
		if current != key {
			dst <- item
			current = key
		}
	}, defaultErrorFuncOperator, defaultEndFuncOperator, opts...)
}

// ElementAt emits only item n emitted by an Observable.
func (o *observable) ElementAt(index uint, opts ...Option) Single {
	takeCount := 0
	sent := false

	return newSingleFromOperator(o, func(item Item, dst chan<- Item, operator operatorOptions) {
		if takeCount == int(index) {
			dst <- item
			sent = true
			operator.stop()
			return
		}
		takeCount++
	}, defaultErrorFuncOperator, func(dst chan<- Item) {
		if !sent {
			dst <- FromError(&IllegalInputError{})
		}
	}, opts...)
}

// Filter emits only those items from an Observable that pass a predicate test.
func (o *observable) Filter(apply Predicate, opts ...Option) Observable {
	return newObservableFromOperator(o, func(item Item, dst chan<- Item, operator operatorOptions) {
		if apply(item.Value) {
			dst <- item
		}
	}, defaultErrorFuncOperator, defaultEndFuncOperator, opts...)
}

// First returns new Observable which emit only first item.
func (o *observable) First(opts ...Option) OptionalSingle {
	return newSingleFromOperator(o, func(item Item, dst chan<- Item, operator operatorOptions) {
		dst <- item
		operator.stop()
	}, defaultErrorFuncOperator, defaultEndFuncOperator, opts...)
}

// FirstOrDefault returns new Observable which emit only first item.
// If the observable fails to emit any items, it emits a default value.
func (o *observable) FirstOrDefault(defaultValue interface{}, opts ...Option) Single {
	sent := false

	return newSingleFromOperator(o, func(item Item, dst chan<- Item, operator operatorOptions) {
		dst <- item
		sent = true
		operator.stop()
	}, defaultErrorFuncOperator, func(dst chan<- Item) {
		if !sent {
			dst <- FromValue(defaultValue)
		}
	}, opts...)
}

// ForEach subscribes to the Observable and receives notifications for each element.
func (o *observable) ForEach(nextFunc NextFunc, errFunc ErrFunc, doneFunc DoneFunc, opts ...Option) {
	handler := func(ctx context.Context, src <-chan Item, dst chan<- Item) {
		for {
			select {
			case <-ctx.Done():
				doneFunc()
				return
			case i, ok := <-src:
				if !ok {
					doneFunc()
					return
				}
				if i.IsError() {
					errFunc(i.Err)
					break
				}
				nextFunc(i.Value)
			}
		}
	}

	option := parseOptions(opts...)
	next := option.buildChannel()
	ctx := option.buildContext()
	go handler(ctx, o.Observe(), next)
}

// IgnoreElements ignores all items emitted by the source ObservableSource and only calls onComplete
// or onError.
func (o *observable) IgnoreElements(opts ...Option) Observable {
	return newObservableFromOperator(o, func(_ Item, _ chan<- Item, _ operatorOptions) {
	}, defaultErrorFuncOperator, defaultEndFuncOperator, opts...)
}

// Last returns a new Observable which emit only last item.
func (o *observable) Last(opts ...Option) OptionalSingle {
	var last Item
	empty := true

	return newOptionalSingleFromOperator(o, func(item Item, dst chan<- Item, operator operatorOptions) {
		last = item
		empty = false
	}, defaultErrorFuncOperator, func(dst chan<- Item) {
		if !empty {
			dst <- last
		}
	}, opts...)
}

// LastOrDefault returns a new Observable which emit only last item.
// If the observable fails to emit any items, it emits a default value.
func (o *observable) LastOrDefault(defaultValue interface{}, opts ...Option) Single {
	var last Item
	empty := true

	return newSingleFromOperator(o, func(item Item, dst chan<- Item, operator operatorOptions) {
		last = item
		empty = false
	}, defaultErrorFuncOperator, func(dst chan<- Item) {
		if !empty {
			dst <- last
		} else {
			dst <- FromValue(defaultValue)
		}
	}, opts...)
}

// Map transforms the items emitted by an Observable by applying a function to each item.
func (o *observable) Map(apply Func, opts ...Option) Observable {
	return newObservableFromOperator(o, func(item Item, dst chan<- Item, operator operatorOptions) {
		res, err := apply(item.Value)
		if err != nil {
			dst <- FromError(err)
			operator.stop()
		}
		dst <- FromValue(res)
	}, defaultErrorFuncOperator, defaultEndFuncOperator, opts...)
}

// Marshal transforms the items emitted by an Observable by applying a marshalling to each item.
func (o *observable) Marshal(marshaler Marshaler, opts ...Option) Observable {
	return o.Map(func(i interface{}) (interface{}, error) {
		return marshaler(i)
	}, opts...)
}

// Max determines and emits the maximum-valued item emitted by an Observable according to a comparator.
func (o *observable) Max(comparator Comparator, opts ...Option) OptionalSingle {
	empty := true
	var max interface{}

	return newObservableFromOperator(o, func(item Item, dst chan<- Item, operator operatorOptions) {
		empty = false

		if max == nil {
			max = item.Value
		} else {
			if comparator(max, item.Value) < 0 {
				max = item.Value
			}
		}
	}, defaultErrorFuncOperator, func(dst chan<- Item) {
		if !empty {
			dst <- FromValue(max)
		}
	}, opts...)
}

// Min determines and emits the minimum-valued item emitted by an Observable according to a comparator.
func (o *observable) Min(comparator Comparator, opts ...Option) OptionalSingle {
	empty := true
	var min interface{}

	return newObservableFromOperator(o, func(item Item, dst chan<- Item, operator operatorOptions) {
		empty = false

		if min == nil {
			min = item.Value
		} else {
			if comparator(min, item.Value) > 0 {
				min = item.Value
			}
		}
	}, defaultErrorFuncOperator, func(dst chan<- Item) {
		if !empty {
			dst <- FromValue(min)
		}
	}, opts...)
}

// Observe observes an observable by returning its channel
func (o *observable) Observe(opts ...Option) <-chan Item {
	return o.iterable.Observe(opts...)
}

func (o *observable) OnErrorResumeNext(resumeSequence ErrorToObservable, opts ...Option) Observable {
	return newObservableFromOperator(o, func(item Item, dst chan<- Item, operator operatorOptions) {
		dst <- item
	}, func(item Item, dst chan<- Item, operator operatorOptions) {
		operator.resetIterable(resumeSequence(item.Err))
	}, defaultEndFuncOperator, opts...)
}

// Retry retries if a source Observable sends an error, resubscribe to it in the hopes that it will complete without error.
func (o *observable) Retry(count int, opts ...Option) Observable {
	option := parseOptions(opts...)
	next := option.buildChannel()
	ctx := option.buildContext()

	go func() {
		observe := o.Observe()
	loop:
		for {
			select {
			case <-ctx.Done():
				break loop
			case i, ok := <-observe:
				if !ok {
					break loop
				}
				if i.IsError() {
					count--
					if count < 0 {
						next <- i
						break loop
					}
					observe = o.Observe()
				} else {
					next <- i
				}
			}
		}
		close(next)
	}()

	return &observable{
		iterable: newChannelIterable(next),
	}
}

// SkipWhile discard items emitted by an Observable until a specified condition becomes false.
func (o *observable) SkipWhile(apply Predicate, opts ...Option) Observable {
	skip := true

	return newObservableFromOperator(o, func(item Item, dst chan<- Item, operator operatorOptions) {
		if !skip {
			dst <- item
		} else {
			if !apply(item.Value) {
				skip = false
				dst <- item
			}
		}
	}, defaultErrorFuncOperator, defaultEndFuncOperator, opts...)
}

// Take emits only the first n items emitted by an Observable.
func (o *observable) Take(nth uint, opts ...Option) Observable {
	takeCount := 0

	return newObservableFromOperator(o, func(item Item, dst chan<- Item, operator operatorOptions) {
		if takeCount < int(nth) {
			takeCount++
			dst <- item
		}
	}, defaultErrorFuncOperator, defaultEndFuncOperator, opts...)
}

// TakeLast emits only the last n items emitted by an Observable.
func (o *observable) TakeLast(nth uint, opts ...Option) Observable {
	n := int(nth)
	r := ring.New(n)
	count := 0

	return newObservableFromOperator(o, func(item Item, dst chan<- Item, operator operatorOptions) {
		count++
		r.Value = item.Value
		r = r.Next()
	}, defaultErrorFuncOperator, func(dst chan<- Item) {
		if count < n {
			remaining := n - count
			if remaining <= count {
				r = r.Move(n - count)
			} else {
				r = r.Move(-count)
			}
			n = count
		}
		for i := 0; i < n; i++ {
			dst <- FromValue(r.Value)
			r = r.Next()
		}
	}, opts...)
}

// ToSlice collects all items from an Observable and emit them as a single slice.
func (o *observable) ToSlice(opts ...Option) Single {
	s := make([]interface{}, 0)
	return newSingleFromOperator(o, func(item Item, dst chan<- Item, operator operatorOptions) {
		s = append(s, item.Value)
	}, defaultErrorFuncOperator, func(dst chan<- Item) {
		dst <- FromValue(s)
	}, opts...)
}

// Marshal transforms the items emitted by an Observable by applying an unmarshalling to each item.
func (o *observable) Unmarshal(unmarshaler Unmarshaler, factory func() interface{}, opts ...Option) Observable {
	return o.Map(func(i interface{}) (interface{}, error) {
		v := factory()
		err := unmarshaler(i.([]byte), v)
		if err != nil {
			return nil, err
		}
		return v, nil
	}, opts...)
}
