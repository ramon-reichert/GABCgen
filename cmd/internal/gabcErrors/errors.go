package gabcErrors

type ErrResponse struct {
	Code    int    `json:"error_code"`
	Message string `json:"error_message"`
}

func (e ErrResponse) Error() string {
	return e.Message
}

// BadRequest:
var ErrShortPhrase = ErrResponse{100, "the phrase is to short to apply the whole melody"}
var ErrShortParagraph = ErrResponse{101, "each paragraph must have at least three phrases, not counting the conclusion phrase - which can start the last paragraph"}
var ErrNoText = ErrResponse{102, "no incoming text to be parsed"}
var ErrNoLetters = ErrResponse{103, "non-letter char not attached to any letter"}
var ErrEntryInvalidJSON = ErrResponse{104, "the entry is not a valid JSON: "}
var ErrEntryBlankFields = ErrResponse{105, "this entry field must not be blank: "}

// Context:
var ErrRequestTimeout = ErrResponse{201, "context deadline exceeded"}
