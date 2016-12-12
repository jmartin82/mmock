package vars

type Filler interface {
	Fill(holders []string) map[string]string
}
