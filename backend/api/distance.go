package api

import (
	"errors"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/rs/zerolog/log"
)

type Data struct {
	ID       int                 `json:"Id"`
	Users    []string            `json:"Labels"`
	Location map[string]UserInfo `json:"Props"`
}

type UserInfo struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Radius    int64   `json:"radius"`
	Status    bool    `json:"status"`
	Username  string  `json:"username"`
}

func convertInterfaceToData(data interface{}) (*Data, error) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, errors.New("data is not a map")
	}

	id, ok := dataMap["ID"].(int)
	if !ok {
		return nil, errors.New("ID is not an integer")
	}

	users, ok := dataMap["Users"].([]interface{})
	if !ok {
		return nil, errors.New("Users is not a slice")
	}
	userNames := make([]string, len(users))
	for i, u := range users {
		userNames[i] = u.(string)
	}

	locationMap, ok := dataMap["Location"].(map[string]interface{})
	if !ok {
		return nil, errors.New("Location is not a map")
	}
	location := make(map[string]UserInfo, len(locationMap))
	for k, v := range locationMap {
		username := k
		userMap, ok := v.(map[string]interface{})
		if !ok {
			return nil, errors.New("user data is not a map")
		}
		lat, ok := userMap["latitude"].(float64)
		if !ok {
			return nil, errors.New("latitude is not a float64")
		}
		lon, ok := userMap["longitude"].(float64)
		if !ok {
			return nil, errors.New("longitude is not a float64")
		}
		radius, ok := userMap["radius"].(int64)
		if !ok {
			return nil, errors.New("radius is not an int64")
		}
		status, ok := userMap["status"].(bool)
		if !ok {
			return nil, errors.New("status is not a boolean")
		}
		userInfo := UserInfo{
			Latitude:  lat,
			Longitude: lon,
			Radius:    radius,
			Status:    status,
			Username:  username,
		}
		location[username] = userInfo
	}

	return &Data{
		ID:       id,
		Users:    userNames,
		Location: location,
	}, nil
}

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
			//fmt.Println(rs)
			//_, ok := rs.(Data)
			//fmt.Println("ok", ok)
			//ax, err := convertInterfaceToData(rs)
			//if err != nil {
			//	log.Error().Err(err)
			//}
			//fmt.Println("ax", ax)
			users = append(users, rs)
		}

		return rs, nil
	})
	return users, nil
}
