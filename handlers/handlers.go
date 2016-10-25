package handlers

import (
	"io/ioutil"
	"net/http"
	"strconv"

	vps "github.com/jianqiu/vm-pool-server"
	"github.com/jianqiu/vm-pool-server/controllers"
	"github.com/jianqiu/vm-pool-server/db"
	"github.com/jianqiu/vm-pool-server/handlers/middleware"
	"code.cloudfoundry.org/bbs/models"
	"code.cloudfoundry.org/lager"
	"github.com/gogo/protobuf/proto"
	"github.com/tedsuo/rata"
)

func New(
logger, accessLogger lager.Logger,
db db.DB,
serviceClient vps.ServiceClient,
migrationsDone <-chan struct{},
exitChan chan struct{},
) http.Handler {
	vmController := controllers.NewVirtualGuestController(db, serviceClient)
	vmHandler := NewVirtualGuestHandler(vmController, exitChan)

	emitter := middleware.NewLatencyEmitter(logger)

	actions := rata.Handlers{
		// Tasks
		vps.VMsRoute:   route(emitter.EmitLatency(middleware.LogWrap(logger, accessLogger, vmHandler.VirtualGuests))),
	}

	handler, err := rata.NewRouter(vps.Routes, actions)
	if err != nil {
		panic("unable to create router: " + err.Error())
	}

	return middleware.RequestCountWrap(
		UnavailableWrap(handler,
			migrationsDone,
		),
	)
}

func route(f http.HandlerFunc) http.Handler {
	return f
}

type MessageValidator interface {
	proto.Message
	Validate() error
	Unmarshal(data []byte) error
}

func parseRequest(logger lager.Logger, req *http.Request, request MessageValidator) error {
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Error("failed-to-read-body", err)
		return models.ErrUnknownError
	}

	err = request.Unmarshal(data)
	if err != nil {
		logger.Error("failed-to-parse-request-body", err)
		return models.ErrBadRequest
	}

	if err := request.Validate(); err != nil {
		logger.Error("invalid-request", err)
		return models.NewError(models.Error_InvalidRequest, err.Error())
	}

	return nil
}

func exitIfUnrecoverable(logger lager.Logger, exitCh chan<- struct{}, err *models.Error) {
	if err != nil && err.Type == models.Error_Unrecoverable {
		logger.Error("unrecoverable-error", err)
		select {
		case exitCh <- struct{}{}:
		default:
		}
	}
}

func writeResponse(w http.ResponseWriter, message proto.Message) {
	responseBytes, err := proto.Marshal(message)
	if err != nil {
		panic("Unable to encode Proto: " + err.Error())
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(responseBytes)))
	w.Header().Set("Content-Type", "application/x-protobuf")
	w.WriteHeader(http.StatusOK)

	w.Write(responseBytes)
}
