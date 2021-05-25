package libs

type Configuration struct {
	Config  Config                 `yaml:"config"`
	Meta    map[string]interface{} `yaml:"meta"`
	Name    string                 `yaml:"name"`
	Channel string                 `yaml:"channel"`
}

// Config
type Config struct {
	Port          string `yaml:"port"`
	Username      string `yaml:"username"`
	Password      string `yaml:"password"`
	Host          string `yaml:"host"`
	Url           string `yaml:"url"`
	Authorization string `yaml:"authorization"`
	From          string `yaml:"from"`
	Accept        string `yaml:"accept"`
}

type Message struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}
