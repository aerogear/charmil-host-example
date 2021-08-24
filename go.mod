module github.com/aerogear/charmil-host-example

go 1.16

require (
	github.com/AlecAivazis/survey/v2 v2.2.16
	github.com/MakeNowJust/heredoc v1.0.0
	github.com/Nerzal/gocloak/v7 v7.11.0
	github.com/aerogear/charmil v0.8.2
	github.com/aerogear/charmil-plugin-example v0.2.2
	github.com/coreos/go-oidc/v3 v3.0.0
	github.com/golang-jwt/jwt/v4 v4.0.0
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/google/go-github v17.0.0+incompatible
	github.com/landoop/tableprinter v0.0.0-20201125135848-89e81fc956e7
	github.com/openconfig/goyang v0.2.8
	github.com/pelletier/go-toml v1.9.3
	github.com/phayes/freeport v0.0.0-20180830031419-95f893ade6f2
	github.com/pkg/errors v0.9.1
	github.com/redhat-developer/app-services-sdk-go/kafkainstance v0.2.0
	github.com/redhat-developer/app-services-sdk-go/kafkamgmt v0.3.2
	github.com/redhat-developer/app-services-sdk-go/registrymgmt v0.1.1
	github.com/redhat-developer/service-binding-operator v0.8.0
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	gitlab.com/c0b/go-ordered-json v0.0.0-20201030195603-febf46534d5a
	golang.org/x/oauth2 v0.0.0-20210819190943-2bc19b11175f
	golang.org/x/text v0.3.7
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.22.0
	k8s.io/apimachinery v0.22.0
	k8s.io/client-go v0.22.0
	sigs.k8s.io/controller-runtime v0.9.6
)

// replace github.com/aerogear/charmil-plugin-example => ../charmil-plugin-example
