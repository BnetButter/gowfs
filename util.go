package main


type Result[T any] struct {
	Result T;
	Error error;
}

func (result Result[T]) Unwrap() T {
	if result.Error != nil {
		panic(result.Error.Error());
	}
	return result.Result;
}

func (result Result[T]) Maybe() (T, error) {
	return result.Result, result.Error;
}

func Ok[T any](result T) Result[T] {
	return Result[T]{ Result:result }
}

func Err[T any](err error) Result[T] {
	var val T;
	return Result[T]{
		Result:val,
		Error:err,
	}
}

func Wrap[T any](result T, err error) Result[T] {
	return Result[T]{
		Result: result,
		Error: err,
	}
}

func Ensure[T any](result T, err error) T {
	if err != nil {
		panic(err.Error());
	}
	return result;
}

