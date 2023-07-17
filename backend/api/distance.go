package api

import (
	"errors"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/rs/zerolog/log"
)

func (s *Server) listUsersInsideRadius(username string, radius int16) ([]interface{}, error) {
	// create session
	session := s.store.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})

	defer func() {
		session.Close()
	}()
	// Check User Exist
	_, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return s.getUser(tx, username)
	})
	//fmt.Println("userData", userData)

	if err != nil {
		if err.Error() == "Result contains no more records" {
			return nil, errors.New("User not exist")
		}
	}

	users := make([]interface{}, 0)

	_, err = session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		rs, err := tx.Run(`MATCH (a:User {name: $name})-[r:LOCATED_AT]->(aLoc:Location)
		WITH point({latitude: aLoc.latitude, longitude: aLoc.longitude}) AS aPoint
		MATCH (b:User)-[s:LOCATED_AT]->(bLoc:Location)
		WITH b, point({latitude: bLoc.latitude, longitude: bLoc.longitude}) AS bPoint, aPoint, bLoc
		WITH b, bLoc, distance(aPoint, bPoint) AS distance
		WHERE distance < $radius
		RETURN b.name AS name, bLoc.latitude AS latitude, bLoc.longitude AS longitude, distance
		ORDER BY distance`, map[string]interface{}{
			"name":   username,
			"radius": radius,
		})
		if err != nil {
			log.Error().Msg(err.Error())
			return User{}, err
		}

		for rs.Next() {
			record := rs.Record()

			userDistance := record.Values

			user := User{
				Name:      userDistance[0].(string),
				Latitude:  userDistance[1].(float64),
				Longitude: userDistance[2].(float64),
				Distance:  userDistance[3].(float64),
			}
			users = append(users, user)
		}

		return rs, nil
	})
	return users, nil
}
