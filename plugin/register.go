package plugin

const RegisterVersion = "pufferpanel:register"

//go:generate msgp
type RegisterTopics struct {
	Topics []string `json:"topics"`
}
