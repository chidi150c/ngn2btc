package apiexch

import (
	"errors"
	"fmt"
	"user-apiv2/model"
)

var nextID ExchID

func init() {
	nextID = 1
}

type ExchService struct {
	session *Session
}

//var _ ExchServicer = &ExchService{}

func (g *ExchService) BuyBitcoin(operator *User)      {}
func (g *ExchService) SellBitcoin(operator *User)     {}
func (g *ExchService) WithdrawBitcoin(operator *User) {}
func (g *ExchService) DepositBitcoin(operator *User)  {}

// func (g *ExchService) AddExch(operator *User, exchang *Exch) (ExchID, error) {

// 	if exchang.OwnerName != operator.Username || operator.Level != "Admin" {
// 		return exchang.ID, model.ErrUnauthorized
// 	}
// 	g.session.dbExchData[nextID] = exc
// 	exchang.ID = nextID
// 	nextID++
// 	return exchang.ID, nil
// }

func (g *ExchService) GetExch(operator *User, exchgID ExchID) (*Exch, error) {
	//Early return
	if operator.Username == "" {
		return nil, model.ErrOperatorNameEmpty
	}
	if exchgID <= 0 {
		return nil, model.ErrExchIDRequired
	}

	exchang := g.session.dbExchData[exchgID]
	if exchang.OwnerName != operator.Username {
		return nil, model.ErrUnauthorized
	}
	if exchang == nil {
		return nil, model.ErrExchNotFound
	}
	return exchang, nil
}
func (g *ExchService) DeleteExch(operator *User, exchgID ExchID) error {
	//Early return
	if operator.Username == "" {
		return model.ErrOperatorNameEmpty
	}
	if exchgID <= 0 {
		return model.ErrExchIDRequired
	}

	exchang := g.session.dbExchData[exchgID]
	if exchang.ID != exchgID && exchang.OwnerName != operator.Username {
		return model.ErrUnauthorized
	}
	if exchang == nil {
		return model.ErrExchNotFound
	}

	delete(g.session.dbExchData, exchgID)
	return nil
}

func (g *ExchService) ListExchs(operator *User) ([]*Exch, error) {
	//Early return
	if operator.Username == "" {
		return nil, model.ErrOperatorNameEmpty
	}

	if operator.Role != "Admin" {
		return nil, model.ErrOperatorNotAdmin
	}
	var exchs []*Exch
	for _, b := range g.session.dbExchData {
		exchs = append(exchs, b)
	}
	return exchs, nil
}

func (g *ExchService) UpdateExch(operator *User, exchang *Exch) error {
	//Early return
	if operator.Username == "" {
		return model.ErrOperatorNameEmpty
	}

	// Only allow Admin operator to update Product.
	exchInDB, ok := g.session.dbExchData[exchang.ID]
	if !ok {
		return fmt.Errorf("memory g.session.dbExchData: product not found with ID %v", exchang.ID)
	} else if exchInDB.OwnerToken != operator.token || exchInDB.ID != exchang.ID || exchInDB.OwnerName != operator.Username {
		return model.ErrUnauthorized
		//return fmt.Errorf("memory g.session.dbExchData: Non player not allowed to update Product %v", b.ID)
	}
	if exchang.ID == 0 {
		return errors.New("memory g.session.dbExchData: product with unassigned ID passed into updateProduct")
	}
	g.session.dbExchData[exchang.ID] = exchang
	return nil
}
