package ams

import (
	"context"
	"errors"

	"github.com/aerogear/charmil-host-example/internal/build"
	"github.com/aerogear/charmil-host-example/pkg/api/ams/amsclient"
	"github.com/aerogear/charmil-host-example/pkg/connection"
)

func CheckTermsAccepted(conn connection.Connection) (accepted bool, redirectURI string, err error) {
	termsReview, _, err := conn.API().AccountMgmt().
		ApiAuthorizationsV1SelfTermsReviewPost(context.Background()).
		SelfTermsReview(amsclient.SelfTermsReview{
			EventCode: &build.TermsReviewEventCode,
			SiteCode:  &build.TermsReviewSiteCode,
		}).
		Execute()
	if err != nil {
		return false, "", err
	}

	if !termsReview.GetTermsAvailable() && !termsReview.GetTermsRequired() {
		return true, "", nil
	}

	if !termsReview.HasRedirectUrl() {
		return false, "", errors.New("terms must be signed, but there is no terms URL")
	}

	return false, termsReview.GetRedirectUrl(), nil
}
