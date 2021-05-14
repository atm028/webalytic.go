package RedisBroker

import (
	inversify "github.com/alekns/go-inversify"
)

func Container(container inversify.Container) inversify.Container {
	broker := Init(container)

	container.Bind("redisBroker").To(broker)
	return container
}
