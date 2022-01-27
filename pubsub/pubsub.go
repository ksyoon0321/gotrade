package pubsub

type PubsubTopic struct {
	id   string
	cmd  string
	data interface{}
}

func NewPubsubTopic(id, cmd string, data interface{}) PubsubTopic {
	return PubsubTopic{id: id, cmd: cmd, data: data}
}

func (p *PubsubTopic) GetId() string {
	return p.id
}

func (p *PubsubTopic) GetCmd() string {
	return p.cmd
}

func (p *PubsubTopic) GetData() interface{} {
	return p.data
}

type IPubsub interface {
	Subscribe(id string, nify interface{})
	Publish(topic PubsubTopic)
}
