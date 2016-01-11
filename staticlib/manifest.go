package staticlib

import (
	"bytes"
	"crypto/md5"
	"sort"
)

type PushAction int

const (
	Skip PushAction = iota
	Create
	Update
	ForceUpdate
	Delete
)

type Entry struct {
	Key        string
	Src        *Content
	Dst        *Object
	PushAction PushAction
}

type Manifest struct {
	entries map[string]*Entry
	digest  [16]byte
}

func NewManifest(src *Source, bucket *Bucket, forceUpdate bool) *Manifest {
	m := &Manifest{}
	keys := sortedKeys(src, bucket)
	m.entries = make(map[string]*Entry, len(keys))
	digestBuffer := bytes.Buffer{}
	for _, key := range keys {
		src := src.ContentForKey(key)
		dst := bucket.objects[key]
		op := Skip

		if src != nil && dst == nil {
			op = Create
		} else if src == nil && dst != nil {
			op = Delete
		} else if src != nil && dst != nil && !bytes.Equal(src.Digest, dst.Digest) {
			op = Update
		} else if forceUpdate {
			op = ForceUpdate
		}

		m.entries[key] = &Entry{
			Key:        key,
			Src:        src,
			Dst:        dst,
			PushAction: op,
		}

		digestBuffer.Write([]byte(key))
		digestBuffer.Write([]byte(string(op)))
	}
	m.digest = md5.Sum(digestBuffer.Bytes())

	return m
}

func (m *Manifest) entriesForPushActions(ops ...PushAction) []Entry {
	var entries []Entry
	for _, entry := range m.entries {
		for _, op := range ops {
			if entry.PushAction == op {
				entries = append(entries, *entry)
				break
			}
		}
	}
	return entries
}

func sortedKeys(src *Source, bucket *Bucket) []string {
	keyMap := make(map[string]bool)
	for key := range src.contents {
		keyMap[key] = true
	}
	for key := range bucket.objects {
		keyMap[key] = true
	}

	keys := make([]string, 0, len(keyMap))
	for key := range keyMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
