package server

func StartGinServer() error {
	r := SetupRouter()
	if err := r.Run("0.0.0.0:8080"); err != nil {
		return err
	} else {
		return nil
	}
}
