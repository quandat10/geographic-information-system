neo4j:
	docker run -d --name neo4j -p 7474:7474 -p 7687:7687 --env NEO4J_AUTH=neo4j/quandat10 neo4j
server:
	cd backend && go run cmd/main.go
client:
	cd frontend && yarn dev
.PHONY: neo4j server backend