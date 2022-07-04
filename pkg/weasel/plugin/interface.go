package plugin

type Plugin interface {
	Enter(string) error
}
