package ChirpyDatabase

type Result struct {
	Code int
  Body *any
	Error  error
}

func GetErrorResult(code int, err error) Result {
	return Result {
		Code: code,
    Body: nil,
		Error:  err,
	}
}

func GetOKResult(code int, body any) Result {
  return Result {
    Code: code,
    Body: &body,
    Error: nil,
  }
}
