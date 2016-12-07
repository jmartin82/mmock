package notify

import "github.com/jmartin82/mmock/definition"

//Notifier notifies the needed parties
type Notifier interface {
	//Notify the needed parties
	Notify(m *definition.Mock) bool
}
