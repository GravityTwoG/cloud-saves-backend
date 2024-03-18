package sessions

import (
	userDTOs "cloud-saves-backend/internal/app/cloud-saves-backend/dto/user"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// this map stores the users sessions. For larger scale applications, you can use a database or cache for this purpose
var sessions = map[string]Session{}

// each session contains the username of the user and the time at which it expires
type Session struct {
	Id string
	UserId uint
	Email string
	Username string
	Role string
	ExpiresAt time.Time
}

// we'll use this method later to determine if the session has expired
func (s Session) isExpired() bool {
	return s.ExpiresAt.Before(time.Now())
}

func Create(user *userDTOs.UserResponseDTO) *Session {
	sessionId := uuid.NewString()
	expiresAt := time.Now().Add(60 * time.Minute)
	
	session := Session{
		Id: sessionId,
		UserId: user.Id,
		Email: user.Email,
		Username: user.Username,
		Role: user.Role,
		ExpiresAt: expiresAt,
	}

	sessions[sessionId] = session

	return &session
}

func Get(sessionId string) (*Session, error) {
	session, exists := sessions[sessionId]
	if !exists {
		return nil, fmt.Errorf("Session not found or expired")
	}

	return &session, nil
}

func Refresh(sessionId string) (string, error) {
	session, exists := sessions[sessionId]
	if !exists || session.isExpired() {
		return "", fmt.Errorf("Session not found or expired")
	}

	newSessionId := uuid.NewString()
	expiresAt := time.Now().Add(60 * time.Minute)
	
	newSession := Session{
		Id: newSessionId,
		UserId: session.UserId,
		Email: session.Email,
		Username: session.Username,
		Role: session.Role,
		ExpiresAt: expiresAt,
	}

	sessions[newSessionId] = newSession

	Delete(sessionId)

	return newSessionId, nil 
}

func Delete(sessionId string) {
	delete(sessions, sessionId)
}