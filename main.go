package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/docker/machine/libmachine"
	"github.com/docker/machine/utils"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"

	_ "github.com/docker/machine/drivers/amazonec2"
	_ "github.com/docker/machine/drivers/digitalocean"
	_ "github.com/docker/machine/drivers/google"
	_ "github.com/docker/machine/drivers/rackspace"
	_ "github.com/docker/machine/drivers/virtualbox"
)

var (
	mcn       *libmachine.Machine
	storePath string
)

type HostResource struct{}

func (h HostResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/hosts").
		Doc("Manage Hosts").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.GET("/").To(h.listHosts).
		// docs
		Doc("list hosts").
		Operation("listHosts").
		Writes(libmachine.Host{}))

	ws.Route(ws.GET("/{name}").To(h.getHost).
		// docs
		Doc("get a host").
		Operation("getHost").
		Param(ws.PathParameter("name", "name of the host").DataType("string")).
		Writes(libmachine.Host{}))

	ws.Route(ws.DELETE("/{name}").To(h.removeHost).
		// docs
		Doc("delete a host").
		Operation("removeHost").
		Param(ws.PathParameter("name", "name of the host").DataType("string")))
	container.Add(ws)
}

// GET http://localhost:8080/hosts
func (h HostResource) listHosts(request *restful.Request, response *restful.Response) {
	hosts, err := mcn.List()
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, err.Error())
		return
	}

	response.WriteEntity(hosts)
}

// GET http://localhost:8080/hosts/1
func (h HostResource) getHost(request *restful.Request, response *restful.Response) {
	name := request.PathParameter("name")

	host, err := mcn.Get(name)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, err.Error())
		return
	}

	response.WriteEntity(host)
}

// DELETE http://localhost:8080/hosts/1
func (h HostResource) removeHost(request *restful.Request, response *restful.Response) {

	//	if err != nil {
	//		response.AddHeader("Content-Type", "text/plain")
	//		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	//		return
	//	}

	response.WriteHeader(http.StatusNoContent)
}

func init() {
	flag.StringVar(&storePath, "storepath", utils.GetBaseDir(), "machine storage path")
}

func main() {
	flag.Parse()

	log.Printf("using store: %s\n", storePath)

	store := libmachine.NewFilestore(storePath, "", "")

	m, err := libmachine.New(store)
	if err != nil {
		log.Fatal(err)
	}

	mcn = m

	wsContainer := restful.NewContainer()
	h := HostResource{}
	h.Register(wsContainer)

	config := swagger.Config{
		WebServices:     wsContainer.RegisteredWebServices(),
		WebServicesUrl:  "http://localhost:8080",
		ApiPath:         "/apidocs.json",
		SwaggerPath:     "/apidocs/",
		SwaggerFilePath: "swagger"}

	swagger.RegisterSwaggerService(config, wsContainer)

	log.Printf("start listening on localhost:8080")
	server := &http.Server{Addr: ":8080", Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
