package routes

import (
	"nanosoft/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(publicRoutes, authenticatedRoutes, adminRoutes *gin.RouterGroup) {
	publicRoutes.POST("/user/register", controllers.Register())
	publicRoutes.POST("/user/login", controllers.Login())
	publicRoutes.GET("/user/refresh-token", controllers.RefreshToken())

	authenticatedRoutes.GET("/user/me", controllers.GetUserInfo())
	authenticatedRoutes.PUT("/user/update-info", controllers.UpdateUserInfo())
	authenticatedRoutes.PUT("/user/update-password", controllers.UpdateUserPassword())

	adminRoutes.GET("/admin/get-all-users", controllers.GetAllUsers())
	adminRoutes.PUT("/admin/update-user-role", controllers.UpdateUserRole())
	adminRoutes.DELETE("/admin/delete-user/:id", controllers.DeleteUser())
}

func ServiceRoutes(publicRoutes, authenticatedRoutes, adminRoutes *gin.RouterGroup) {
	publicRoutes.GET("/service/get-all", controllers.GetAllServices())
	publicRoutes.GET("/service/get-one/:id", controllers.GetOneService())

	adminRoutes.POST("/service/create", controllers.CreateService())
	adminRoutes.PUT("/service/update/:id", controllers.UpdateService())
	adminRoutes.DELETE("/service/delete/:id", controllers.DeleteService())
}

func ProjectRoutes(publicRoutes, authenticatedRoutes, adminRoutes *gin.RouterGroup) {
	publicRoutes.GET("/project/get-all", controllers.GetAllProjects())
	publicRoutes.GET("/project/get-one/:id", controllers.GetOneProject())

	adminRoutes.POST("/project/create", controllers.CreateProject())
	adminRoutes.PUT("/project/update/:id", controllers.UpdateProject())
	adminRoutes.DELETE("/project/delete/:id", controllers.DeleteProject())
}

func RemarkRoutes(publicRoutes, authenticatedRoutes, adminRoutes *gin.RouterGroup) {
	publicRoutes.GET("/remark/get-all", controllers.GetAllRemarks())
	publicRoutes.GET("/remark/get-one/:id", controllers.GetOneRemark())

	adminRoutes.POST("/remark/create", controllers.CreateRemark())
	adminRoutes.PUT("/remark/update/:id", controllers.UpdateRemark())
	adminRoutes.DELETE("/remark/delete/:id", controllers.DeleteRemark())
}

func EmailRoutes(publicRoutes, authenticatedRoutes, adminRoutes *gin.RouterGroup) {
	publicRoutes.POST("/email/create", controllers.CreateEmail())

	adminRoutes.GET("/email/get-all", controllers.GetAllEmails())
	adminRoutes.GET("/email/get-one/:id", controllers.GetOneEmail())
	adminRoutes.DELETE("/email/delete/:id", controllers.DeleteEmail())
}
