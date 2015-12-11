package staticlib

import (
	"bytes"
	"crypto/md5"
	"sync"
)

type Operation int

const (
	Skip Operation = iota
	Create
	Update
	ForceUpdate
	Delete
)

type Entry struct {
	Key       string
	Src       *File
	Dst       *Object
	Operation Operation
}

type ManifestStats struct {
	ObjCount   int64
	FileCount  int64
	RedirCount int64
}

type Manifest struct {
	*sync.Mutex
	entries map[string]*Entry
	digest  [16]byte
	Stats   *ManifestStats
}

func newManifest() *Manifest {
	return &Manifest{
		Mutex:   &sync.Mutex{},
		entries: make(map[string]*Entry),
		Stats:   &ManifestStats{},
	}
}

func (m *Manifest) plan(forceUpdate bool) {
	for _, e := range m.entries {
		if e.Src != nil && e.Dst == nil {
			e.Operation = Create
		} else if e.Src == nil && e.Dst != nil {
			e.Operation = Delete
		} else if e.Src != nil && e.Dst != nil && !bytes.Equal(e.Src.Digest, e.Dst.Digest) {
			e.Operation = Update
		} else if forceUpdate {
			e.Operation = ForceUpdate
		} else {
			e.Operation = Skip
		}
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

func (m *Manifest) addFile(file File) {
	m.Lock()
	defer m.Unlock()
	e := m.entryForKey(file.Key)
	e.Src = &file
	if file.IsRedirect() {
		m.Stats.RedirCount++
	} else {
		m.Stats.FileCount++
	}
}

func (m *Manifest) addObject(obj Object) {
	m.Lock()
	defer m.Unlock()
	e := m.entryForKey(obj.Key)
	e.Dst = &obj
	m.Stats.ObjCount++
}

func (m *Manifest) entriesForOperations(ops ...Operation) []Entry {
	var entries []Entry
	for _, entry := range m.entries {
		for _, op := range ops {
			if entry.Operation == op {
				entries = append(entries, *entry)
				break
			}
		}
	}
	return entries
}

func (m *Manifest) createDigest() {
	b := bytes.Buffer{}
	for _, entry := range m.entries {
		b.Write([]byte(entry.Key))
		b.Write([]byte(string(entry.Operation)))
	}
	m.digest = md5.Sum(b.Bytes())
}

func (m *Manifest) entryStream() <-chan *Entry {
	out := make(chan *Entry)
	go func() {
		for _, e := range m.entries {
			out <- e
		}
		close(out)
	}()
	return out
}
