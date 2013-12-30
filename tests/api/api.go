package api

type UseMap struct {
	Map map[string]string
}

type Service interface {
	Authorize(name string) (useMap UseMap, err error)
}
