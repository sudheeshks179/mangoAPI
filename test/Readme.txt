Add bird
curl -X POST -v -H "Content-Type: application/json" -d @birdAdd.json http://localhost:8080/birds
curl -X POST -vsb -H "Content-Type: application/json" \
    -d '{   "name":    "mybird", "family":   "mybirdFamily", "continents": ["mybird C1", "mybird C2"],  "added": "2018-01-02", "visible" : false }' \
    http://localhost:8080/birds

Addig error
curl -X POST -v -H "Content-Type: application/json" -d @birdAddError.json http://localhost:8080/birds

List birds
curl -v -H "Content-Type: application/json" http://localhost:8080/birds

Delete Birds
curl -v -X DELETE -H "Content-Type: application/json"  http://localhost:8080/birds/1112

Get bird based on id
curl -v -H "Content-Type: application/json" http://localhost:8080/birds/1111


Mogo DB Help

Drop all documents i collections
db.birds.drop()

To Find
db.birds.find( )

db.birds.find( { visible : true })
