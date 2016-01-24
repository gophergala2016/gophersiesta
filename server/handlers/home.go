package handlers
import (
	"github.com/gophergala2016/gophersiesta/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"net/http"
)


func GetHome(c *gin.Context) {

	logo :=  `
.      .__                 __
  _____|__| ____   _______/  |______
 /  ___/  |/ __ \ /  ___/\   __\__  \
 \___ \|  \  ___/ \___ \  |  |  / __ \_
/____  >__|\___  >____  > |__| (____  /
     \/        \/     \/            \/
`
	c.String(http.StatusOK, logo)

}
