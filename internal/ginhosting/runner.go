package ginhosting

import (
	"github.com/arpanrec/secureserver/internal/appconfig"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func GinRunner(serverHosting appconfig.ApplicationServerConfig) {
	if serverHosting.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(JsonLoggerMiddleware())
	err := r.SetTrustedProxies(nil)
	if err != nil {
		log.Fatalln("Error setting trusted proxies: ", err)
	}
	apiRouter := r.Group("/api")
	log.Println("Starting server on port 8080")
	apiRouterV1 := apiRouter.Group("/v1")
	apiRouterV1.Use(AuthMiddleWare(), NameSpaceMiddleWare())
	apiRouterV1.Match([]string{http.MethodGet, http.MethodPost, http.MethodPut, "LOCK", "UNLOCK"},
		"/tfstate/*any", TfStateHandler())
	apiRouterV1.Match([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		"/files/*any", FileServerHandler())
	apiRouterV1.PUT("/pki/*any", PkiHandler())

	if serverHosting.TlsEnable {
		log.Fatal(r.RunTLS("0.0.0.0"+
			":"+strconv.Itoa(serverHosting.Port),
			serverHosting.TlsCertFile,
			serverHosting.TlsKeyFile))
	}
	log.Fatal(r.Run(":" + strconv.Itoa(serverHosting.Port)))
}
