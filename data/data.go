package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/cschep/trix"
)

//WeddingData holds all the data associated with a wedding
type WeddingData struct {
	t          *trix.Trix
	InviteList []string
}

//NewWeddingData makes a new WeddingData struct with a trix object initialized
//and InviteList populated.
func NewWeddingData(spreadsheetID string) (*WeddingData, error) {
	t, err := trix.NewTrix(spreadsheetID)
	if err != nil {
		return nil, err
	}

	//read from a file
	inviteList, err := readInviteList()
	if err != nil {
		log.Println("retrieving inviteList from google")
		invites, err := t.Get("RSVP!J:J")
		if err != nil {
			return nil, err
		}

		for _, invite := range invites.Values {
			inviteList = append(inviteList, invite[0].(string))
		}

		err = saveInviteList(inviteList)
		if err != nil {
			log.Println("saving failed", err)
		}
	} else {
		log.Println("using cached inviteList")
	}

	wd := &WeddingData{
		t:          t,
		InviteList: inviteList,
	}

	return wd, nil
}

func saveInviteList(list []string) error {
	b, err := json.Marshal(list)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("invite_list.json", b, 0777)
	if err != nil {
		return err
	}

	return nil
}

func readInviteList() ([]string, error) {
	b, err := ioutil.ReadFile("invite_list.json")
	if err != nil {
		return nil, err
	}

	var inviteList []string
	err = json.Unmarshal(b, &inviteList)
	if err != nil {
		return nil, err
	}

	return inviteList, nil
}

//GetUserList returns a list of people that need to RSVP
// func (wd *WeddingData) GetUserList() []string {
// 	return []string{"James Pozdena", "Ruffenach Family"}
// }

//TODO: refactor, these are basically the same

//RespondNo submits the data as a NO response
func (wd *WeddingData) RespondNo(who string, note string) error {
	// var values [][]interface{}
	// values = append(values, []interface{}{who, note, "NO"})
	// wd.t.InsertRow("RSVP", values)

	readResp, err := wd.t.Get("RSVP!J:J")
	if err != nil || len(readResp.Values) < 1 {
		log.Println("No Values.", err)
		return err
	}

	var writeRow int
	for i, row := range readResp.Values {
		if row[0] == who {
			writeRow = i + 1 //the spreadsheet has a
		}
	}

	updateRange := fmt.Sprintf("RSVP!K%d:Q%d", writeRow, writeRow)

	var values [][]interface{}
	values = append(values, []interface{}{note, "NO"})
	_, err = wd.t.Update(updateRange, values)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

//RespondYes submits the data as a YES response
func (wd *WeddingData) RespondYes(who string, note string) error {
	// var values [][]interface{}
	// values = append(values, []interface{}{who, note, "NO"})
	// wd.t.InsertRow("RSVP", values)

	readResp, err := wd.t.Get("RSVP!J:J")
	if err != nil || len(readResp.Values) < 1 {
		log.Println("No Values.", err)
		return err
	}

	var writeRow int
	for i, row := range readResp.Values {
		if row[0] == who {
			writeRow = i + 1 //the spreadsheet has a
		}
	}

	updateRange := fmt.Sprintf("RSVP!K%d:Q%d", writeRow, writeRow)

	var values [][]interface{}
	values = append(values, []interface{}{note, "YES"})
	_, err = wd.t.Update(updateRange, values)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
