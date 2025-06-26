package httpGabcGen

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/ramon-reichert/GABCgen/cmd/internal/gabcErrors"
	"github.com/ramon-reichert/GABCgen/cmd/internal/gabcGen"
	"github.com/ramon-reichert/GABCgen/cmd/internal/preface"
)

type GabcHandler struct {
	gabcGenAPI     gabcGen.GabcGen
	requestTimeout time.Duration
}

func NewGabcHandler(gabc gabcGen.GabcGen, reqTimeout time.Duration) GabcHandler {
	return GabcHandler{
		gabcGenAPI:     gabc,
		requestTimeout: reqTimeout,
	}
}

/* Addresses a call to "/preface" according to the requested action.  */
func (h *GabcHandler) preface(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(h.requestTimeout))
	defer cancel()
	r = r.WithContext(ctx)

	method := r.Method
	switch method {
	case http.MethodPost:
		h.genPreface(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

type PrefaceJSON struct {
	Header   preface.PrefaceHeader `json:"header"`
	Dialogue string                `json:"dialogue"`
	Text     string                `json:"text"`
}

/* Validates the entry, then generates a preface GABC. */
func (h *GabcHandler) genPreface(w http.ResponseWriter, r *http.Request) {
	var prefaceEntry PrefaceJSON
	err := json.NewDecoder(r.Body).Decode(&prefaceEntry)
	if err != nil {
		log.Println(err)
		errR := gabcErrors.ErrResponse{
			Code:    gabcErrors.ErrEntryInvalidJSON.Code,
			Message: gabcErrors.ErrEntryInvalidJSON.Message + err.Error(),
		}
		responseJSON(w, http.StatusBadRequest, errR)
		return
	}

	if prefaceEntry.Text == "" {
		errR := gabcErrors.ErrResponse{
			Code:    gabcErrors.ErrEntryBlankFields.Code,
			Message: gabcErrors.ErrEntryBlankFields.Message + "Text",
		}
		responseJSON(w, http.StatusBadRequest, errR)
		return
	}

	prefaceGABC, err := h.gabcGenAPI.GeneratePreface(r.Context(), entryToPreface(prefaceEntry))
	if err != nil {
		handleError(err, w)
		return
	}

	responseJSON(w, http.StatusOK, prefaceToResponse(prefaceGABC))
}

/*Writes a JSON response into a http.ResponseWriter. */
func responseJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

/*Hydrates the preface object with json tags*/
func prefaceToResponse(p preface.Preface) PrefaceJSON {
	return PrefaceJSON{
		Header:   p.Header,
		Dialogue: string(p.Dialogue),
		Text:     p.Text.ComposedGABC,
	}
}

func entryToPreface(pEntry PrefaceJSON) preface.Preface {
	return preface.Preface{
		Header:   pEntry.Header,
		Dialogue: setDialogueTone(pEntry),
		Text:     preface.PrefaceText{LinedText: pEntry.Text},
	}
}

func setDialogueTone(pEntry PrefaceJSON) preface.Dialogue {
	switch pEntry.Dialogue {
	case "regional":
		return preface.Regional
	default:
		return preface.Solemn
	}
}

func handleError(err error, w http.ResponseWriter) {
	log.Println(err)
	if errors.As(err, &gabcErrors.ErrResponse{}) {
		responseJSON(w, http.StatusBadRequest, err)
		return
	} else if errors.Is(err, context.DeadlineExceeded) {
		responseJSON(w, http.StatusGatewayTimeout, gabcErrors.ErrRequestTimeout)
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
}
