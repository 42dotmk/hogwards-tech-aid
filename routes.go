package main

func InitializeRoutes() {

	router.GET("/", Home)
	router.GET("/donors/:id", HandleDonor)
	router.GET("/donors", ListDonors)
	router.POST("/donors/:id", UpsertDonor)
	router.PUT("/donors", HandleDonor)
}
