package worker

type Job struct {
	Id       string `yaml:"id"`
	Input    string `yaml:"input"`
	Output   string `yaml:"output"`
	Type     string `yaml:"type"`
	User     string `yaml:"user"`
	Password string `yaml:"pass"`
	Database string `yaml:"database"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}

type ConfigStructure struct {
	Server string `yaml:"server"`
	Jobs   []Job  `yaml:"jobs"`
}

const JobTypeMySqlDump = "mysqldump"
const JobTypeArchive = "archive"
