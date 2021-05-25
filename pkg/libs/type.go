package libs

type Configuration struct {
	Config  Config                 `yaml:"config"`
	Meta    map[string]interface{} `yaml:"meta"`
	Name    string                 `yaml:"name"`
	Channel string                 `yaml:"channel"`
}

// Config
type Config struct {
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Url      string `yaml:"url"`
	Token    string `yaml:"token"`
}

type Message struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}
