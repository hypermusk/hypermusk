package api

type UseMap struct {
	Map map[string]string
}

type ReservedKeywordsForObjC struct {
	New         string
	Alloc       string
	Copy        string
	MutableCopy string
	Description string
	NormalName  string
}

type Service interface {
	Authorize(name string) (useMap UseMap, err error)
	PermiessionDenied() (err error)
	GetReservedKeywordsForObjC() (r ReservedKeywordsForObjC, err error)
}
