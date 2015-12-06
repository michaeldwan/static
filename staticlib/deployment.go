package staticlib

import "sync"

type Deployment struct {
	config   Config
	bucket   Bucket
	source   Source
	Manifest *Manifest
}

func NewDeployment(cfg Config) Deployment {
	d := Deployment{
		config:   cfg,
		source:   newSource(cfg),
		bucket:   newBucket(cfg),
		Manifest: newManifest(),
	}
	return d
}

func (d *Deployment) Compile(forceUpdate bool) <-chan *ManifestStats {
	progress := make(chan *ManifestStats)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		for file := range d.source.process(d.config) {
			d.Manifest.addFile(file)
			progress <- d.Manifest.Stats
		}
		wg.Done()
	}()

	go func() {
		for obj := range d.bucket.Scan() {
			d.Manifest.addObject(obj)
			progress <- d.Manifest.Stats
		}
		wg.Done()
	}()

	go func() {
		wg.Wait()
		d.Manifest.plan(forceUpdate)
		close(progress)
	}()

	return progress
}

func (d *Deployment) Clean() {
	if err := d.source.clean(); err != nil {
		panic(err)
	}
}
