package kate

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"hash"
	"strconv"
	"time"

	"github.com/k81/kate/context"
	"github.com/k81/kate/log"
	"github.com/k81/kate/utils"
)

var (
	RequiredFields = []string{"from", "nonce", "timestamp"}
)

func VerifySign(keyStore KeyStore) Middleware {
	return func(h ContextHandler) ContextHandler {
		f := func(ctx context.Context, w ResponseWriter, r *Request) {
			var (
				timestamp  int64
				appKey     string
				buf        bytes.Buffer
				mac        hash.Hash
				clientSign string
				serverMAC  []byte
				clientMAC  []byte
				err        error
			)

			for _, k := range RequiredFields {
				if r.Form.Get(k) == "" {
					panic(ErrUnauthorized)
				}
			}

			clientSign = r.Header.Get("X-Signature")

			if clientSign == "" {
				log.Warning(ctx, "missing client sign")
				panic(ErrUnauthorized)
			}

			if clientMAC, err = hex.DecodeString(clientSign); err != nil {
				log.Warning(ctx, "decode client sign", "client_sign", clientSign, "error", err)
				panic(ErrUnauthorized)
			}

			if timestamp, err = strconv.ParseInt(r.Form.Get("timestamp"), 10, 64); err != nil {
				log.Warning(ctx, "parse timestamp", "timestamp", r.Form.Get("timestamp"), "error", err)
				panic(ErrBadRequest)
			}

			if utils.Abs(time.Now().Unix()-timestamp) > 60 {
				panic(ErrRequestExpired)
			}

			appKey = keyStore.GetKey(r.Form.Get("from"))
			if appKey == "" {
				log.Error(ctx, "appkey not found", "from", r.Form.Get("from"))
				panic(ErrUnauthorized)
			}

			buf.WriteString(r.RequestURI)

			if len(r.RawBody) > 0 {
				buf.WriteString(string(r.RawBody))
			}

			mac = hmac.New(sha1.New, []byte(appKey))
			mac.Write(buf.Bytes())
			serverMAC = mac.Sum(nil)

			if !hmac.Equal(clientMAC, serverMAC) {
				log.Error(ctx, "sign not match", "client_sign", clientSign, "server_sign", hex.EncodeToString(serverMAC), "sign_input", buf.String())
				panic(ErrUnauthorized)
			}

			h.ServeHTTP(ctx, w, r)
		}
		return ContextHandlerFunc(f)
	}
}
