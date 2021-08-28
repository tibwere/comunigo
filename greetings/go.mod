module gitlab.com/tibwere/comunigo/greetings

go 1.16

require (
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1
)

replace gitlab.com/tibwere/comunigo/greetings/encoding/proto => ../encoding/proto
