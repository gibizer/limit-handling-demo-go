package limits

import "testing"

func TestPrivateAppLimitIsTheUsersPlan(t *testing.T) {
	for _, plan := range []Limit{FreePlan, DeveloperPlan, OrganizationPlan} {
		user := User{&plan}
		app := user.UploadPrivateApp()
		limit := app.Limit()
		if limit != plan {
			t.Errorf(
				"The private app should have limit %v but got %v", plan, limit)
		}

	}
}

func TestPublicAppHasAnIndependentDefaultLimit(t *testing.T) {
	for _, plan := range []Limit{FreePlan, DeveloperPlan, OrganizationPlan} {
		user := User{&plan}
		app := user.UploadPublicApp()
		limit := app.Limit()
		if limit != PublicAppPlan {
			t.Errorf(
				"The public app should have limit %v but got %v",
				FreePlan, limit)
		}

	}
}

func TestPublicAppCanHaveCustomLimit(t *testing.T) {
	user := User{&FreePlan}
	app := user.UploadPublicApp()
	app.CustomizeLimit(&OrganizationPlan)
	limit := app.Limit()
	if limit != OrganizationPlan {
		t.Errorf(
			"The public app should have limit %v but got %v",
			FreePlan, limit)
	}

}

func TestPublicAppCanOptOutToGetUserLimit(t *testing.T) {
	for _, plan := range []Limit{FreePlan, DeveloperPlan, OrganizationPlan} {
		user := User{&plan}
		app := user.UploadPublicApp()
		app.OptOutFromDefaultLimits()
		limit := app.Limit()
		if limit != plan {
			t.Errorf(
				"The public app should have limit %v but got %v",
				FreePlan, limit)
		}

	}
}

func TestPublicAppCanOputOutFromCustomLimit(t *testing.T) {
	user := User{&OrganizationPlan}
	app := user.UploadPublicApp()
	app.CustomizeLimit(&DeveloperPlan)
	app.OptOutFromDefaultLimits()
	limit := app.Limit()
	if limit != OrganizationPlan {
		t.Errorf(
			"The public app should have limit %v but got %v",
			FreePlan, limit)
	}

}
