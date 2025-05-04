package httpGabcGen

import "github.com/ramon-reichert/GABCgen/cmd/internal/gabcGen"

type GabcHandler struct {
	gabcGenAPI gabcGen.GabcGen
	//	requestTimeout time.Duration
}

func NewGabcHandler(gabc gabcGen.GabcGen /*, reqTimeout time.Duration */) GabcHandler {
	return GabcHandler{
		gabcGenAPI: gabc,
		//	requestTimeout: reqTimeout,
	}
}
