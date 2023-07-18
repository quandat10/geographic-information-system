neo4j:
	docker run \
    -p 7474:7474 -p 7687:7687 \
    -v ./data:/data -v ./plugins:/plugins \
    --name neo4j \
    -e NEO4J_apoc_export_file_enabled=true \
    -e NEO4J_apoc_import_file_enabled=true \
    -e NEO4J_apoc_import_file_use__neo4j__config=true \
    -e NEO4JLABS_PLUGINS=\[\"apoc\"\] \
	-e NEO4J_dbms_security_procedures_unrestricted=apoc.*,apoc.spatial.* \
	-e NEO4J_AUTH=neo4j/quandat10 neo4j:4.0
neo4j-restart:
	docker restart neo4j
server:
	cd backend && go run cmd/main.go
client:
	cd frontend && yarn dev
.PHONY: neo4j server backend neo4j-restart
