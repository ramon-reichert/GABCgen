package gabcErrors

type ErrResponse struct {
	Code    int    `json:"error_code"`
	Message string `json:"error_message"`
}

func (e ErrResponse) Error() string {
	return e.Message
}

var ErrToShort = ErrResponse{100, "the phrase is to short to apply the whole melody"}
var ErrNoMarks = ErrResponse{101, "each sentence must end with a structure mark. See documentation"}
var ErrNoText = ErrResponse{102, "no incoming text to be parsed"}
