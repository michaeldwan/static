package push

import (
	"sync"
	"webmaster/context"
)

type DestObject struct {
	Key    string
	Digest []byte
}

type Manifest struct {
	*sync.Mutex
	context *context.Context
	entries map[string]*Entry

	objCount int64
	fileCount int64
	redirCount int64
}

func newManifest(context *context.Context) *Manifest {
	return &Manifest{
		Mutex: &sync.Mutex{},
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
		entries := make([]Entry, 0)
		for _, entry := range m.entries {
			if entry.Operation() == op {
				entries = append(entries, *entry)
			}
		}
		return entries
}


// --------------------

// func (m *Manifest) Inspect() {
// 	fmt.Println(m.entries)
// 	for key, e := range m.entries {
// 		destDigest := ""
// 		if e.Dst != nil {
// 			destDigest = hex.EncodeToString(e.Dst.Digest)
// 		}
// 		localDigest := ""
// 		if e.Src != nil {
// 			localDigest = hex.EncodeToString(e.Src.Digest())
// 		}
// 		fmt.Printf("%s: %s - %s <> %s\n", key, e.Operation(), localDigest, destDigest)
// 	}
// }

// func (m *Manifest) updatedKeys() []string {
// 	keys := make([]string, 0)
// 	for key, entry := range m.entries {
// 		if entry.Operation() == Update {
// 			keys = append(keys, key)
// 		}
// 	}
// 	return keys
// }
