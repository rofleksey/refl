package runtime

import "iter"

type Variable struct {
	value Object
}

type Environment struct {
	globalEnv *Environment
	parent    *Environment
	values    map[string]*Variable // not thread safe
}

func NewEnvironment(parent *Environment) *Environment {
	result := &Environment{
		parent: parent,
		values: make(map[string]*Variable),
	}

	var globalEnv *Environment
	if parent != nil {
		globalEnv = parent.globalEnv
	} else {
		globalEnv = result
	}

	result.globalEnv = globalEnv

	return result
}

func (e *Environment) Clone() *Environment {
	cloned := NewEnvironment(e.parent)
	for key, variable := range e.values {
		cloned.values[key] = variable
	}
	return cloned
}

func (e *Environment) Get(name string) (Object, bool) {
	variable, ok := e.values[name]
	if !ok && e.parent != nil {
		return e.parent.Get(name)
	}
	if ok {
		return variable.value, true
	}
	return nil, false
}

func (e *Environment) Set(name string, value Object) {
	current := e
	for current != nil {
		if variable, ok := current.values[name]; ok {
			variable.value = value
			return
		}
		current = current.parent
	}

	e.globalEnv.Define(name, value)
}

func (e *Environment) Define(name string, value Object) {
	e.values[name] = &Variable{value}
}

func (e *Environment) Delete(name string) {
	delete(e.values, name)
}

func (e *Environment) GetAll() map[string]Object {
	result := make(map[string]Object)

	// Walk up the chain and collect all variables
	current := e
	for current != nil {
		for k, v := range current.values {
			if _, exists := result[k]; !exists {
				result[k] = v.value
			}
		}
		current = current.parent
	}

	return result
}

func (e *Environment) GlobalsIterator() iter.Seq2[string, Object] {
	return func(yield func(string, Object) bool) {
		for name, variable := range e.globalEnv.values {
			if !yield(name, variable.value) {
				return
			}
		}
	}
}
