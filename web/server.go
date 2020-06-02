package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	jsonutils "github.com/sayuthisobri/goutils/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

//
// Server - server struct
//
type Server struct {
	*http.Server
	*gin.Engine
	*Router
	Translator ut.Translator
}

//
// NewServer - create new server
//
func NewServer(server *http.Server, handlers ...gin.HandlerFunc) *Server {
	engine := gin.New()
	engine.Use(handlers...)
	s := &Server{Server: server, Engine: engine, Router: NewRouter(engine)}
	eng := en.New()
	uni := ut.New(eng, eng)
	s.Translator, _ = uni.GetTranslator("en")
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validate := v
		_ = enTranslations.RegisterDefaultTranslations(validate, s.Translator)
	}
	s.Router.DI = &jsonutils.J{
		"translator": s.Translator,
	}
	s.Handler = engine
	return s
}

//os.Interrupt, syscall.SIGTERM
func (s Server) InterruptBy(onInterruptFn func(), sig ...os.Signal) {
	c := make(chan os.Signal)
	signal.Notify(c, sig...)
	go func() {
		<-c
		onInterruptFn()
	}()
}

func (s Server) WaitForInterrupt(exitCode int) {
	s.InterruptBy(func() {
		fmt.Printf("Interupted!")
		os.Exit(exitCode)
	}, os.Interrupt, syscall.SIGTERM)
}
