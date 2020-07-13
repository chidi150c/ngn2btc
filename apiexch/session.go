package apiexch

import (
	"net/http"
	"time"
	"user-apiv2/model"
)

/*Session has the database handle the services can reference them. By making the
ExchService a non-pointer field we reduce the allocations required when creating
a new session.*/
type Session struct {
	dbExchUser  UDBType
	dbExchData  ExDBType
	Exchservice ExchService

	now         time.Time
	Userservice UserService
	session     model.Sessioner
}

func NewSession(uDB UDBType, exDB ExDBType, msess model.Sessioner) *Session {

	s := &Session{
		dbExchUser: uDB,
		dbExchData: exDB,
		session:    msess,
	}
	s.Exchservice.session = s
	s.Userservice.session = s
	s.Userservice.Userservice = msess.Userservice()
	return s
}

func (s *Session) UserFromRequest(r *http.Request) (string, *User, error) {
	red, usr, err := s.session.UserFromRequest(r)
	otherusr := s.Userservice.ModelToOtherUser(usr)
	s.cachedUser = otherusr
	//s.dbExchUser[otherusr.ID] = s.cachedUser
	return red, s.cachedUser, err
}
