package try

import "fmt"

type TryFunc func(args ...interface{}) (res interface{})

// return a different newErr to change origErr
// return caught true to immediately return newRes and stop catch processing
type CatchFunc func(origErr error) (caught bool, newRes interface{}, newErr error)

// cfs are the catch blocks, tried one at a time
// tf is the function to run safely
// res is desired output result(s)
func Try(tf TryFunc, cfs []CatchFunc, args ...interface{}) (res interface{}, err error) {
	wait := make(chan interface{})
	go func() { // try func
		defer func() { // recover func
			r := recover()
			if origErr, isError := r.(error); isError {
				for _, c := range cfs {
					caught, newRes, newErr := c(origErr)
					if caught {
						wait <- newRes
						return
					} else if newErr != origErr {
						err = newErr
						break
					}
				}
				if err == nil {
					err = origErr
				}
			} else if r != nil { // r isnt an error, make it one
				err = fmt.Errorf("caught: %v", r)
			}
			wait <- nil
		}() // recover func
		wait <- tf(args...)
	}() // try func

	res = <-wait // wait synchronously for try or recover to give us a value
	return
}
