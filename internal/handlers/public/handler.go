package public

type PublicHandler interface {
}

type PublicHandlerImpl struct {
}

func NewPublicHandlerImpl() *PublicHandlerImpl {
	return &PublicHandlerImpl{}
}
