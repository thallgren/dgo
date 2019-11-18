package loader

import (
	"fmt"
	"strings"
	"sync"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/vf"
)

type (
	multipleEntries struct {
		dgo.Map
	}

	mapLoader struct {
		name     string
		parentNs dgo.Loader
		entries  dgo.Map
	}

	loader struct {
		mapLoader
		lock       sync.RWMutex
		namespaces dgo.Map
		finder     dgo.Finder
		nsCreator  dgo.NsCreator
	}

	childLoader struct {
		dgo.Loader
		parent dgo.Loader
	}
)

// Multiple creates an value that holds multiple loader entries. The function is used by a finder that
// wishes to return more than the one entry that was requested. The requested element must be one of the
// entries.
func Multiple(m dgo.Map) dgo.Value {
	return multipleEntries{m}
}

func load(l dgo.Loader, name string) dgo.Value {
	parts := strings.Split(name, `/`)
	last := len(parts) - 1
	for i := 0; i < last; i++ {
		l = l.Namespace(parts[i])
		if l == nil {
			return nil
		}
	}
	return l.Get(parts[last])
}

// New returns a new Loader instance
func New(parentNs dgo.Loader, name string, entries dgo.Map, finder dgo.Finder, nsCreator dgo.NsCreator) dgo.Loader {
	if entries == nil {
		entries = vf.Map()
	}
	if finder == nil && nsCreator == nil {
		// Immutable map based loader
		return &mapLoader{parentNs: parentNs, name: name, entries: entries.FrozenCopy().(dgo.Map)}
	}

	var namespaces dgo.Map
	if nsCreator == nil {
		namespaces = vf.Map()
	} else {
		namespaces = vf.MutableMap()
	}
	return &loader{
		mapLoader:  mapLoader{parentNs: parentNs, name: name, entries: entries.Copy(finder == nil)},
		namespaces: namespaces,
		finder:     finder,
		nsCreator:  nsCreator}
}

func (l *mapLoader) Get(key interface{}) dgo.Value {
	return l.entries.Get(key)
}

func (l *mapLoader) Load(name string) dgo.Value {
	return load(l, name)
}

func (l *mapLoader) AbsoluteName() string {
	an := `/` + l.name
	if l.parentNs != nil {
		if pn := l.parentNs.AbsoluteName(); pn != `/` {
			an = pn + an
		}
	}
	return an
}

func (l *mapLoader) Name() string {
	return l.name
}

func (l *mapLoader) Namespace(name string) dgo.Loader {
	if name == `` {
		return l
	}
	return nil
}

func (l *mapLoader) NewChild(finder dgo.Finder, nsCreator dgo.NsCreator) dgo.Loader {
	return loaderWithParent(l, finder, nsCreator)
}

func (l *mapLoader) ParentNamespace() dgo.Loader {
	return l.parentNs
}

func (l *loader) add(key, value dgo.Value) dgo.Value {
	l.lock.Lock()
	defer l.lock.Unlock()

	addEntry := func(key, value dgo.Value) {
		if old := l.entries.Get(key); old == nil {
			l.entries.Put(key, value)
		} else if !old.Equals(value) {
			panic(fmt.Errorf(`attempt to override entry %q`, key))
		}
	}

	if m, ok := value.(multipleEntries); ok {
		value = m.Get(key)
		if value == nil {
			panic(fmt.Errorf(`map returned from finder doesn't contain original key %q`, key))
		}
		m.EachEntry(func(e dgo.MapEntry) { addEntry(e.Key(), e.Value()) })
	} else {
		addEntry(key, value)
	}
	return value
}

func (l *loader) Get(ki interface{}) dgo.Value {
	key, ok := vf.Value(ki).(dgo.String)
	if !ok {
		return nil
	}
	l.lock.RLock()
	v := l.entries.Get(key)
	l.lock.RUnlock()
	if v == nil && l.finder != nil {
		v = vf.Value(l.finder(l, key.String()))
		v = l.add(key, v)
	}
	if vf.Nil == v {
		v = nil
	}
	return v
}

func (l *loader) Load(name string) dgo.Value {
	return load(l, name)
}

func (l *loader) Namespace(name string) dgo.Loader {
	if name == `` {
		return l
	}

	l.lock.RLock()
	nv, ok := l.namespaces.Get(name).(dgo.Native)
	l.lock.RUnlock()

	if ok {
		return nv.GoValue().(dgo.Loader)
	}
	if l.nsCreator == nil {
		return nil
	}

	ns := l.nsCreator(l, name)
	if ns != nil {
		var old dgo.Value

		l.lock.Lock()
		if old = l.namespaces.Get(name); old == nil {
			l.namespaces.Put(name, ns)
		}
		l.lock.Unlock()

		if nil != old {
			// Either the nsCreator did something wrong that resulted in the creation of this
			// namespace or a another one has been created from another go routine.
			if !old.Equals(ns) {
				panic(fmt.Errorf(`namespace %q is already defined`, name))
			}

			// Get rid of the duplicate
			ns = old.(dgo.Loader)
		}
	}
	return ns
}

func (l *loader) NewChild(finder dgo.Finder, nsCreator dgo.NsCreator) dgo.Loader {
	return loaderWithParent(l, finder, nsCreator)
}

func (l *childLoader) Get(key interface{}) dgo.Value {
	v := l.parent.Get(key)
	if v == nil {
		v = l.Loader.Get(key)
	}
	return v
}

func (l *childLoader) Load(name string) dgo.Value {
	return load(l, name)
}

func (l *childLoader) Namespace(name string) dgo.Loader {
	if name == `` {
		return l
	}
	pv := l.parent.Namespace(name)
	v := l.Loader.Namespace(name)
	switch {
	case v == nil:
		v = pv
	case pv == nil:
	default:
		v = &childLoader{Loader: v, parent: pv}
	}
	return v
}

func (l *childLoader) NewChild(finder dgo.Finder, nsCreator dgo.NsCreator) dgo.Loader {
	return loaderWithParent(l, finder, nsCreator)
}

func loaderWithParent(parent dgo.Loader, finder dgo.Finder, nsCreator dgo.NsCreator) dgo.Loader {
	return &childLoader{Loader: New(parent.ParentNamespace(), parent.Name(), vf.Map(), finder, nsCreator), parent: parent}
}
