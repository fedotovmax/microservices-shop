package httphelper

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/fedotovmax/microservices-shop/api-gateway/internal/domain"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func DecodeJSON(r io.Reader, v any) error {
	defer io.Copy(io.Discard, r)
	return json.NewDecoder(r).Decode(v)
}

func WriteJSON(w http.ResponseWriter, status int, v any) {

	w.Header().Set(HeaderContentType, HeaderContentTypeJSON)
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, `{"message": "failed to encode json"}`, http.StatusInternalServerError)
	}
}

func HandleErrorFromGrpc(err error, w http.ResponseWriter) {

	st, ok := status.FromError(err)

	if ok {
		for _, d := range st.Details() {
			switch info := d.(type) {
			case *errdetails.BadRequest:
				WriteJSON(w, http.StatusBadRequest, info)
				return
			}
		}
	}

	WriteJSON(w, http.StatusBadRequest, domain.NewError(err.Error()))
}

/*
OK	200 OK
Canceled	499 Client Closed Request
Unknown	500 Internal Server Error
InvalidArgument	400 Bad Request
DeadlineExceeded	504 Gateway Timeout
NotFound	404 Not Found
AlreadyExists	409 Conflict
PermissionDenied	403 Forbidden
ResourceExhausted	429 Too Many Requests
FailedPrecondition	400 Bad Request
Aborted	409 Conflict
OutOfRange	400 Bad Request
Unimplemented	501 Not Implemented
Internal	500 Internal Server Error
Unavailable	503 Service Unavailable
DataLoss	500 Internal Server Error
Unauthenticated	401 Unauthorized
*/
func GrpcCodeToHttp(grpcCode codes.Code, fallback ...int) int {

	switch grpcCode {
	case codes.OK:
		return http.StatusOK
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.Internal:
		return http.StatusInternalServerError
	default:
		if len(fallback) > 0 {
			return fallback[0]
		}
		return http.StatusInternalServerError
	}
}
