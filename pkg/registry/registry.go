package registry

import (
	"context"
	"github.com/vksir/vkiss-lib/pkg/log"
	"sync"
)

var topics = make(map[string]map[string]Callback)
var topicLock sync.RWMutex

type Callback = func(ctx context.Context, msgAny any)

func Notify(ctx context.Context, topic string, msg any) {
	topicLock.RLock()
	defer topicLock.RUnlock()

	log.InfoC(ctx, "notify", "topic", topic, "msg", msg)
	tp, ok := topics[topic]
	if ok {
		for _, callback := range tp {
			callback(ctx, msg)
		}
	}
}

func Subscribe(topic, subscriber string, callback Callback) {
	topicLock.Lock()
	defer topicLock.Unlock()

	tp, ok := topics[topic]
	if !ok {
		topics[topic] = make(map[string]Callback)
	}
	tp[subscriber] = callback
	log.Warn("subscribe success", "topic", topic, "subscriber", subscriber, "callback", callback)
}

func Unsubscribe(topic, subscriber string) {
	topicLock.Lock()
	defer topicLock.Unlock()

	tp, ok := topics[topic]
	if ok {
		delete(tp, subscriber)
		if len(tp) == 0 {
			delete(topics, topic)
		}
	}
	log.Warn("unsubscribe success", "topic", topic, "subscriber", subscriber)
}
