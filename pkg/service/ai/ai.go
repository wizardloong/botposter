package ai

type AI interface {
	Completion(string) (string, error)
}
