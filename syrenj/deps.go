package syrenj

import "reflect"

type Deps struct {
	deps map[string]interface{}
	providers map[string]interface{}
}

func (w Deps) AddDep(id string, dep interface{}) Deps {
	w.deps[id] = dep
	return w
}

func (w Deps) AddProvider(id string, provider interface{}) Deps {
	w.providers[id] = provider
	return w
}

func (w Deps) Get(id string) interface{} {
	return w.deps[id]
}

func Run() {
	
}

func Stop()  {
	
}
