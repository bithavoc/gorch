package gorch

type ProcedureCall interface {
	Resulting(result interface{}) error
}
