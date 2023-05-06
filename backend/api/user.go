package api

import (
	"fmt"
	"net/http"
	"quandat10/htttdl/backend/utils"

	"github.com/labstack/echo"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/rs/zerolog/log"
)

type User struct {
	Username string  `json:"username"`
	Latitude      float32 `json:"latitude"`
	Longitude     float32 `json:"longitude"`
	Radius   int16   `json:"radius"`
	Status   bool    `json:"status"`
}

func (s *Server) NewUser(c echo.Context) error {
	user := User{}

	err := c.Bind(&user)
	if err != nil {
		log.Error().Msg(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	session := s.store.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})

	defer func() {
		err = session.Close()
	}()

	// Check User Exist
	_, err = session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return getUser(tx, user.Username)
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
			"username": user.Username,
			"longitude":     user.Longitude,
			"latitude":      user.Latitude,
			"radius":   user.Radius,
			"status":   true,
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

	session := s.store.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})

	// Check User Exist
	userData, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return getUser(tx, user.Username)
	})
	if err != nil {
		if err.Error() == "Result contains no more records" {
			return c.JSON(http.StatusConflict, &utils.ErrorMsg{
				Message: "User does not exist",
			})
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
			"username": user.Username,
			"longitude":     user.Longitude,
			"latitude":      user.Latitude,
			"radius":   user.Radius,
			"status":   user.Status,
		}
		_, err = tx.Run(query, parameters)
		return nil, err
	}); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, &utils.ResponseMsg{
		Status: "Update success",
		Data:   user,
	})
}

func getUser(tx neo4j.Transaction, username string) (User, error) {
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
	fmt.Println("record", record)
	name, _ := record.Get("username")
	latitude, _ := record.Get("latitude")
	longitude, _ := record.Get("longitude")
	radius, _ := record.Get("radius")

	user := User{
		Username: name.(string),
		Latitude:      float32(latitude.(float64)),
		Longitude:     float32(longitude.(float64)),
		Radius:   int16(radius.(int64)),
	}

	return user, nil
}
