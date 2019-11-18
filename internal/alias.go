package internal

import "github.com/lyraproj/dgo/dgo"

type aliasMap struct {
	typeNames  hashMap
	namedTypes hashMap
}

// NewAliasMap creates a new dgo.Alias map to be used as a scope when parsing types
func NewAliasMap() dgo.AliasMap {
	return &aliasMap{}
}

func (a *aliasMap) GetName(t dgo.Value) dgo.String {
	if v := a.typeNames.Get(t); v != nil {
		return v.(dgo.String)
	}
	return nil
}

func (a *aliasMap) GetType(n dgo.String) dgo.Value {
	return a.namedTypes.Get(n)
}

func (a *aliasMap) Add(t dgo.Value, name dgo.String) {
	a.typeNames.Put(t, name)
	a.namedTypes.Put(name, t)
}
