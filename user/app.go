package user

import (
	linesHttp "lines/lines/http"
	"lines/user/ingress/http"
)

type UserApp struct {
	http http.UserHttpIngressInterface
}

func NewUserApp() UserApp {
	ingress := http.NewUserHttpIngress(nil)
	return UserApp{
		http: &ingress,
	}
}

func (a *UserApp) Initialise() error { return nil }

func (a *UserApp) RegisterHTTPRoutes(engine linesHttp.HttpEngine) {
	a.http.RegisterRoutes(engine)
}

func (a *UserApp) RegisterGRPCServices() error {
	return nil
}
