
## Query samples
### Create users
```
WITH [
  {name: 'User 1', latitude: 37.7749, longitude: -122.4194},
  {name: 'User 2', latitude: 37.7749, longitude: -122.4194},
  {name: 'User 3', latitude: 37.7749, longitude: -122.4194},
  {name: 'User 4', latitude: 37.7749, longitude: -122.4194},
  {name: 'User 5', latitude: 37.7749, longitude: -122.4194},
  {name: 'User 6', latitude: 37.7749, longitude: -122.4194},
  {name: 'User 7', latitude: 37.7749, longitude: -122.4194},
  {name: 'User 8', latitude: 37.7749, longitude: -122.4194},
  {name: 'User 9', latitude: 37.7749, longitude: -122.4194},
  {name: 'User 10', latitude: 37.7749, longitude: -122.4194}
] AS userData
FOREACH (data IN userData |
  CREATE (:User {name: data.name})-[:LOCATED_AT]->(:Location {latitude: data.latitude + (RAND() / 10), longitude: data.longitude + (RAND() / 10)})
)
```
### Get users
```
MATCH (u:User)-[:LOCATED_AT]->(l:Location)
RETURN u.name AS name, l.latitude AS latitude, l.longitude AS longitude
```
### Get User by Name
```
MATCH (u:User {name: 'User 1'})-[:LOCATED_AT]->(l:Location)
RETURN u.name AS name, l.latitude AS latitude, l.longitude AS longitude
```
### Get Users by radius
```
MATCH (a:User {name: 'User 1'})-[r:LOCATED_AT]->(aLoc:Location)
WITH point({latitude: aLoc.latitude, longitude: aLoc.longitude}) AS aPoint
MATCH (b:User)-[s:LOCATED_AT]->(bLoc:Location)
WITH aPoint, b, point({latitude: bLoc.latitude, longitude: bLoc.longitude}) AS bPoint, bLoc
WITH aPoint, b, distance(aPoint, bPoint) AS distance, bLoc
WHERE distance < 10000 // 10km in meters
RETURN b.name AS name, bLoc.latitude AS latitude, bLoc.longitude AS longitude,distance as distance
```