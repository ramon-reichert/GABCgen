package httpGabcGen

import "github.com/ramon-reichert/GABCgen/cmd/internal/gabcGen"

type GabcHandler struct {
	gabcGenAPI gabcGen.GabcGenAPI
	//	requestTimeout time.Duration
}

func NewGabcHandler(gabc gabcGen.GabcGenAPI /*, reqTimeout time.Duration */) GabcHandler {
	return GabcHandler{
		gabcGenAPI: gabc,
		//	requestTimeout: reqTimeout,
	}
}
