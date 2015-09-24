package push

import (
	"sync"

	"github.com/michaeldwan/webmaster/context"
)

type Manifest struct {
	*sync.Mutex
	context *context.Context
	entries map[string]*Entry

	objCount   int64
	fileCount  int64
	redirCount int64
}

func newManifest(context *context.Context) *Manifest {
	return &Manifest{
		Mutex:   &sync.Mutex{},
		context: context,
		entries: make(map[string]*Entry),
	}
}

func (m *Manifest) scan() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for file := range buildSource(m.context) {
			m.addSourceFile(file)
		}
	}()

	go func() {
		defer wg.Done()
		for obj := range scanBucket(m.context) {
			m.addDestObject(obj)
		}
	}()

	wg.Wait()

	for _, e := range m.entries {
		e.plan(m.context)
	}
}

func (m *Manifest) entryForKey(key string) *Entry {
	entry, ok := m.entries[key]
	if !ok {
		entry = &Entry{Key: key}
		m.entries[key] = entry
	}
	return entry
}

func (m *Manifest) addSourceFile(file *File) {
	m.Lock()
	defer m.Unlock()
	e := m.entryForKey(file.Key())
	e.Src = file
	if file.IsRedirect() {
		m.redirCount++
	} else {
		m.fileCount++
	}
}

func (m *Manifest) addDestObject(obj *DestObject) {
	m.Lock()
	defer m.Unlock()
	e := m.entryForKey(obj.Key)
	e.Dst = obj
	m.objCount++
}

func (m *Manifest) entriesForOperation(op Operation) []Entry {
	var entries []Entry
	for _, entry := range m.entries {
		if entry.Operation() == op {
			entries = append(entries, *entry)
		}
	}
	return entries
}
