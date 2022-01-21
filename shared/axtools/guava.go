package axtools

type CreateObjectFunc func(id uint64) error

type Guava struct {
	cache  map[uint64]holdObject
	create CreateObjectFunc
	exit   chan bool
}

type holdObject struct {
	obj interface{}
}

func NewGuava(create CreateObjectFunc) *Guava {
	return &Guava{
		cache:  map[uint64]holdObject{},
		create: create,
		exit:   make(chan bool, 1),
	}
}
