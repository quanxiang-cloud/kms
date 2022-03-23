package adaptor

import "fmt"

// ServiceInstance service instance
type ServiceInstance interface {
}

type instances map[string]ServiceInstance

var (
	adaptorIns = instances{}
)

func setInstance(name string, instance ServiceInstance) {
	adaptorIns[name] = instance
}

func getInstance(name string) (ServiceInstance, error) {
	if ins, ok := adaptorIns[name]; ok {
		return ins, nil
	}
	return nil, fmt.Errorf("not exist instance of %s", name)
}
