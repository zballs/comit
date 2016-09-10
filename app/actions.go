package app

import (
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	// . "github.com/tendermint/go-crypto"
	lib "github.com/zballs/3ii/lib"
	util "github.com/zballs/3ii/util"
	"log"
)

type ActionListener struct {
	*socketio.Server
}

func CreateActionListener() (ActionListener, error) {
	server, err := socketio.NewServer(nil)
	return ActionListener{server}, err
}

// Change print statements to socket emit statements

func (al ActionListener) Run(app *Application) {
	al.On("connection", func(so socketio.Socket) {
		log.Println("connected")

		// Create Accounts
		so.On("create-account", func(passphrase string) {
			pubkey, privkey, err := app.account_manager.CreateAccount(passphrase)
			if err != nil {
				log.Println(err.Error())
			}
			so.Emit("return-keys", util.PubKeyToString(pubkey), util.PrivKeyToString(privkey))
		})

		// Remove Accounts
		so.On("remove-account", func(pubKeyString string, privKeyString string) {
			err := app.account_manager.RemoveAccount(pubKeyString, privKeyString)
			if err != nil {
				log.Println(err.Error())
			} else {
				log.Println(fmt.Sprintf("removed account [PubKeyEd25519{%v}, PrivKeyEd25519{%v}]", pubKeyString, privKeyString))
			}
		})

		// Submit Forms
		so.On("submit-form", func(_type string, _address string, _description string, _specfield string, _pubkey string, _privkey string) {
			str := lib.SERVICE.WriteType(_type) + lib.SERVICE.WriteAddress(_address) + lib.SERVICE.WriteDescription(_description) + lib.SERVICE.WriteSpecField(_specfield, _type) + lib.SERVICE.WritePubkeyString(_pubkey) + lib.SERVICE.WritePrivkeyString(_privkey)
			result := app.account_manager.SubmitForm(str, app)
			log.Println(result.Log)
		})

		// Query Forms
		so.On("query-form", func(str string) {
			form := app.account_manager.QueryForm(str, app.cache)
			if form != nil {
				log.Println(*form)
			} else {
				log.Println("no form found")
			}
		})

		// Query Resolved Forms
		so.On("query-resolved", func(str string) {
			form := app.account_manager.QueryResolved(str, app.cache)
			if form != nil {
				log.Println(*form)
			} else {
				log.Println("no form found")
			}
		})

		// Disconnect
		al.On("disconnection", func() {
			log.Println("disconnected")
		})
	})

	// Errors
	al.On("error", func(so socketio.Socket, err error) {
		log.Println(err.Error())
	})
}
