package config

type Config interface {
	IsDefaultConfig() bool
	IsUserConfig() bool

	Unmarshal([]byte) error

	AddPath(string, string, string, string, string) error
	RemovePath(string, string) error
	EditPath(string, string, string, string, string, string, string) error

	Log(string) Config
	LogRemoved(map[string]struct{}) Config // struct for smaller size; only need unique keys
	LogEdited(map[string]string) Config
}
