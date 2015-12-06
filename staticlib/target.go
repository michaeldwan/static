package staticlib

type Target struct {
	Bucket Bucket
}

func newTarget(cfg Config) *Target {
	t := Target{}
	t.Bucket = newBucket(cfg)
	return &t
}
