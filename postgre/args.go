package postgre

type NamedArg struct {
	value any
	name  string
}

func Named(name string, value any) NamedArg {
	return NamedArg{name: name, value: value}
}

func (a NamedArg) Name() string {
	return a.name
}

func (a NamedArg) Value() any {
	return a.value
}
