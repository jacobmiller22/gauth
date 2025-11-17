package oauth2

import (
	"context"
	"net/url"

	"github.com/jacobmiller22/gauth/pkg/oauth2"
)

type Service struct {
	RequestRepo RequestRepo
	CodeRepo    CodeRepo

	ClientReader          oauth2.ClientReader
	AuthCodeReader        oauth2.AuthCodeReader
	TokenRequestProcessor oauth2.TokenRequestProcessor
}

func (s *Service) Authorize(ctx context.Context, req *oauth2.AuthorizationReq, userId string) (*oauth2.AuthorizationRes, error) {
	if err := req.Verify(ctx, s.ClientReader); err != nil {
		return nil, err
	}

	// if req.Expired(func() time.Time { return time.Now().Add(time.Second * -5) }) {
	// 	return nil, fmt.Errorf("expired oauth2 request")
	// }

	redirectUri, err := url.ParseRequestURI(req.RedirectUri)
	if err != nil {
		return nil, err
	}

	if err := s.RequestRepo.Save(ctx, *req); err != nil {
		return nil, err
	}

	code, err := req.AuthorizationCodeSimple(userId)
	if err != nil {
		return nil, err
	}

	codeString := code.String()

	if err := s.CodeRepo.Save(ctx, codeString, req.Id.String()); err != nil {
		return nil, err
	}

	// TODO: Use NewAuthorizationResponse(req)
	return &oauth2.AuthorizationRes{
		State:        req.State,
		Code:         codeString,
		ResponseType: req.ResponseType,
		RedirectUri:  redirectUri,
	}, nil
}

func (s *Service) Token(ctx context.Context, req *oauth2.TokenRequest) (*oauth2.TokenResponse, error) {
	res, err := s.TokenRequestProcessor.BearerToken(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
