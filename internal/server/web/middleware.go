package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/seftomsk/abf/internal/access"
	"github.com/seftomsk/abf/internal/limiter"
)

func writeResponse(w http.ResponseWriter, status, msg string, code int) {
	resp := ResponseDTO{
		Status: status,
		Code:   code,
		Msg:    msg,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println(err)
	}
}

func CheckRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") == "" {
			writeResponse(
				w,
				http.StatusText(http.StatusBadRequest),
				"Unsupported Content-Type",
				http.StatusBadRequest)
			return
		}
		if r.Header.Get("Content-Type") != "application/json" {
			writeResponse(
				w,
				http.StatusText(http.StatusBadRequest),
				"Content-Type is not application/json",
				http.StatusBadRequest)
			return
		}
		dec := json.NewDecoder(r.Body)

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			writeResponse(
				w,
				http.StatusText(http.StatusBadRequest),
				"Cannot read the body",
				http.StatusBadRequest)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		var dto RequestDTO
		err = json.Unmarshal(bodyBytes, &dto)
		if err != nil {
			var syntaxError *json.SyntaxError
			var unmarshalTypeError *json.UnmarshalTypeError
			status := http.StatusText(http.StatusBadRequest)
			code := http.StatusBadRequest
			var msg string
			switch {
			// Catch any syntax errors in the JSON and send an error message
			// which interpolates the location of the problem to make it
			// easier for the client to fix.
			case errors.As(err, &syntaxError):
				msg = fmt.Sprintf(
					"Request body contains badly-formed "+
						"JSON (at position %d)", syntaxError.Offset)
			// In some circumstances Decode() may also return an
			// io.ErrUnexpectedEOF error for syntax errors in the JSON. There
			// is an open issue regarding this at
			// https://github.com/golang/go/issues/25956.
			case errors.Is(err, io.ErrUnexpectedEOF):
				msg = fmt.Sprintf("Request body contains" +
					"badly-formed JSON")
			// Catch any type errors, like trying to assign a string in the
			// JSON request body to a int field in our Person struct. We can
			// interpolate the relevant field name and position into the error
			// message to make it easier for the client to fix.
			case errors.As(err, &unmarshalTypeError):
				msg = fmt.Sprintf("Request body contains an invalid"+
					"value for the %q field (at position %d)",
					unmarshalTypeError.Field, unmarshalTypeError.Offset)
			// An io.EOF error is returned by Decode() if the request body is
			// empty.
			case errors.Is(err, io.EOF):
				msg = "Request body must not be empty"
			// Catch the error caused by the request body being too large. Again
			// there is an open issue regarding turning this into a sentinel
			// error at https://github.com/golang/go/issues/30715.
			default:
				log.Println(err)
				msg = http.StatusText(http.StatusInternalServerError)
				status = http.StatusText(http.StatusInternalServerError)
				code = http.StatusInternalServerError
			}
			writeResponse(
				w,
				status,
				msg,
				code)
			return
		}

		// Call decode again, using a pointer to an empty anonymous struct as
		// the destination. If the request body only contained a single JSON
		// object this will return an io.EOF error. So if we get anything else,
		// we know that there is additional data in the request body.
		err = dec.Decode(&struct{}{})
		if !errors.Is(err, io.EOF) {
			writeResponse(
				w,
				http.StatusText(http.StatusBadRequest),
				"Request body must only contain a single JSON object",
				http.StatusBadRequest)
			return
		}

		if len(dto.IP) == 0 {
			writeResponse(
				w,
				http.StatusText(http.StatusBadRequest),
				"field \"ip\" must not be empty",
				http.StatusBadRequest)
			return
		}
		if len(dto.Login) == 0 {
			writeResponse(
				w,
				http.StatusText(http.StatusBadRequest),
				"field \"login\" must not be empty",
				http.StatusBadRequest)
			return
		}
		if len(dto.Password) == 0 {
			writeResponse(
				w,
				http.StatusText(http.StatusBadRequest),
				"field \"password\" must not be empty",
				http.StatusBadRequest)
			return
		}

		next(w, r)
	}
}

func BlackAndWhite(a *access.IPAccess, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			writeResponse(
				w,
				http.StatusText(http.StatusBadRequest),
				"Cannot read the body",
				http.StatusBadRequest)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		var dto RequestDTO
		_ = json.Unmarshal(bodyBytes, &dto)
		ipDTO := access.IPDTO{IP: dto.IP}

		ok, err := a.IsInBList(r.Context(), ipDTO)
		if err != nil && !errors.Is(err, access.ErrNotFound) {
			var errParseIP *access.ErrParseIP
			if errors.As(err, &errParseIP) {
				writeResponse(
					w,
					http.StatusText(http.StatusBadRequest),
					"Invalid Ip Address. Use IP:Mask",
					http.StatusBadRequest)
				return
			}
			if r.Context().Err() != nil {
				writeResponse(
					w,
					http.StatusText(http.StatusBadRequest),
					"Very long request. Try again",
					http.StatusBadRequest)
				return
			}
			fmt.Println(err)
			writeResponse(
				w,
				http.StatusText(http.StatusInternalServerError),
				"Ooops... Something happened wrong. Try again.",
				http.StatusInternalServerError)
			return
		}
		if ok {
			writeResponse(
				w,
				http.StatusText(http.StatusOK),
				"true",
				http.StatusOK)
			return
		}

		ok, err = a.IsInBList(r.Context(), ipDTO)
		if err != nil && !errors.Is(err, access.ErrNotFound) {
			var errParseIP *access.ErrParseIP
			if errors.As(err, &errParseIP) {
				writeResponse(
					w,
					http.StatusText(http.StatusBadRequest),
					"Invalid Ip Address. Use IP:Mask",
					http.StatusBadRequest)
				return
			}
			if r.Context().Err() != nil {
				writeResponse(
					w,
					http.StatusText(http.StatusBadRequest),
					"Very long request. Try again",
					http.StatusBadRequest)
				return
			}
			fmt.Println(err)
			writeResponse(
				w,
				http.StatusText(http.StatusInternalServerError),
				"Ooops... Something happened wrong. Try again.",
				http.StatusInternalServerError)
			return
		}
		if ok {
			writeResponse(
				w,
				http.StatusText(http.StatusBadRequest),
				"false",
				http.StatusBadRequest)
			return
		}

		next(w, r)
	}
}

func Limiter(l *limiter.MultiLimiter) func(
	w http.ResponseWriter,
	r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var dto RequestDTO
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&dto); err != nil {
			log.Println(err)
			resp := ResponseDTO{
				Status: http.StatusText(http.StatusBadRequest),
				Code:   http.StatusBadRequest,
				Msg:    "Cannot read the body",
			}
			if err = json.NewEncoder(w).Encode(resp); err != nil {
				log.Println(err)
			}
			return
		}
		bucket := l.GetBucket(
			dto.Login,
			dto.Password,
			dto.IP)
		bucket.AddTokens()
		if bucket.CheckTokensExist() {
			bucket.DeleteToken()
			resp := ResponseDTO{
				Status: http.StatusText(http.StatusOK),
				Code:   http.StatusOK,
				Msg:    "true",
			}
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				log.Println(err)
			}
			return
		}
		resp := ResponseDTO{
			Status: http.StatusText(http.StatusBadRequest),
			Code:   http.StatusBadRequest,
			Msg:    "false",
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Println(err)
		}
	}
}
