package server

import (
	"net/http"
	"strings"

	"code.google.com/p/go-uuid/uuid"

	"github.com/dgrijalva/jwt-go"
	"github.com/emicklei/go-restful"
)

type AuthResource struct {
	signingKey []byte
	sessions   map[string]jwt.Token
}

type TokenRequest struct {
	GrantType string `json:"grant_type"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	ClientId  string `json:"client_id"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	IdToken     string `json:"id_token"`
}

func RegisterAuth(c *restful.Container, signingKey []byte) *AuthResource {
	h := &AuthResource{sessions: make(map[string]jwt.Token), signingKey: signingKey}
	c.Filter(h.AuthorizationFilter)

	ws := new(restful.WebService)

	ws.Path("/auth").
		Doc("Manages authorization").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_XML, restful.MIME_JSON)

	ws.Route(ws.POST("/token").To(h.createSession).
		Doc("create a new session").
		Operation("createSession").
		Reads(TokenRequest{}).
		Writes(""))

	c.Add(ws)

	return h
}

func (h *AuthResource) AuthorizationFilter(req *restful.Request, res *restful.Response, chain *restful.FilterChain) {
	// auth/token is exempt
	if req.SelectedRoutePath() == "/auth/token" {
		chain.ProcessFilter(req, res)
		return
	}

	token := req.Request.Header.Get("Authorization")

	if strings.HasPrefix(token, "Bearer ") {
		if t, ok := h.sessions[token[7:]]; ok {
			req.SetAttribute("token", t)
			chain.ProcessFilter(req, res)
			return
		}
	}

	res.AddHeader("WWW-Authenticate", `OAuth realm="http://localhost:3001/"`)
	res.WriteErrorString(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
}

func (h *AuthResource) createSession(req *restful.Request, res *restful.Response) {
	tr := new(TokenRequest)
	req.ReadEntity(tr)

	if tr.Username == "username" && tr.Password == "password" {
		accessToken := uuid.New()
		token := jwt.New(jwt.SigningMethodHS256)
		token.Claims["sub"] = tr.Username
		token.Claims["name"] = tr.Username
		token.Claims["roles"] = []string{"admin", "shmurda"}
		tokenString, err := token.SignedString(h.signingKey)
		if err != nil {
			res.WriteErrorString(500, err.Error())
		}

		h.sessions[tokenString] = *token

		response := TokenResponse{
			AccessToken: accessToken,
			TokenType:   "bearer",
			ExpiresIn:   3600,
			IdToken:     tokenString,
		}

		res.WriteHeader(http.StatusFound)
		res.WriteEntity(response)
	} else {
		res.WriteErrorString(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}
}