package hermes


import "github.com/BurntSushi/toml"


type HermesConfig struct {
	DatabaseName	string `toml:"database_path"`

	LoggingDir  	string `toml:"logging_dir"`
	LoggingFormat 	string `toml:"logging_format"`

	Workers      	int    `toml:"workers"`
}


func LoadConfig(path string) (*HermesConfig, error) {
	var cfg HermesConfig
	_, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
