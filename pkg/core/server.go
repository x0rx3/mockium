package core

// Server is interface for managing server lifecycle
type Server interface {
	Start() error
	Stop()
	Restart() error
	Configure() error
	IsRunning() bool
}
