package header

type Header struct { // Metadata to generate the preface GABC
	Name           string `json:"name"`            // Name of the preface, default: ""
	InitialStyle   string `json:"initial-style"`   // Style of the initial, default: "0" (plain)
	FontSize       string `json:"fontsize"`        // Font size, default: "12"
	Font           string `json:"font"`            // Font family, default: "OFLSortsMillGoudy"
	Width          string `json:"width"`           // Width of the page, default: "7.3"
	Height         string `json:"height"`          // Height of the page, default: "11.7"
	Clef           string `json:"clef"`            // Line of the C clef, default: "c3"
	ComposedHeader string `json:"composed-string"` // Composed string from all header fields
}

func (h *Header) SetHeaderOptions() {
	h.InitialStyle = "0"

	if h.FontSize == "" {
		h.FontSize = "12"
	}

	if h.Font == "" {
		h.Font = "OFLSortsMillGoudy"
	}

	if h.Width == "" {
		h.Width = "7.3"
	}

	if h.Height == "" {
		h.Height = "11.7"
	}

	if h.Clef != "c3" && h.Clef != "c4" {
		h.Clef = "c3"
	}

	// Compose the header string:
	h.ComposedHeader = "name: " + h.Name + ";\n" + "initial-style: " + "0" + ";\n" + "fontsize: " + h.FontSize + ";\n" + "font: " + h.Font + ";\n" + "width: " + h.Width + ";\n" + "height: " + h.Height + ";\n" + "%%\n(" + h.Clef + ")\n"
}

/*
name: ;
user-notes: ;
commentary: ;
annotation: ;
centering-scheme: english;
%fontsize: 12;
%spacing: vichi;
%font: OFLSortsMillGoudy;
%width: 4.5;
%height: 11;
%%
(c3)
*/
