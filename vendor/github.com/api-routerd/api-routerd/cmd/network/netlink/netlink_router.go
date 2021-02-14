package netlink

import (
	"net/http"

	"github.com/api-routerd/api-routerd/cmd/network/netlink/address"
	"github.com/api-routerd/api-routerd/cmd/network/netlink/link"
	"github.com/api-routerd/api-routerd/cmd/network/netlink/route"
	"github.com/gorilla/mux"
)

func routerLinkGet(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	linkVar := vars["link"]

	switch r.Method {
	case "GET":

		l := link.Link{
			Link: linkVar,
		}

		err := l.Get(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routerLinkAdd(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		link, err := link.DecodeJSONRequest(r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		err = link.Create()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routerLinkDelete(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":
		link, err := link.DecodeJSONRequest(r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		err = link.Delete()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routerLinkSet(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		link, err := link.DecodeJSONRequest(r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		err = link.Set()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routerGetAddress(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	link := vars["link"]

	switch r.Method {
	case "GET":

		a := address.Address{
			Link: link,
		}

		a.Get(rw)
	}
}

func routerAddAddress(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":

		address, err := address.DecodeJSONRequest(r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		err = address.Add()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func routerDeleteAddress(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":

		address, err := address.DecodeJSONRequest(r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		err = address.Del()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		break
	}
}

func routerAddRoute(rw http.ResponseWriter, r *http.Request) {
	route, err := route.DecodeJSONRequest(r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "POST":
		err = route.Configure()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	case "PUT":
		err = route.Configure()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}
}

func routerDeleteRoute(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":
		route, err := route.DecodeJSONRequest(r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		err = route.DeleteGateWay()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		break
	}
}

func routerGetRoute(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	link := vars["link"]

	switch r.Method {
	case "GET":
		r := route.Route{
			Link: link,
		}

		err := r.Get(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// RegisterRouterNetlink register with mux
func RegisterRouterNetlink(n *mux.Router) {
	// Link
	n.HandleFunc("/link/set", routerLinkSet)
	n.HandleFunc("/link/add", routerLinkAdd)
	n.HandleFunc("/link/delete", routerLinkDelete)
	n.HandleFunc("/link/get/{link}", routerLinkGet)
	n.HandleFunc("/link/get", routerLinkGet)

	// Address
	n.HandleFunc("/address/add", routerAddAddress)
	n.HandleFunc("/address/delete", routerDeleteAddress)
	n.HandleFunc("/address/get", routerGetAddress)
	n.HandleFunc("/address/get/{link}", routerGetAddress)

	// Route
	n.HandleFunc("/route/add", routerAddRoute)
	n.HandleFunc("/route/del", routerDeleteRoute)
	n.HandleFunc("/route/get/{link}", routerGetRoute)
}
