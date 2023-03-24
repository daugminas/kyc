package delivery

func (s *server) RegisterRouter() {
	s.e.POST("/user", s.createUser)
	s.e.GET("/user/:userId", s.getUser)
	s.e.PUT("/user/:userId", s.editUser)
	s.e.DELETE("/user/:userId", s.deleteUser)
	// e.GET("/users", controllers.GetAllUsers)
}
