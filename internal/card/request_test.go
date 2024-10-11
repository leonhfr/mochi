package card

var (
	_ Request = &createCardRequest{}
	_ Request = &updateCardRequest{}
	_ Request = &archiveCardRequest{}
	_ Request = &deleteCardRequest{}
)
