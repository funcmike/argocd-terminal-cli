package argocd

type OptionsAuth struct {
	ArgoCDServer    string `json:"argoCDServer"`
	ArgoCDAuthToken string `json:"argoCDAuthToken"`
}

type OptionsApp struct {
	AppName      string `json:"appName"`
	AppNamespace string `json:"appNamespace"`
}

type OptionsK8s struct {
	Namespace string `json:"namespace"`
}
