package pubsub

import (
	"sync"
)

type PubusbChan struct {
	mutex sync.RWMutex
	subs  map[string]chan PubsubTopic
}

func NewPubSubChan() IPubsub {
	return &PubusbChan{subs: make(map[string]chan PubsubTopic)}
}

func (p *PubusbChan) Subscribe(id string, nify interface{}) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.subs[id] = nify.(chan PubsubTopic)
}

func (p *PubusbChan) Publish(topic PubsubTopic) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	for _, ch := range p.subs {
		ch <- topic
	}
}
