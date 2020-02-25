package pufferpanel

type ServerIdResponse struct {
	Id string `json:"id"`
}

type ServerStats struct {
	Cpu    float64 `json:"cpu"`
	Memory float64 `json:"memory"`
}

type ServerLogs struct {
	Epoch int64  `json:"epoch"`
	Logs  string `json:"logs"`
}

type ServerRunning struct {
	Running bool `json:"running"`
}

type ServerData struct {
	Variables map[string]Variable `json:"data"`
}

type ServerDataAdmin struct {
	*Server
}

type PufferdRunning struct {
	Message string `json:"message"`
}
