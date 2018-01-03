curl -X POST -vsb -H "Content-Type: application/json" \
    -d '{   "name":    "mybird", "family":   "mybirdFamily", "continents": ["mybird C1", "mybird C2"],  "added": "2018-01-02", "visible" : true }' \
    http://localhost:8080/birds
