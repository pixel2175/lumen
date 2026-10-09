package errs

func Catch(err *error, fn ...func(*error)) {
	if r := recover(); r != nil {
		e, ok := r.(error)
		if !ok {
			panic(r)
		}
		*err = e
		if len(fn) > 0 && fn[0] != nil {
			fn[0](err)
		}
	}
}

func TryE(err error) {
	if err != nil {
		panic(err)
	}
}

func TryV[T any](v T, err error, custom ...error) T {
	if err != nil {
		if len(custom) > 0 {
			err = custom[0]
		}
		panic(err)
	}
	return v
}
