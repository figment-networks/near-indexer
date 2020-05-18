package pipeline

type chainFunc func() error

func runChain(funcs ...chainFunc) error {
	for _, f := range funcs {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}
