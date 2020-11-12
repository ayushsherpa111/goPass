package schema

type CSVData struct {
	CsvPath string
	Key     string
}

type Item interface {
	InsertKeys() []string
	GetTbl() string
	GetId() string
	SetTbl(string)
	InsertVals() []interface{}
}

type ItemGenerator interface {
	Generate(Item) Item
}

type MatchPass struct {
	Key     string
	Pattern []interface{}
	Fields  []string
	KeyVal  map[string]interface{}
	Tbl     string
}
