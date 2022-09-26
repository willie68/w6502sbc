package config

type RAM struct {
	Start  int `yaml:"start"`
	Length int `yaml:"length"`
}

type ROM struct {
	Start  int    `yaml:"start"`
	Length int    `yaml:"length"`
	File   string `yaml:"length"`
}

type P6522 struct {
	Start int `yaml:"start"`
}

type P6551 struct {
	Start int `yaml:"start"`
}
