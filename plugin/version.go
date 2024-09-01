package plugin

const TopicVersion = "pufferpanel:version"

//go:generate msgp
type Version struct {
	Version string `msg:"version"`
}
