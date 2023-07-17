# GIS
## Required
### Frontend
```
Nextjs
```
### Backend
```
Golang
```
### Database
```
Neo4j (apoc plugin for spatial indexing)
```
## How to Start
```
make neo4j
make server
make client
```

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

## Tree Folders
```
├── backend
│   ├── api
│   │   ├── distance.go
│   │   ├── server.go
│   │   ├── socket.go
│   │   └── user.go
│   ├── cmd
│   │   └── main.go
│   ├── go.mod
│   ├── go.sum
│   └── utils
│       └── helper.go
├── frontend
│   ├── components
│   │   ├── maker-map.tsx
│   │   ├── mapbox-map.tsx
│   │   ├── map-loading-holder.tsx
│   │   ├── modal.tsx
│   │   ├── Sidebar.tsx
│   │   └── world-icon.tsx
│   ├── context
│   │   └── context.tsx
│   ├── lib
│   │   └── map-wrapper.ts
│   ├── next.config.js
│   ├── next-env.d.ts
│   ├── package.json
│   ├── package-lock.json
│   ├── pages
│   │   ├── api
│   │   │   ├── axios-client.ts
│   │   │   └── type.ts
│   │   ├── _app.tsx
│   │   └── index.tsx
│   ├── postcss.config.js
│   ├── public
│   │   ├── avatar-svgrepo-com.svg
│   │   ├── favicon.ico
│   │   ├── hero-image.png
│   │   ├── map-loading-screen.gif
│   │   └── vercel.svg
│   ├── README.md
│   ├── styles
│   │   └── globals.css
│   ├── tailwind.config.js
│   ├── tsconfig.json
│   ├── ultis
│   │   └── popupMap.ts
│   └── yarn.lock
├── Makefile
├── plugins
│   └── apoc.jar
└── README.md
```