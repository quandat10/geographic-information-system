package api

import (
	"errors"
	"net/http"
	"quandat10/htttdl/backend/utils"
	"strconv"

	"github.com/labstack/echo"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/rs/zerolog/log"
)

type User struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Radius    int64   `json:"radius"`
	Status    bool    `json:"status"`
	Distance  float64 `json:"distance"`
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
		return s.getUser(tx, user.Name)
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
		query := "CREATE (:User {name: $name, longitude: $longitude, latitude: $latitude,radius: $radius,status: $status})"
		parameters := map[string]interface{}{
			"name":      user.Name,
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
	radiusParam := c.QueryParam("radius")
	nameParam := c.QueryParam("name")

	session := s.store.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})

	query := `
		MATCH (a:User {name: $name})-[r:LOCATED_AT]->(aLoc:Location)
		WITH point({latitude: aLoc.latitude, longitude: aLoc.longitude}) AS aPoint
		MATCH (b:User)-[s:LOCATED_AT]->(bLoc:Location)
		WITH b, point({latitude: bLoc.latitude, longitude: bLoc.longitude}) AS bPoint, aPoint, bLoc
		WITH b, bLoc, distance(aPoint, bPoint) AS distance
		WHERE distance < $radius // 10km in meters
		RETURN b.name AS name, bLoc.latitude AS latitude, bLoc.longitude AS longitude, distance
		ORDER BY distance
	`

	radius, err := strconv.ParseInt(radiusParam, 10, 16)
	if err != nil {
		return errors.New(err.Error())
	}

	params := map[string]interface{}{
		"name":   nameParam,
		"radius": int16(radius),
	}

	result, err := session.Run(query, params)
	if err != nil {
		return errors.New(err.Error())
	}

	// Parse the query results into a slice of User structs
	var users []User
	for result.Next() {
		record := result.Record()

		userDistance := record.Values

		user := User{
			Name:      userDistance[0].(string),
			Latitude:  userDistance[1].(float64),
			Longitude: userDistance[2].(float64),
			Distance:  userDistance[3].(float64),
		}
		users = append(users, user)
	}

	return c.JSON(http.StatusOK, &utils.ResponseMsg{
		Status: "success",
		Data:   users,
	})
}

func (s *Server) updateStore(user User) error {
	session := s.store.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})

	// Check User Exist
	userData, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return s.getUser(tx, user.Name)
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
		query := "MATCH (u:User {name: $name}) SET u.status=$status,u.longitude=$longitude,u.latitude=$latitude,u.radius=$radius RETURN u"
		parameters := map[string]interface{}{
			"name":      user.Name,
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

func (s *Server) getUser(tx neo4j.Transaction, name string) (User, error) {
	result, err := tx.Run(
		`MATCH (u:User {name: 'User 1'})-[:LOCATED_AT]->(l:Location)
		RETURN u.name AS name, l.latitude AS latitude, l.longitude AS longitude`,
		map[string]interface{}{
			"name": name,
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

	username, _ := record.Get("name")
	latitude, _ := record.Get("latitude")
	longitude, _ := record.Get("longitude")

	user := User{
		Name:      username.(string),
		Latitude:  latitude.(float64),
		Longitude: longitude.(float64),
	}

	return user, nil
}
