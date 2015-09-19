package context

var current *Config

func Init(path string) {
	current = newConfig()
	current.loadFile(path)
}

func Current() *Config {
	if current == nil {
		Init(ConfigFileName)
	}
	return current
}
