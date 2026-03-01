package commands

type CommandFactory func() Command

var registry = map[string]CommandFactory{}

func Register(name string, factoryCommand CommandFactory) {
	if _, exists := registry[name]; exists {
		panic("command already registerd: " + name)
	}

	registry[name] = factoryCommand
}

func GetFactory(name string) (CommandFactory, bool) {
	f, ok := registry[name]

	return f, ok
}

func GetAll() map[string]CommandFactory {
	return registry
}
