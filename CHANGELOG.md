
<a name="v0.26.0"></a>
## [v0.26.0](https://github.com/aerogear/charmil-host-example/compare/v0.25.0...v0.26.0) (2021-07-22)

### Bug Fixes

* change default pagination flag ([#816](https://github.com/aerogear/charmil-host-example/issues/816))
* remove old workaround for migrating config file name ([#795](https://github.com/aerogear/charmil-host-example/issues/795))
* add owner to registry list ([#802](https://github.com/aerogear/charmil-host-example/issues/802))
* create folder for the initial config ([#806](https://github.com/aerogear/charmil-host-example/issues/806))
* remove URL from table view for serviceregistry list command ([#809](https://github.com/aerogear/charmil-host-example/issues/809))
* invalid location of shared service i18n files ([#808](https://github.com/aerogear/charmil-host-example/issues/808))
* update charmil & validatorOptions ([#814](https://github.com/aerogear/charmil-host-example/issues/814))
* cannot delete service registry by name ([#786](https://github.com/aerogear/charmil-host-example/issues/786))
* **deps:** update module github.com/redhat-developer/app-services-sdk-go to v0.6.0
* **deps:** update module github.com/redhat-developer/app-services-sdk-go to v0.6.2
* **deps:** update module github.com/redhat-developer/app-services-sdk-go to v0.7.0
* **deps:** update all to v0.21.3
* **deps:** update module sigs.k8s.io/controller-runtime to v0.9.3
* **deps:** update all
* **kafka consumer-group:** use group id filter for dynamic completions ([#827](https://github.com/aerogear/charmil-host-example/issues/827))
* **service-account:** reset-credentials prompt ([#838](https://github.com/aerogear/charmil-host-example/issues/838))

### Features

* enable auto completion for output flag in service-registry commands ([#805](https://github.com/aerogear/charmil-host-example/issues/805))
* config command ([#798](https://github.com/aerogear/charmil-host-example/issues/798))
* **consumer-group list:** add flags for pagination ([#821](https://github.com/aerogear/charmil-host-example/issues/821))
* **consumer-group list:** add search flag ([#813](https://github.com/aerogear/charmil-host-example/issues/813))
* **kafka topic:** partitions flag for update ([#823](https://github.com/aerogear/charmil-host-example/issues/823))
* **kafka topic list:** add flags for pagination ([#810](https://github.com/aerogear/charmil-host-example/issues/810))


<a name="v0.25.0"></a>
## [v0.25.0](https://github.com/aerogear/charmil-host-example/compare/v0.24.4...v0.25.0) (2021-07-05)

### Bug Fixes

* fix not working insecure login ([#738](https://github.com/aerogear/charmil-host-example/issues/738))
* add cobra commands validator ([#767](https://github.com/aerogear/charmil-host-example/issues/767))
* **deps:** update module github.com/redhat-developer/app-services-sdk-go to v0.3.7 ([#749](https://github.com/aerogear/charmil-host-example/issues/749))
* **deps:** update module github.com/redhat-developer/app-services-sdk-go to v0.3.6 ([#740](https://github.com/aerogear/charmil-host-example/issues/740))
* **deps:** update module github.com/openconfig/goyang to v0.2.6 ([#737](https://github.com/aerogear/charmil-host-example/issues/737))
* **deps:** update golang.org/x/oauth2 commit hash to 14747e6 ([#741](https://github.com/aerogear/charmil-host-example/issues/741))
* **deps:** update golang.org/x/oauth2 commit hash to bce0382 ([#742](https://github.com/aerogear/charmil-host-example/issues/742))
* **deps:** update golang.org/x/oauth2 commit hash to a8dc77f ([#743](https://github.com/aerogear/charmil-host-example/issues/743))
* **deps:** update module github.com/redhat-developer/app-services-sdk-go to v0.4.0 ([#758](https://github.com/aerogear/charmil-host-example/issues/758))
* **deps:** update all ([#755](https://github.com/aerogear/charmil-host-example/issues/755))
* **deps:** update module github.com/redhat-developer/app-services-sdk-go to v0.3.9 ([#754](https://github.com/aerogear/charmil-host-example/issues/754))
* **deps:** update module github.com/redhat-developer/app-services-sdk-go to v0.3.8 ([#752](https://github.com/aerogear/charmil-host-example/issues/752))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.9.2 ([#751](https://github.com/aerogear/charmil-host-example/issues/751))
* **deps:** update all
* **deps:** update module sigs.k8s.io/controller-runtime to v0.9.1 ([#736](https://github.com/aerogear/charmil-host-example/issues/736))
* **error:** remove blank line from error message ([#727](https://github.com/aerogear/charmil-host-example/issues/727))
* **kafka admin:** validations and show proper error messages ([#729](https://github.com/aerogear/charmil-host-example/issues/729))
* **kafka consumergroup list:** fix reference to the wrong i18n ID ([#775](https://github.com/aerogear/charmil-host-example/issues/775))
* **serviceaccount reset-credentials:** validate serviceaccount ID in prompt ([#720](https://github.com/aerogear/charmil-host-example/issues/720))
* **topic:** use the updated methods in KafkaInstance SDK

### Features

* **kafka topic create:** add cleanup policy flag ([#771](https://github.com/aerogear/charmil-host-example/issues/771))
* **kafka topic update:** add cleanup policy flag ([#776](https://github.com/aerogear/charmil-host-example/issues/776))


<a name="v0.24.4"></a>
## [v0.24.4](https://github.com/aerogear/charmil-host-example/compare/0.24.3...v0.24.4) (2021-06-14)

### Bug Fixes

* fix invalid i18n arg ([#681](https://github.com/aerogear/charmil-host-example/issues/681))
* minor bug fixes and validations ([#696](https://github.com/aerogear/charmil-host-example/issues/696))
* misformatted error message ([#669](https://github.com/aerogear/charmil-host-example/issues/669))
* erroneous message for invalid option ([#689](https://github.com/aerogear/charmil-host-example/issues/689))
* localize id for cluster info ([#688](https://github.com/aerogear/charmil-host-example/issues/688))
* fix invalid i18n ID ([#680](https://github.com/aerogear/charmil-host-example/issues/680))
* **config:** create config directory if it does not exist ([#683](https://github.com/aerogear/charmil-host-example/issues/683))
* **kafka create:** run interactive if no name provided and fix typos ([#705](https://github.com/aerogear/charmil-host-example/issues/705))
* **kafka delete:** show proper message for delete operation ([#700](https://github.com/aerogear/charmil-host-example/issues/700))
* **kafka topic:** update regexp for topic name ([#663](https://github.com/aerogear/charmil-host-example/issues/663))

### Features

* Add ability to specify a binding name using binding-name flag
* **completion:** dynamic completion for suitable flags ([#702](https://github.com/aerogear/charmil-host-example/issues/702))
* **completion:** static completion for suitable flags ([#686](https://github.com/aerogear/charmil-host-example/issues/686))
* **kafka topic:** add search flag to list subcommand ([#709](https://github.com/aerogear/charmil-host-example/issues/709))


<a name="0.24.3"></a>
## [0.24.3](https://github.com/aerogear/charmil-host-example/compare/0.24.2...0.24.3) (2021-04-26)

### Bug Fixes

* fix panic when no kafkas available to selected ([#628](https://github.com/aerogear/charmil-host-example/issues/628))


<a name="0.24.2"></a>
## [0.24.2](https://github.com/aerogear/charmil-host-example/compare/0.24.1...0.24.2) (2021-04-23)

### Bug Fixes

* **iostreams:** make coloured output work on Windows ([#625](https://github.com/aerogear/charmil-host-example/issues/625))


<a name="0.24.1"></a>
## [0.24.1](https://github.com/aerogear/charmil-host-example/compare/0.24.0...0.24.1) (2021-04-20)

### Bug Fixes

* **version:** ignore pre-releases when checking version ([#610](https://github.com/aerogear/charmil-host-example/issues/610))


<a name="0.24.0"></a>
## [0.24.0](https://github.com/aerogear/charmil-host-example/compare/0.23.1...0.24.0) (2021-04-20)

### Features

* use production APIs by default ([#606](https://github.com/aerogear/charmil-host-example/issues/606))

### BREAKING CHANGE


The production environment is now the default environment used. To continue using staging, run `rhoas login --api-gateway=stage --auth-url=stage --mas-auth-url=stage`


<a name="0.23.1"></a>
## [0.23.1](https://github.com/aerogear/charmil-host-example/compare/0.23.0...0.23.1) (2021-04-20)

### Bug Fixes

* **consumergroup:** fix wrong active members value ([#608](https://github.com/aerogear/charmil-host-example/issues/608))
* **login:** clear MAS tokens from config when using token login ([#605](https://github.com/aerogear/charmil-host-example/issues/605))


<a name="0.23.0"></a>
## [0.23.0](https://github.com/aerogear/charmil-host-example/compare/0.22.2...0.23.0) (2021-04-20)

### Bug Fixes

* add support for creating operator based resource ([#599](https://github.com/aerogear/charmil-host-example/issues/599))
* return nil when no Kafka was selected ([#602](https://github.com/aerogear/charmil-host-example/issues/602))
* update pkged file ([#592](https://github.com/aerogear/charmil-host-example/issues/592))
* pointer error ([#588](https://github.com/aerogear/charmil-host-example/issues/588))
* set explicit valid argument number for command ([#585](https://github.com/aerogear/charmil-host-example/issues/585))
* return request output format when list is empty ([#584](https://github.com/aerogear/charmil-host-example/issues/584))
* **completion:** fix dynamic completions for Bash ([#587](https://github.com/aerogear/charmil-host-example/issues/587))

### Features

* show when new version is available ([#598](https://github.com/aerogear/charmil-host-example/issues/598))
* Add bind command using SBO SDK ([#534](https://github.com/aerogear/charmil-host-example/issues/534))
* add consumer group describe command ([#536](https://github.com/aerogear/charmil-host-example/issues/536))
* **consumergroup:** add consumer group commands ([#596](https://github.com/aerogear/charmil-host-example/issues/596))

### BREAKING CHANGE


The `list` commands now return the original response
object in JSON or YAML, instead of nil, depending on the format requested.


<a name="0.22.2"></a>
## [0.22.2](https://github.com/aerogear/charmil-host-example/compare/0.22.1...0.22.2) (2021-04-15)

### Features

* **serviceaccount:** add owner column to table ([#578](https://github.com/aerogear/charmil-host-example/issues/578))


<a name="0.22.1"></a>
## [0.22.1](https://github.com/aerogear/charmil-host-example/compare/0.22.0...0.22.1) (2021-04-14)

### Bug Fixes

* use the OpenShift online terms ([#572](https://github.com/aerogear/charmil-host-example/issues/572))

### Features

* **consumergroup:** add delete command ([#542](https://github.com/aerogear/charmil-host-example/issues/542))
* **topic:** add retention size flag for topic create ([#563](https://github.com/aerogear/charmil-host-example/issues/563))


<a name="0.22.0"></a>
## [0.22.0](https://github.com/aerogear/charmil-host-example/compare/0.21.4...0.22.0) (2021-04-13)

### Bug Fixes

* **topic:** set maximum partition value to 100 ([#560](https://github.com/aerogear/charmil-host-example/issues/560))

### Features

* **auth:** remove double-login for code flow ([#561](https://github.com/aerogear/charmil-host-example/issues/561))
* **topic:** Show 'Unlimited' when value is -1 ([#559](https://github.com/aerogear/charmil-host-example/issues/559))


<a name="0.21.4"></a>
## [0.21.4](https://github.com/aerogear/charmil-host-example/compare/0.21.3...0.21.4) (2021-04-12)

### Bug Fixes

* MAS-SSO token refresh was not enabled ([#558](https://github.com/aerogear/charmil-host-example/issues/558))


<a name="0.21.3"></a>
## [0.21.3](https://github.com/aerogear/charmil-host-example/compare/0.21.2...0.21.3) (2021-04-12)

### Bug Fixes

* use direct link to the operator repository in the status ([#551](https://github.com/aerogear/charmil-host-example/issues/551))
* **serviceaccount:** update regex pattern for description ([#552](https://github.com/aerogear/charmil-host-example/issues/552))


<a name="0.21.2"></a>
## [0.21.2](https://github.com/aerogear/charmil-host-example/compare/0.21.1...0.21.2) (2021-04-09)

### Bug Fixes

* **serviceaccount:** allow capital letters in description ([#550](https://github.com/aerogear/charmil-host-example/issues/550))


<a name="0.21.1"></a>
## [0.21.1](https://github.com/aerogear/charmil-host-example/compare/0.21.0...0.21.1) (2021-04-09)

### Bug Fixes

* update mas-sso url ([#545](https://github.com/aerogear/charmil-host-example/issues/545))
* increase timeout for watching managed kafka to 60 seconds ([#521](https://github.com/aerogear/charmil-host-example/issues/521))

### Features

* **consumergroup:** add consumergroup cmd with list subcommand ([#530](https://github.com/aerogear/charmil-host-example/issues/530))
* **kafka:** add a terms and conditions check ([#529](https://github.com/aerogear/charmil-host-example/issues/529))


<a name="0.21.0"></a>
## [0.21.0](https://github.com/aerogear/charmil-host-example/compare/0.20.6...0.21.0) (2021-04-01)

### Bug Fixes

* switch to new mas-sso url ([#524](https://github.com/aerogear/charmil-host-example/issues/524))

### BREAKING CHANGE


This change will mean that old Kafka instances are inaccessible without overriding the MAS-SSO URL


<a name="0.20.6"></a>
## [0.20.6](https://github.com/aerogear/charmil-host-example/compare/0.20.5...0.20.6) (2021-04-01)

### Bug Fixes

* **topic:** remove partition update code ([#526](https://github.com/aerogear/charmil-host-example/issues/526))


<a name="0.20.5"></a>
## [0.20.5](https://github.com/aerogear/charmil-host-example/compare/0.20.4...0.20.5) (2021-03-31)

### Bug Fixes

* **topic:** set default retention to 7 days ([#516](https://github.com/aerogear/charmil-host-example/issues/516))

### Features

* **kafka:** add interactive prompt for kafka use ([#510](https://github.com/aerogear/charmil-host-example/issues/510))


<a name="0.20.4"></a>
## [0.20.4](https://github.com/aerogear/charmil-host-example/compare/0.20.3...0.20.4) (2021-03-30)

### Bug Fixes

* **cluster:** uniform name for service account ([#517](https://github.com/aerogear/charmil-host-example/issues/517))
* **serviceaccount:** add service account input validation ([#512](https://github.com/aerogear/charmil-host-example/issues/512))


<a name="0.20.3"></a>
## [0.20.3](https://github.com/aerogear/charmil-host-example/compare/0.20.2...0.20.3) (2021-03-29)

### Bug Fixes

* **serviceaccount:** fix invalid i18n message ([#509](https://github.com/aerogear/charmil-host-example/issues/509))
* **serviceaccount reset-credentials:** files should use clientID, clientSecret instead of user, password ([#502](https://github.com/aerogear/charmil-host-example/issues/502))

### Features

* add support for generating modular docs ([#504](https://github.com/aerogear/charmil-host-example/issues/504))


<a name="0.20.2"></a>
## [0.20.2](https://github.com/aerogear/charmil-host-example/compare/0.20.1...0.20.2) (2021-03-26)

### Bug Fixes

* **config:** check if .config directory exists ([#498](https://github.com/aerogear/charmil-host-example/issues/498))
* **kafka topic:** creation in interactive mode should check if name is available ([#492](https://github.com/aerogear/charmil-host-example/issues/492))

### Features

* **kafka create:** add --use flag to set current Kafka instance ([#491](https://github.com/aerogear/charmil-host-example/issues/491))


<a name="0.20.1"></a>
## [0.20.1](https://github.com/aerogear/charmil-host-example/compare/0.20.0...0.20.1) (2021-03-24)

### Bug Fixes

* update kafka admin API client ([#484](https://github.com/aerogear/charmil-host-example/issues/484))
* add Bearer to authorization token ([#480](https://github.com/aerogear/charmil-host-example/issues/480))
* show 500 message from admin server ([#482](https://github.com/aerogear/charmil-host-example/issues/482))
* place the config file in XDG_CONFIG_HOME instead of HOME ([#467](https://github.com/aerogear/charmil-host-example/issues/467))
* lint errors ([#460](https://github.com/aerogear/charmil-host-example/issues/460))
* **serviceaccount create:** display processing text while creation ([#465](https://github.com/aerogear/charmil-host-example/issues/465))
* **topic:** log response body ([#483](https://github.com/aerogear/charmil-host-example/issues/483))

### Features

* add version command ([#471](https://github.com/aerogear/charmil-host-example/issues/471))
* **kafka topic:** display missing columns from topic list ([#466](https://github.com/aerogear/charmil-host-example/issues/466))
* **login:** add flag to skip MAS-SSO login ([#477](https://github.com/aerogear/charmil-host-example/issues/477))
* **status:** display failed_reason for a failing Kafka instance ([#476](https://github.com/aerogear/charmil-host-example/issues/476))


<a name="0.20.0"></a>
## [0.20.0](https://github.com/aerogear/charmil-host-example/compare/0.19.0...0.20.0) (2021-03-15)

### Bug Fixes

* check http response for nil pointer error ([#451](https://github.com/aerogear/charmil-host-example/issues/451))
* appropriate error message when TTY is unavailable for kafka create ([#449](https://github.com/aerogear/charmil-host-example/issues/449))
* removing Managed parts from the CLI ([#448](https://github.com/aerogear/charmil-host-example/issues/448))
* lint error ([#421](https://github.com/aerogear/charmil-host-example/issues/421))
* make binding executable directly in the bash ([#419](https://github.com/aerogear/charmil-host-example/issues/419))
* rename command from info to status in description ([#417](https://github.com/aerogear/charmil-host-example/issues/417))
* **auth:** add dual-login to RH-SSO and MAS-SSO ([#404](https://github.com/aerogear/charmil-host-example/issues/404))
* **serviceaccount create:** allow absolute paths when passing custom file location ([#438](https://github.com/aerogear/charmil-host-example/issues/438))

### Features

* replace --force with --yes
* **kafka topic:** interactive mode for create/update topic ([#436](https://github.com/aerogear/charmil-host-example/issues/436))
* **login:** add the ability to log in using an offline token ([#450](https://github.com/aerogear/charmil-host-example/issues/450))


<a name="0.19.0"></a>
## [0.19.0](https://github.com/aerogear/charmil-host-example/compare/0.18.0...0.19.0) (2021-03-02)

### Bug Fixes

* Add  bindAsFiles by default and enforce proper name for right mo… ([#410](https://github.com/aerogear/charmil-host-example/issues/410))
* invalid oc command for connect operation ([#405](https://github.com/aerogear/charmil-host-example/issues/405))

### Features

* **kafka:** dynamic kafka name completions ([#389](https://github.com/aerogear/charmil-host-example/issues/389))
* **serviceaccount describe:** add describe command ([#406](https://github.com/aerogear/charmil-host-example/issues/406))


<a name="0.18.0"></a>
## [0.18.0](https://github.com/aerogear/charmil-host-example/compare/0.17.2...0.18.0) (2021-02-24)

### Bug Fixes

* Improvements to the CLI to aling with binding format ([#351](https://github.com/aerogear/charmil-host-example/issues/351))
* do not throw error when --force is passed ([#391](https://github.com/aerogear/charmil-host-example/issues/391))
* remove ServiceAuth from Config type ([#369](https://github.com/aerogear/charmil-host-example/issues/369))
* ci: install pkger ([#378](https://github.com/aerogear/charmil-host-example/issues/378))

### Features

* add native asciidoc renderer for docs ([#362](https://github.com/aerogear/charmil-host-example/issues/362))
* **kafka list:** add search flag ([#364](https://github.com/aerogear/charmil-host-example/issues/364))


<a name="0.17.2"></a>
## [0.17.2](https://github.com/aerogear/charmil-host-example/compare/0.17.1...0.17.2) (2021-02-22)

### Bug Fixes

* **i18n:** fix error where locale file not being loaded ([#374](https://github.com/aerogear/charmil-host-example/issues/374))


<a name="0.17.1"></a>
## [0.17.1](https://github.com/aerogear/charmil-host-example/compare/0.17.0...0.17.1) (2021-02-22)

### Bug Fixes

* **login:** fix nil-pointer error ([#373](https://github.com/aerogear/charmil-host-example/issues/373))


<a name="0.17.0"></a>
## [0.17.0](https://github.com/aerogear/charmil-host-example/compare/0.16.0...0.17.0) (2021-02-19)

### Bug Fixes

* invalid YAML
* use yq only if version >= 4 ([#367](https://github.com/aerogear/charmil-host-example/issues/367))
* i18n errors ([#353](https://github.com/aerogear/charmil-host-example/issues/353))
* service account i18n ([#344](https://github.com/aerogear/charmil-host-example/issues/344))

### Features

* **kafka topic:** add topic commands ([#309](https://github.com/aerogear/charmil-host-example/issues/309))
* **whoami:** add whoami command ([#356](https://github.com/aerogear/charmil-host-example/issues/356))


<a name="0.16.0"></a>
## [0.16.0](https://github.com/aerogear/charmil-host-example/compare/0.15.1...0.16.0) (2021-02-10)

### Bug Fixes

* add ability to force delete ([#329](https://github.com/aerogear/charmil-host-example/issues/329))
* refresh token if no access token is provided ([#326](https://github.com/aerogear/charmil-host-example/issues/326))
* **kafka delete:** confirm name only to delete ([#321](https://github.com/aerogear/charmil-host-example/issues/321))

### Features

* **kafka create:** use a positional argument for Kafka create ([#330](https://github.com/aerogear/charmil-host-example/issues/330))


<a name="0.15.1"></a>
## [0.15.1](https://github.com/aerogear/charmil-host-example/compare/0.15.0...0.15.1) (2021-02-04)

### Bug Fixes

* **kafka delete:** add async=true to ensure Kafka can be deleted ([#314](https://github.com/aerogear/charmil-host-example/issues/314))
* **kafka topic:** change topic command to singular form ([#308](https://github.com/aerogear/charmil-host-example/issues/308))


<a name="0.15.0"></a>
## [0.15.0](https://github.com/aerogear/charmil-host-example/compare/0.14.1...0.15.0) (2021-01-28)

### Bug Fixes

* handle "MGD-SERV-API-36" error code ([#305](https://github.com/aerogear/charmil-host-example/issues/305))

### Features

* **status:** add root-level status command ([#301](https://github.com/aerogear/charmil-host-example/issues/301))


<a name="0.14.1"></a>
## [0.14.1](https://github.com/aerogear/charmil-host-example/compare/0.14.0...0.14.1) (2021-01-28)

### Bug Fixes

* print only single topics ([#300](https://github.com/aerogear/charmil-host-example/issues/300))


<a name="0.14.0"></a>
## [0.14.0](https://github.com/aerogear/charmil-host-example/compare/0.13.2...0.14.0) (2021-01-26)

### Bug Fixes

* remove unused function ([#275](https://github.com/aerogear/charmil-host-example/issues/275))
* BootstrapServerHost nil pointer ([#269](https://github.com/aerogear/charmil-host-example/issues/269))
* refactor cluster connect to use new format of the CRD's ([#247](https://github.com/aerogear/charmil-host-example/issues/247))
* **cluster info:** rename command info to status ([#289](https://github.com/aerogear/charmil-host-example/issues/289))
* **connection:** only refresh tokens when needed ([#274](https://github.com/aerogear/charmil-host-example/issues/274))
* **docs:** remove the docs command ([#267](https://github.com/aerogear/charmil-host-example/issues/267))

### Features

* standardise colors for printing to console ([#291](https://github.com/aerogear/charmil-host-example/issues/291))
* **login page:** use Patternfly empty state template ([#292](https://github.com/aerogear/charmil-host-example/issues/292))


<a name="0.13.2"></a>
## [0.13.2](https://github.com/aerogear/charmil-host-example/compare/0.13.1...0.13.2) (2021-01-21)

### Bug Fixes

* pointer error when bootstrap host is empty ([#266](https://github.com/aerogear/charmil-host-example/issues/266))


<a name="0.13.1"></a>
## [0.13.1](https://github.com/aerogear/charmil-host-example/compare/0.13.0...0.13.1) (2021-01-21)

### Bug Fixes

* **status:** fix pointer error ([#262](https://github.com/aerogear/charmil-host-example/issues/262))


<a name="0.13.0"></a>
## [0.13.0](https://github.com/aerogear/charmil-host-example/compare/0.12.0...0.13.0) (2021-01-21)

### Bug Fixes

* negate flag value check ([#254](https://github.com/aerogear/charmil-host-example/issues/254))

### Features

* **serviceaccount:** add interactive mode for the reset credentials command ([#248](https://github.com/aerogear/charmil-host-example/issues/248))


<a name="0.12.0"></a>
## [0.12.0](https://github.com/aerogear/charmil-host-example/compare/0.11.0...0.12.0) (2021-01-20)

### Bug Fixes

* remove kafka credentials format ([#245](https://github.com/aerogear/charmil-host-example/issues/245))


<a name="0.11.0"></a>
## [0.11.0](https://github.com/aerogear/charmil-host-example/compare/0.10.0...0.11.0) (2021-01-19)

### Bug Fixes

* standardize table output format flag ([#233](https://github.com/aerogear/charmil-host-example/issues/233))
* usused option value ([#231](https://github.com/aerogear/charmil-host-example/issues/231))
* **serviceaccount:** remove ability to force delete service accounts ([#230](https://github.com/aerogear/charmil-host-example/issues/230))

### Features

* **kafka:** require name confirmation ([#227](https://github.com/aerogear/charmil-host-example/issues/227))
* **status:** print Bootstrap URL ([#235](https://github.com/aerogear/charmil-host-example/issues/235))


<a name="0.10.0"></a>
## [0.10.0](https://github.com/aerogear/charmil-host-example/compare/0.9.3...0.10.0) (2021-01-14)

### Bug Fixes

* **topics:** missing connection option ([#223](https://github.com/aerogear/charmil-host-example/issues/223))

### Features

* add service account CRUD commands ([#216](https://github.com/aerogear/charmil-host-example/issues/216))


<a name="0.9.3"></a>
## [0.9.3](https://github.com/aerogear/charmil-host-example/compare/0.9.2...0.9.3) (2021-01-11)

### Bug Fixes

* pointer error when bootstrap host is empty ([#214](https://github.com/aerogear/charmil-host-example/issues/214))

### Features

* **login:** add ability to provide custom openid scope ([#210](https://github.com/aerogear/charmil-host-example/issues/210))


<a name="0.9.2"></a>
## [0.9.2](https://github.com/aerogear/charmil-host-example/compare/0.9.1...0.9.2) (2021-01-05)

### Bug Fixes

* ensure context is cancelled when finished ([#198](https://github.com/aerogear/charmil-host-example/issues/198))


<a name="0.9.1"></a>
## [0.9.1](https://github.com/aerogear/charmil-host-example/compare/0.9.0...0.9.1) (2021-01-05)


<a name="0.9.0"></a>
## [0.9.0](https://github.com/aerogear/charmil-host-example/compare/0.8.0...0.9.0) (2020-12-15)

### Bug Fixes

* do not use a pointer for a slice
* append :443 to BootstrapServerHost ([#176](https://github.com/aerogear/charmil-host-example/issues/176))

### Features

* add insecure data plane ([#127](https://github.com/aerogear/charmil-host-example/issues/127))


<a name="0.8.0"></a>
## [0.8.0](https://github.com/aerogear/charmil-host-example/compare/0.7.1...0.8.0) (2020-12-14)

### Features

* print sso url in login ([#167](https://github.com/aerogear/charmil-host-example/issues/167))


<a name="0.7.1"></a>
## [0.7.1](https://github.com/aerogear/charmil-host-example/compare/0.7.0...0.7.1) (2020-12-14)

### Bug Fixes

* display API error reason ([#164](https://github.com/aerogear/charmil-host-example/issues/164))


<a name="0.7.0"></a>
## [0.7.0](https://github.com/aerogear/charmil-host-example/compare/0.6.0...0.7.0) (2020-12-11)

### Bug Fixes

* Initial version of SASL/Plain support for topic creation ([#161](https://github.com/aerogear/charmil-host-example/issues/161))
* remove credentials file
* return error ([#159](https://github.com/aerogear/charmil-host-example/issues/159))
* list command with pagination ([#156](https://github.com/aerogear/charmil-host-example/issues/156))


<a name="0.6.0"></a>
## [0.6.0](https://github.com/aerogear/charmil-host-example/compare/0.5.0...0.6.0) (2020-12-10)

### Bug Fixes

* pandoc trying to remove twice ([#152](https://github.com/aerogear/charmil-host-example/issues/152))
* bump version to 0.6.0
* navigation for cli documentation ([#150](https://github.com/aerogear/charmil-host-example/issues/150))
* remove trailing % from stdout/stderr messages ([#147](https://github.com/aerogear/charmil-host-example/issues/147))


<a name="0.5.0"></a>
## [0.5.0](https://github.com/aerogear/charmil-host-example/compare/0.4.0...0.5.0) (2020-12-10)

### Bug Fixes

* change default client ID and remove token login ([#146](https://github.com/aerogear/charmil-host-example/issues/146))


<a name="0.4.0"></a>
## [0.4.0](https://github.com/aerogear/charmil-host-example/compare/0.3.0...0.4.0) (2020-12-09)

### Bug Fixes

* CR name in credentials
* adding kuberentes secret as output ([#138](https://github.com/aerogear/charmil-host-example/issues/138))
* rename kafka cluster to kafka instance ([#144](https://github.com/aerogear/charmil-host-example/issues/144))

### Features

* refactor connect to use top level group ([#139](https://github.com/aerogear/charmil-host-example/issues/139))
* auto-use kafka cluster after creation ([#142](https://github.com/aerogear/charmil-host-example/issues/142))


<a name="0.3.0"></a>
## [0.3.0](https://github.com/aerogear/charmil-host-example/compare/0.2.0...0.3.0) (2020-12-08)

### Bug Fixes

* unused flag for linting
* make create work ([#133](https://github.com/aerogear/charmil-host-example/issues/133))
* update branch
* add -n flag for create ([#119](https://github.com/aerogear/charmil-host-example/issues/119))
* Make CR using namespaced scope ([#116](https://github.com/aerogear/charmil-host-example/issues/116))
* Rename cr version ([#113](https://github.com/aerogear/charmil-host-example/issues/113))
* change apiversion for connect command
* parse API URL to get host and scheme ([#106](https://github.com/aerogear/charmil-host-example/issues/106))
* remove trailing slash from url ([#103](https://github.com/aerogear/charmil-host-example/issues/103))
* make auth url hard-coded ([#102](https://github.com/aerogear/charmil-host-example/issues/102))
* add missing builders file
* Cleanup of the documentation topics

### Features

* wip: validate kafka name ([#131](https://github.com/aerogear/charmil-host-example/issues/131))
* token-based login ([#132](https://github.com/aerogear/charmil-host-example/issues/132))
* update OPENAPI spec for Service Account ([#121](https://github.com/aerogear/charmil-host-example/issues/121))
* expanded help for credentials command ([#120](https://github.com/aerogear/charmil-host-example/issues/120))
* allow using the currently selected Kafka cluster in the describe command ([#114](https://github.com/aerogear/charmil-host-example/issues/114))
* show message on login success
* rhoas kafka connect command ([#85](https://github.com/aerogear/charmil-host-example/issues/85))
* **cmd:** add YAML output format


<a name="0.2.0"></a>
## [0.2.0](https://github.com/aerogear/charmil-host-example/compare/0.1.0...0.2.0) (2020-11-20)


<a name="0.1.0"></a>
## 0.1.0 (2020-11-18)

### Bug Fixes

* cleanup commands documents for usability ([#69](https://github.com/aerogear/charmil-host-example/issues/69))
* add basic documentation ([#67](https://github.com/aerogear/charmil-host-example/issues/67))
* Remove token mock ([#66](https://github.com/aerogear/charmil-host-example/issues/66))
* add missing elements to guide
* apply fedback by [@wtrocki](https://github.com/wtrocki)
* provide script for the provisioning of the clusters
* rename folder
* make credentials file more secure
* Update gomod version
* add authz
* add minor fixes
* add package
* resolve formatting problems
* reorganize script for api updates
* add handy kafka docker compose to the mock
* add release process docs
* resolve confusion around authorization command
* Remove architecture for cli
* update api
* add initial version of goreleaser
* remove operator from the repository
* disable invalid printing for login/logout
* remove function used to test bot
* rename yml file
* general improvements to make file
* reorganization of the structure
* build for mac and linux
* formatting of the status command
* Add dummy test targetr
* remove vendor folder. It should not be used with packages
* add formatting check to PR's
* revert changes for formatting
* openapi make file
* add missing files to client
* move package to root
* minor fixes
* minor changes for the demo
* CMD backbone
* support for help in browser
* Base for the unit and integration tests
* use packge name
* Guide for running this docs
* disable documentation creator
* documentation generator
* Do not require gopath on build
* Use make when building command
* reduce golang versions
* Use golang setup action
* makefile install problem
* build issue with wrong arg
* minor fixes based on the approved spec
* switch to github package name
* Add logout
* minor improvements
* list command
* add error handling
* Improve formatting
* formatting
* name issue
* Add support for credentials
* rename cli
* rename operator
* add demo setup
* change namespace
* format for the cli
* typo
* additional commands and formatting
* command completion
* rename command
* functional operator
* add spec for operator to read config
* remove duplicate
* add extra commands
* website backbone
* Improve commands
* add docusaurus for the demo
* improve deletion script
* support loging flow
* add new info to readme
* mock
* support for all commands
* mock index page
* multi_az to boolean
* support for the create with some missing environment abstraction
* rename client
* rename cli
* build pipeline
* improve architecture
* Initial architecture
* **cmd:** typo in command name
* **kafka:** delete status code results is 204 and not 200;
* **kafka:** stop command execution when user is not loggen in
* **kafka:** change default region to "us-west-1"
* **kafka:** create command returns 202 and always require async=true
* **login:** make staging the default environment and do not require "url"
* **login:** check token expiration before sending request to control plane
* **login:** make token required for now until a proper login flow is figured out

### Features

* positional argument to reference Kafka
* open browser according to OS
* add status command
* add config
* mock server used for the demo purposes
* print kafka instances to table
* Operator using SDK
* OpenAPI generated client
* Openshift CR's
* **cmd:** Display message if there are no clusters ([#45](https://github.com/aerogear/charmil-host-example/issues/45))
* **kafka:** add mocked version of topics command
* **login:** login using the --token flow

