package models

var Conf = new(Config)

type Config struct {
	*AppConfig `yaml:"app"`
}
type AppConfig struct {
	Name      string `json:"name"`
	Mode      string `json:"mode"`
	StartTime string `json:"start_time"`
	MachineID int64  `json:"machine_id"`
	Port      int    `json:"port"`
}
