package staticlib

type Pusher struct {
	deployment *Deployment
	stats      *PushStats
	err        error
}

func NewPush(d *Deployment) *Pusher {
	return &Pusher{deployment: d, stats: &PushStats{}}
}

func (p *Pusher) Err() error {
	return p.err
}

func (p *Pusher) Stats() PushStats {
	return *p.stats
}

func (p *Pusher) Push(concurrency int, forceUpdate bool, simulate bool) <-chan PushEntryResult {
	out := make(chan PushEntryResult)
	results := make(chan PushEntryResult)
	pool := NewPool(concurrency)

	go func() {
		defer close(out)
		for result := range results {
			out <- result
			p.stats.merge(result.Stats)
			if result.Error != nil {
				p.err = result.Error
				return
			}
		}
	}()

	go func() {
		defer func() {
			pool.Wait()
			close(results)
		}()
		for _, e := range p.deployment.Manifest.entries {
			e := e
			if p.err != nil {
				return
			}
			pool.Run(func() {
				results <- p.pushEntry(e, simulate)
			})
		}
	}()
	return out
}

func (p *Pusher) Invalidate() {
	if len(p.deployment.Manifest.entriesForOperations(Update, ForceUpdate)) == 0 {
		return
	}
	for _, distro := range findDistributionsForOrigin(p.deployment.bucket.WebsiteEndpoint()) {
		distro.invalidate(p.deployment.Manifest)
	}
}

type PushStats struct {
	Created int
	Updated int
	Deleted int
	Skipped int
	Bytes   int64
}

func (s *PushStats) merge(s2 PushStats) {
	s.Bytes += s2.Bytes
	s.Created += s2.Created
	s.Updated += s2.Updated
	s.Deleted += s2.Deleted
	s.Skipped += s2.Skipped
}

type PushEntryResult struct {
	Entry *Entry
	Error error
	Stats PushStats
}

func (p *Pusher) pushEntry(e *Entry, simulate bool) PushEntryResult {
	result := &PushEntryResult{Entry: e}
	switch e.Operation {
	case Create:
		if result.Error = p.deployment.bucket.putFile(e.Src, simulate); result.Error == nil {
			result.Stats = PushStats{Bytes: e.Src.Size(), Created: 1}
		}
	case Update, ForceUpdate:
		if result.Error = p.deployment.bucket.putFile(e.Src, simulate); result.Error == nil {
			result.Stats = PushStats{Bytes: e.Src.Size(), Updated: 1}
		}
	case Delete:
		if result.Error = p.deployment.bucket.deleteKey(e.Key, simulate); result.Error == nil {
			result.Stats = PushStats{Deleted: 1}
		}
	case Skip:
		result.Stats = PushStats{Skipped: 1}
	}
	return *result
}
