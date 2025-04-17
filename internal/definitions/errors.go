package definitions

type ErrResponse struct {
	Code    int    `json:"error_code"`
	Message string `json:"error_message"`
}

func (e ErrResponse) Error() string {
	return e.Message
}

var ErrResponseToShort = ErrResponse{100, "the phrase is to short to apply the whole melody"}
var ErrResponseNoMarks = ErrResponse{101, "each sentence must end with a structure mark. See documentation"}
