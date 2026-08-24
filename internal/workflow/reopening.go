package workflow

import (
	"fmt"
	"strings"

	"wedding-sign/internal/model"
	"wedding-sign/internal/service"
)

func OpenDisplaySession(svc *service.Service, sessionID, recordID, profileID string) (model.Session, error) {
	if strings.TrimSpace(sessionID) == "" {
		return model.Session{}, fmt.Errorf("session id is required")
	}
	if _, err := svc.GetRecord(recordID); err != nil {
		return model.Session{}, err
	}
	if _, err := svc.Store().GetProfile(profileID); err != nil {
		return model.Session{}, err
	}
	session := model.NewSession(sessionID, recordID, profileID)
	if err := svc.Store().PutSession(session); err != nil {
		return model.Session{}, err
	}
	return session, nil
}

func UpdateDisplaySession(svc *service.Service, sessionID, screen string) (model.Session, error) {
	session, err := svc.Store().GetSession(sessionID)
	if err != nil {
		return model.Session{}, err
	}
	if !session.SetScreen(screen) {
		return model.Session{}, fmt.Errorf("session %s cannot update screen", sessionID)
	}
	if err := svc.Store().PutSession(session); err != nil {
		return model.Session{}, err
	}
	return session, nil
}

func CloseDisplaySession(svc *service.Service, sessionID string) (model.Session, error) {
	session, err := svc.Store().GetSession(sessionID)
	if err != nil {
		return model.Session{}, err
	}
	session.Close()
	if err := svc.Store().PutSession(session); err != nil {
		return model.Session{}, err
	}
	return session, nil
}
