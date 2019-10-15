package mergeban

type Notifier interface {
	Notify(responseURL string, responseBody string)
}

type slackNotifier struct {
}

func NewSlackNotifier() *slackNotifier {
	return &slackNotifier{}
}

func (n *slackNotifier) Notify(responseURL string, responseBody string) {
}
