package dispatcher

import "github.com/callevo/ari/arievent"

// Listener type for defining functions as listeners
type Listener func(*arievent.StasisEvent)
