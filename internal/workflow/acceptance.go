package workflow

import (
	"fmt"

	"wedding-sign/internal/model"
	"wedding-sign/internal/service"
)

type AcceptanceResult struct {
	Record    model.Record
	Profile   model.Profile
	ViewReady bool
}

func AcceptWelcome(svc *service.Service, profileID, coupleNames, date, venue, image string) (AcceptanceResult, error) {
	profile, err := svc.CreateProfile(profileID, profileID)
	if err != nil {
		return AcceptanceResult{}, err
	}
	record, err := svc.CreateRecord(coupleNames, date, venue, image)
	if err != nil {
		return AcceptanceResult{}, err
	}
	record, err = svc.ValidateRecord(record.ID)
	if err != nil {
		return AcceptanceResult{}, err
	}
	record, err = svc.ConfirmRecord(record.ID, "welcome-desk")
	if err != nil {
		return AcceptanceResult{}, err
	}
	view, err := svc.BuildView(record.ID, profile.ID)
	if err != nil {
		return AcceptanceResult{}, err
	}
	if !view.Layout.Valid() {
		return AcceptanceResult{}, fmt.Errorf("welcome view is incomplete")
	}
	return AcceptanceResult{Record: record, Profile: profile, ViewReady: view.Background.IsAccessible()}, nil
}

func ConfirmedRecord(svc *service.Service, id string) (model.Record, error) {
	record, err := svc.ValidateRecord(id)
	if err != nil {
		return model.Record{}, err
	}
	return svc.ConfirmRecord(record.ID, "workflow")
}
