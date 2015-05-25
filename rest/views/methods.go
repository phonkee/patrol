package views

import "net/http"

var (
	methods = map[string]func(Viewer) (http.HandlerFunc, error){
		"GET": func(view Viewer) (http.HandlerFunc, error) {
			if t, ok := view.(GetViewer); ok {
				return t.GET, nil
			}
			return nil, ErrMethodNotFound
		},
		"POST": func(view Viewer) (http.HandlerFunc, error) {
			if t, ok := view.(PostViewer); ok {
				return t.POST, nil
			}
			return nil, ErrMethodNotFound
		},
		"PUT": func(view Viewer) (http.HandlerFunc, error) {
			if t, ok := view.(PutViewer); ok {
				return t.PUT, nil
			}
			return nil, ErrMethodNotFound
		},
		"PATCH": func(view Viewer) (http.HandlerFunc, error) {
			if t, ok := view.(PatchViewer); ok {
				return t.PATCH, nil
			}
			return nil, ErrMethodNotFound
		},
		"DELETE": func(view Viewer) (http.HandlerFunc, error) {
			if t, ok := view.(DeleteViewer); ok {
				return t.DELETE, nil
			}
			return nil, ErrMethodNotFound
		},
		"OPTIONS": func(view Viewer) (http.HandlerFunc, error) {
			if t, ok := view.(OptionsViewer); ok {
				return t.OPTIONS, nil
			}
			return nil, ErrMethodNotFound
		},
		"HEAD": func(view Viewer) (http.HandlerFunc, error) {
			if t, ok := view.(HeadViewer); ok {
				return t.HEAD, nil
			}
			return nil, ErrMethodNotFound
		},
	}
)
