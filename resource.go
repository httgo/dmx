package dmx

import (
	"net/http"
)

type resource struct {
	http.Handler
	pat string
}

func NewResource(pat string, h http.Handler) *resource {
	return &resource{
		Handler: h,
		pat:     pat,
	}
}

		if ok {
		}
	}
}

	}
}
