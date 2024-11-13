package handler

type Handler interface {
}

const (
	VISITOR_HANDLER = iota
	POST_HANDLER
	AUTH_HANDLER
)
