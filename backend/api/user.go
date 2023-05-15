package api

import (
	"errors"
	"net/http"
	"quandat10/htttdl/backend/utils"

	"github.com/labstack/echo"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/rs/zerolog/log"
)

type User struct {
	Username  string  `json:"username"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Radius    int64   `json:"radius"`
	Status    bool    `json:"status"`
}

type Node struct {
	Id     int64                  `json:"Id"`
	Labels []string               `json:"Labels"`
	Props  map[string]interface{} `json:"Props"`
}

func (s *Server) NewUser(c echo.Context) error {
	user := User{}

	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	session := s.store.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})

	defer func() {
		session.Close()
	}()

	// Check User Exist
	_, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return s.getUser(tx, user.Username)
	})
	if err != nil {
		if err.Error() == "Result contains more than one record" {
			return c.JSON(http.StatusConflict, &utils.ErrorMsg{
				Message: "User already exits",
			})
		}
	}

	if err == nil {
		return c.JSON(http.StatusConflict, &utils.ErrorMsg{
			Message: "User already exits",
		})
	}

	// Create user
	if _, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := "CREATE (:User {username: $username, longitude: $longitude, latitude: $latitude,radius: $radius,status: $status})"
		parameters := map[string]interface{}{
			"username":  user.Username,
			"longitude": user.Longitude,
			"latitude":  user.Latitude,
			"radius":    user.Radius,
			"status":    true,
		}
		_, err = tx.Run(query, parameters)
		return nil, err
	}); err != nil {
		return err
	}

	if err != nil {
		log.Error().Msg("Create user fail")
	}

	return c.JSON(http.StatusCreated, &utils.ResponseMsg{
		Status: "Create success",
		Data:   user,
	})
}

func (s *Server) UpdateUser(c echo.Context) error {
	user := User{}
	err := c.Bind(&user)
	if err != nil {
		log.Error().Msg(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	err = s.updateStore(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, &utils.ErrorMsg{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, &utils.ResponseMsg{
		Status: "Update success",
		Data:   user,
	})
}

func (s *Server) findUser(c echo.Context) error {
	userName := c.Param("username")
	session := s.store.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})

	userData, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return s.getUser(tx, userName)
	})
	if err != nil {
		if err.Error() == "Result contains no more records" {
			return c.JSON(http.StatusConflict, &utils.ErrorMsg{
				Message: "User does not exist",
			})
		}
	}

	return c.JSON(http.StatusOK, &utils.ResponseMsg{
		Status: "success",
		Data:   userData,
	})
}

func (s *Server) updateStore(user User) error {
	session := s.store.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})

	// Check User Exist
	userData, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return s.getUser(tx, user.Username)
	})
	if err != nil {
		if err.Error() == "Result contains no more records" {
			return errors.New("User does not exist")
		}
	}
	//Modify User
	if user.Latitude == 0.0 {
		user.Latitude = userData.(User).Latitude
	}

	if user.Longitude == 0.0 {
		user.Longitude = userData.(User).Longitude
	}

	if user.Radius == 0.0 {
		user.Radius = userData.(User).Radius
	}

	// Update User
	if _, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := "MATCH (u:User {username: $username}) SET u.status=$status,u.longitude=$longitude,u.latitude=$latitude,u.radius=$radius RETURN u"
		parameters := map[string]interface{}{
			"username":  user.Username,
			"longitude": user.Longitude,
			"latitude":  user.Latitude,
			"radius":    user.Radius,
			"status":    user.Status,
		}
		_, err = tx.Run(query, parameters)
		return nil, err
	}); err != nil {
		return err
	}

	return nil
}

func (s *Server) getUser(tx neo4j.Transaction, username string) (User, error) {
	result, err := tx.Run(
		"MATCH (u:User {username: $username}) RETURN"+
			" u.username AS username, u.latitude AS latitude, "+
			"u.longitude AS longitude, u.radius AS radius",
		map[string]interface{}{
			"username": username,
		},
	)
	if err != nil {
		log.Error().Msg(err.Error())
		return User{}, err
	}
	record, err := result.Single()
	if err != nil {
		log.Error().Msg(err.Error())
		return User{}, err
	}
	name, _ := record.Get("username")
	latitude, _ := record.Get("latitude")
	longitude, _ := record.Get("longitude")
	radius, _ := record.Get("radius")

	user := User{
		Username:  name.(string),
		Latitude:  latitude.(float64),
		Longitude: longitude.(float64),
		Radius:    radius.(int64),
	}

	return user, nil
}
