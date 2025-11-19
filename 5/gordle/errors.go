package gordle

// corpusError defines a sentinel error
type corpusError string

func (e corpusError) Error() string {
	return string(e)
}