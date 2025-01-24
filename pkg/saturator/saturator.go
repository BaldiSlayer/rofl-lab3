package saturator

func WithContinue(f func(goOn func())) {
	c := true

	goOn := func() {
		c = true
	}

	for c {
		c = false

		f(goOn)
	}
}

func WithBreak(f func(stop func())) {
	c := true

	stop := func() {
		c = false
	}

	for c {
		f(stop)
	}
}
