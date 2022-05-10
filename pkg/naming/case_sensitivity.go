package naming

type caseSensitivity int

const (
	CaseInsensitive caseSensitivity = 0
	LowerCase       caseSensitivity = 1
	UpperCase       caseSensitivity = 2
)
