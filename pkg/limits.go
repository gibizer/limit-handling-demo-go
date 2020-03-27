package limits

// Limit represents build quotas in the CI system.
type Limit struct {
	// Maximum number of concurrent builds
	ConcurrentBuilds int
	// Maximum number of minutes any build can take
	MaxBuildTime int
	// Maximum number of builds per months, or -1 if unlimited
	BuildsPerMonth int
	// Maximum number of accounts, or -1 if unlimited
	MaxTeamMembers int
}

// NOTE(gibi): These pre-defined limits are hardcoded now but I would put them
// in the DB (out of scope) so eventually they can be changed via an admin API.

// FreePlan is a free build plan with limited quotas
var FreePlan = Limit{
	ConcurrentBuilds: 1,
	MaxBuildTime:     10,
	BuildsPerMonth:   200,
	MaxTeamMembers:   2,
}

// DeveloperPlan is an basic plan
var DeveloperPlan = Limit{
	ConcurrentBuilds: 2,
	MaxBuildTime:     45,
	BuildsPerMonth:   -1,
	MaxTeamMembers:   -1,
}

// OrganizationPlan is a top tier plan
var OrganizationPlan = Limit{
	ConcurrentBuilds: 4,
	MaxBuildTime:     90,
	BuildsPerMonth:   -1,
	MaxTeamMembers:   -1,
}

// PublicAppPlan is an app specific limit for public open source applications
var PublicAppPlan = Limit{
	ConcurrentBuilds: 2,
	MaxBuildTime:     45,
	BuildsPerMonth:   -1,
	MaxTeamMembers:   -1,
}

// app is a base type for different kind of application supported by the CI
// system
// NOTE(gibi): having a base struct might be a bit too much for this simple
// example but I expect that in the future there will be a lot more shared data
// and code between app types.
type app struct {
	owner *User
}

// PrivateApp inherits the limits of its owner
type PrivateApp struct {
	app
}

// Limit represents the app's build limits
// The limit of a private app is inherited from its owner
// NOTE(gibi): return a copy of the Limit so the caller cannot change the
// app's limit via the return value
func (a *PrivateApp) Limit() Limit {
	return *a.owner.plan
}

// PublicApp has its own default limit indenpendent from its owner but the
// owner can opt out from such limit. Also admins can customize the default app
// limit.
type PublicApp struct {
	app
	limit *Limit
}

// NewPublicApp creates a PublicApp
// TODO(gibi): Is there a way to default the limit without adding an explicit
// constructor?
func NewPublicApp(owner *User) *PublicApp {
	return &PublicApp{app: app{owner: owner}, limit: &PublicAppPlan}
}

// Limit represents the app's build limits
// The limit of a public app depends on multiple factors see PublicApp for
// details
func (a *PublicApp) Limit() Limit {
	return *a.limit
}

// CustomizeLimit allows specifying a new limit for the given app
func (a *PublicApp) CustomizeLimit(limit *Limit) {
	// NOTE(gibi): this operation needs escalated privileges (e.g. admin)
	// as the owner of the app cannot set its own app's limits. But
	// authorization is out of scope.
	a.limit = limit
}

// OptOutFromDefaultLimits allows the user to change from the default public
// app limit to the limits of its own plan
func (a *PublicApp) OptOutFromDefaultLimits() {
	a.limit = a.owner.plan
}

// User represents the user of a build system
// Each user has a plan associated that describe how much build such user can
// execute. See Limit for the details of such plan.
type User struct {
	// a Limit instance describing how much build a user can execute
	plan *Limit
}

// UploadPrivateApp creates a private application for the user
func (u *User) UploadPrivateApp() *PrivateApp {
	// TODO(gibi): Is there a nicer way to initialize promoted fields?
	// NOTE(gibi): Listing apps for a user is out of scope so an app ref is
	// not stored in the user now. In a list-app-by-user scenario it might
	// be beneficial to have a bidirectional link between App and User.
	return &PrivateApp{app{owner: u}}
}

// UploadPublicApp creates a public application for the user
func (u *User) UploadPublicApp() *PublicApp {
	// NOTE(gibi): Listing apps for a user is out of scope so an app ref is
	// not stored in the user now. In a list-app-by-user scenario it might
	// be beneficial to have a bidirectional link between App and User.
	return NewPublicApp(u)
}

// NOTE(gibi): we might consider adding an opt out helper here depending on
// the needs of the presentation layer that would call
// app.OptOutFromDefaultLimits() on public apps and raise an error on
// private apps.
