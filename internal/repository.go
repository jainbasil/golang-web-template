package internal

type DbRepository interface {
	Ping() error
	Disconnect() error
}
