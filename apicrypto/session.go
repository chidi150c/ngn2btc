package apicrypto

import(
	"user-apiv2/model"
	"time"
	"net/http"
)

type Session struct{
	addTdDBChan  chan *Transaction
	getTdDBChan chan getDBData
	deleteTdDBChan chan TransactionID
	worker TransactionService
	now         time.Time
	Userservice UserService
	session     model.Sessioner
}
type getDBData struct{
	id TransactionID
	transChan chan *Transaction
}

func NewSession(uDB UserDBType, tDB  TransDBType, msess model.Sessioner)*Session{
	dBAddTdChan := make(chan *Transaction)
	dBGetTdChan := make(chan getDBData)
	dBDeleteTdChan := make(chan TransactionID)
	go func (){
		for{
			select{
			case td := <-dBAddTdChan:
				tDB[td.ID] = td
			case getd := <-dBGetTdChan:
				if v, ok := tDB[getd.id]; ok{
					getd.transChan <-v
				}else{
					getd.transChan <-nil
				}				  
			case id := <-dBDeleteTdChan:
				delete(tDB, id)
			}
		}
	}()
	s := &Session{
		addTdDBChan: dBAddTdChan,
		getTdDBChan: dBGetTdChan,
		deleteTdDBChan: dBDeleteTdChan,
		session: msess,
	}
	s.worker.SetSession(s) 
	s.Userservice.session = s
	return s
}

func (s *Session) UserFromRequest(r *http.Request) (string, *User, error) {
	red, usr, err := s.session.UserFromRequest(r)
	otherusr := s.Userservice.ModelToOtherUser(usr)
	return red, otherusr, err
}