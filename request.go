package nanohttp

type Request struct {
	Method      string
	Path        string
	Version     string
	Headers     map[string]string
	Body        []byte
	QueryParams map[string]string
	PathParams  map[string]string
}
