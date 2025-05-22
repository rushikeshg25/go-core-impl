package durablelogs

import (
	"durablelogs/durablelogs/pb"

	"google.golang.org/protobuf/proto"
)

func MustMarshal(message *pb.Log) []byte {
	marshaledLog, err := proto.Marshal(message)
	if err != nil {
		panic(err)
	}
	return marshaledLog
}

func MustUnmarshal(marshaledLog []byte) *pb.Log {
	log := &pb.Log{}
	err := proto.Unmarshal(marshaledLog, log)
	if err != nil {
		panic(err)
	}
	return log
}
