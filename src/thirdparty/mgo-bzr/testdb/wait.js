
// We know the master of the first set (pri=1), but not of the second.
var rs1cfg = {_id: "rs1",
              members: [{_id: 1, host: "127.0.0.1:40011", priority: 1},
                        {_id: 2, host: "127.0.0.1:40012", priority: 0},
                        {_id: 3, host: "127.0.0.1:40013", priority: 0}]}
var rs2cfg = {_id: "rs2",
              members: [{_id: 1, host: "127.0.0.1:40021", priority: 1},
                        {_id: 2, host: "127.0.0.1:40022", priority: 1},
                        {_id: 3, host: "127.0.0.1:40023", priority: 0}]}

rs1a = new Mongo("127.0.0.1:40011").getDB("admin")
rs2a = new Mongo("127.0.0.1:40021").getDB("admin")

function countHealthy(rs) {
    var status = rs.runCommand({replSetGetStatus: 1})
    var count = 0
    if (typeof status.members != "undefined") {
        for (var i = 0; i != status.members.length; i++) {
            count += status.members[i].health
        }
    }
    return count
}

var totalRSMembers = rs1cfg.members.length + rs2cfg.members.length

for (var i = 0; i != 10; i++) {
    var count = countHealthy(rs1a) + countHealthy(rs2a)
    print("Replica sets have", count, "healthy nodes.")
    if (count == totalRSMembers) {
        quit(0)
    }
    sleep(3000)
}

print("Replica sets didn't sync up properly.")
quit(12)
