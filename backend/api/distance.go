package api

import (
	"errors"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/rs/zerolog/log"
)

func (s *Server) listUsersInsideRadius(username string) ([]interface{}, error) {
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
		rs, err := tx.Run(`match (u:User) where u.username<>$username return (u);`, map[string]interface{}{
			"username": username,
		})
		if err != nil {
			log.Error().Msg(err.Error())
			return User{}, err
		}

		for rs.Next() {
			record := rs.Record()
			rs := record.Values[0]
			users = append(users, rs)
		}

		return rs, nil
	})
	return users, nil
}
