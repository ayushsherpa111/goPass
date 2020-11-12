package kit

// type Item interface {
// }

type ReEnter struct {
	Path     string
	Password []byte
}

type CSVData struct {
	CsvPath string
	Key     string
}

type MatchPass struct {
	Key    string
	KeyVal map[string]interface{}
	Tbl    string
}
