package excel

// Operation represents a dental operation with dates, tooth numbers, and comments
type Operation struct {
	Dates   []string `json:"dates"`
	Numbers []string `json:"numbers"`
	Comment string   `json:"comment"`
}

// Patient represents a patient with implant information
type Patient struct {
	ID                  string               `json:"id"`
	FIO                 string               `json:"fio"`
	ImplantNumber       int                  `json:"implantNumber"`
	Operations          map[string]Operation `json:"operations"`
	ControlHalfYear     string               `json:"controlHalfYear"`
	ControlYear         string               `json:"controlYear"`
	OccupationalHygiene string               `json:"occupationalHygiene"`
}
