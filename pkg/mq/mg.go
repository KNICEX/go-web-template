package mq

import (
	"sync"
	"time"
)

type Message struct {
	Publisher   string
	PublishTime int64
	Event       string
	Content     interface{}
}

type CallbackFunc func(Message)

type MQ interface {
	Publish(string, Message)
	Subscribe(string, int) <-chan Message
	SubscribeCallback(string, CallbackFunc)
	Unsubscribe(string, <-chan Message)
}

var GlobalMQ = NewMQ()

type memoryMQ struct {
	topics    map[string][]chan Message
	callbacks map[string][]CallbackFunc
	sync.RWMutex
}

func NewMQ() MQ {
	return &memoryMQ{
		topics:    make(map[string][]chan Message),
		callbacks: make(map[string][]CallbackFunc),
	}
}

func (mq *memoryMQ) Publish(topic string, message Message) {
	mq.RLock()
	subscribersChan, okChan := mq.topics[topic]
	subscribeCallback, okCallback := mq.callbacks[topic]
	mq.RUnlock()

	if okChan {
		go func(subscribersChan []chan Message) {
			for _, ch := range subscribersChan {
				select {
				case ch <- message:
				case <-time.After(time.Millisecond * 500):
				}
			}

		}(subscribersChan)
	}

	if okCallback {
		for _, callback := range subscribeCallback {
			callback(message)
		}
	}
}

func (mq *memoryMQ) Subscribe(topic string, bufferSize int) <-chan Message {
	ch := make(chan Message, bufferSize)
	mq.Lock()
	defer mq.Unlock()
	mq.topics[topic] = append(mq.topics[topic], ch)
	return ch
}

func (mq *memoryMQ) SubscribeCallback(topic string, callback CallbackFunc) {
	mq.Lock()
	defer mq.Unlock()
	mq.callbacks[topic] = append(mq.callbacks[topic], callback)
}

func (mq *memoryMQ) Unsubscribe(topic string, sub <-chan Message) {
	mq.Lock()
	defer mq.Unlock()
	subs := mq.topics[topic]
	for i, ch := range subs {
		if ch == sub {
			subs = append(subs[:i], subs[i+1:]...)
			break
		}
	}
	mq.topics[topic] = subs
}
