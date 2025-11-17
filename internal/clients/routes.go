package applications

import (
	"net/http"

	"github.com/jacobmiller22/gauth/pkg/clog"
)

type ClientRoutes struct {
	S *ClientService
}

func (rte *ClientRoutes) GetApplications() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func (rte *ClientRoutes) GetApplication() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := clog.FromContext(ctx)
		l.DebugContext(ctx, "get-client-by-id", "client-id", r.PathValue("client-id"))
	}
}

// func GetApplications

// func (app *gauthApp) applications(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case http.MethodGet:
// 		apps, err := app.db.GetApplications()
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
//
// 		data, err := json.Marshal(apps)
//
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
//
// 		if _, err = w.Write(data); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 	case http.MethodPost:
// 		// Load request body
// 		var data database.ApplicationData
// 		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
//
// 		clientId, err := app.db.CreateApplication(data)
//
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
//
// 		w.Write([]byte(clientId))
//
// 	default:
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}
//
// }
//
// func (app *gauthApp) application(w http.ResponseWriter, r *http.Request) {
// 	clientId := r.PathValue("clientId")
// 	switch r.Method {
// 	case http.MethodGet:
// 		client, err := app.db.GetApplication(clientId)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
//
// 		data, err := json.Marshal(client)
//
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
//
// 		if _, err = w.Write(data); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 	case http.MethodPut:
//
// 		// Load request body
// 		var data database.ApplicationData
// 		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
//
// 		if err := app.db.UpdateApplication(clientId, data); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
//
// 	case http.MethodDelete:
// 		if err := app.db.DeleteApplication(clientId); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 	default:
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}
// }
