// Built from ideas in http://github.com/mattn/go-pubsub

// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and 
// associated documentation files (the “Software”), to deal in the Software without restriction, including 
// without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell 
// copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the 
// following conditions:
//
// The above copyright notice and this permission notice shall be included in all copies or substantial 
// portions of the Software.

package tap

import (
	"reflect"
	"sync"
)

type EventManager struct {
	channel chan interface{}
	listeners []interface{}
	mutex sync.Mutex
	eventType reflect.Type
}

func NewEventManager(eventType reflect.Type) *EventManager {
	result := new(EventManager)
	result.eventType = eventType
	result.channel = make(chan interface{})

	go func() {
		for event := range result.channel {
			result.mutex.Lock()
			for _, listener := range result.listeners {
				go listener.(func(interface{}))(event)
			}
			result.mutex.Unlock()
		}
	}()
	
	return result
}

func (manager *EventManager) On(listener interface{}) {

	t := reflect.TypeOf(listener)
	if t.Kind() != reflect.Func {
		return
	}

	paramTypes := make([]reflect.Type, t.NumIn())
	if (paramTypes[0] != manager.eventType) {
		return
	}

	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	manager.listeners = append(manager.listeners, listener)
}

func (manager *EventManager) Fire(event interface{}) {
	if reflect.TypeOf(event).Implements(manager.eventType) {
		manager.channel <- event
	}
}

func (manager *EventManager) Close() {
	close(manager.channel)
	manager.listeners = nil
}