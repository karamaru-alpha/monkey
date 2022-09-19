package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Object, 0)}
}

func (e *Environment) Get(key string) (Object, bool) {
	obj, ok := e.store[key]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(key)
	}
	return obj, ok
}

func (e *Environment) Set(key string, val Object) Object {
	e.store[key] = val
	return val

}
