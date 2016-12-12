package vars

type FunctionOperator interface {
	Operate(holders []string, values map[string]string) map[string]string
	GetFunctionVars([]string) []string
}
