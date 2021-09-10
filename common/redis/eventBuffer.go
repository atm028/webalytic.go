package RedisBroker

type IEventBuffer interface {
}

type EventBufferConfg struct {
	FlushTimeout int
	MaxLength    int
}
