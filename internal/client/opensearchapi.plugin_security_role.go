package client

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

type PluginSecurityRoleGetRequest struct {
	Name string

	Header http.Header

	ctx context.Context
}

func (r PluginSecurityRoleGetRequest) Do(ctx context.Context, transport opensearchapi.Transport) (*opensearchapi.Response, error) {
	var (
		path string
		req  *http.Request
		res  *http.Response
		err  error
	)

	path = fmt.Sprintf("/_plugins/_security/api/roles/%s", r.Name)

	if req, err = newOpenSearchRequest(ctx, http.MethodGet, path, nil); err != nil {
		return nil, err
	}

	addOpenSearchRequestHeaders(req, r.Header)

	if res, err = transport.Perform(req); err != nil {
		return nil, err
	}

	return buildOpenSearchAPIResponse(res), nil
}

type PluginSecurityRoleGet func(o ...func(*PluginSecurityRoleGetRequest)) (*opensearchapi.Response, error)

func (f PluginSecurityRoleGet) WithContext(v context.Context) func(*PluginSecurityRoleGetRequest) {
	return func(r *PluginSecurityRoleGetRequest) {
		r.ctx = v
	}
}

func (f PluginSecurityRoleGet) WithName(v string) func(*PluginSecurityRoleGetRequest) {
	return func(r *PluginSecurityRoleGetRequest) {
		r.Name = v
	}
}

func (f PluginSecurityRoleGet) WithHeader(n map[string]string) func(*PluginSecurityRoleGetRequest) {
	return func(r *PluginSecurityRoleGetRequest) {
		if r.Header == nil {
			r.Header = make(http.Header, 0)
		}
		for k, v := range n {
			r.Header.Add(k, v)
		}
	}
}

type PluginSecurityRolesGetRequest struct {
	Header http.Header

	ctx context.Context
}

func (r PluginSecurityRolesGetRequest) Do(ctx context.Context, transport opensearchapi.Transport) (*opensearchapi.Response, error) {
	var (
		path string
		req  *http.Request
		res  *http.Response
		err  error
	)

	path = "/_plugins/_security/api/roles/"

	if req, err = newOpenSearchRequest(ctx, http.MethodGet, path, nil); err != nil {
		return nil, err
	}

	addOpenSearchRequestHeaders(req, r.Header)

	if res, err = transport.Perform(req); err != nil {
		return nil, err
	}

	return buildOpenSearchAPIResponse(res), nil
}

type PluginSecurityRoles func(o ...func(*PluginSecurityRolesGetRequest)) (*opensearchapi.Response, error)

func (f PluginSecurityRoles) WithContext(v context.Context) func(*PluginSecurityRolesGetRequest) {
	return func(r *PluginSecurityRolesGetRequest) {
		r.ctx = v
	}
}

func (f PluginSecurityRoles) WithHeader(n map[string]string) func(*PluginSecurityRolesGetRequest) {
	return func(r *PluginSecurityRolesGetRequest) {
		if r.Header == nil {
			r.Header = make(http.Header, 0)
		}
		for k, v := range n {
			r.Header.Add(k, v)
		}
	}
}

type PluginSecurityRoleDeleteRequest struct {
	Name string

	Header http.Header

	ctx context.Context
}

func (r *PluginSecurityRoleDeleteRequest) Do(ctx context.Context, transport opensearchapi.Transport) (*opensearchapi.Response, error) {
	var (
		path string
		req  *http.Request
		res  *http.Response
		err  error
	)

	path = fmt.Sprintf("/_plugins/_security/roles/%s", r.Name)

	if req, err = newOpenSearchRequest(ctx, http.MethodDelete, path, nil); err != nil {
		return nil, err
	}

	addOpenSearchRequestHeaders(req, r.Header)

	if res, err = transport.Perform(req); err != nil {
		return nil, err
	}

	return buildOpenSearchAPIResponse(res), nil
}

type PluginSecurityRoleDelete func(o ...func(*PluginSecurityRoleDeleteRequest)) (*opensearchapi.Response, error)

func (f PluginSecurityRoleDelete) WithContext(v context.Context) func(*PluginSecurityRoleDeleteRequest) {
	return func(r *PluginSecurityRoleDeleteRequest) {
		r.ctx = v
	}
}

func (f PluginSecurityRoleDelete) WithName(v string) func(*PluginSecurityRoleDeleteRequest) {
	return func(r *PluginSecurityRoleDeleteRequest) {
		r.Name = v
	}
}

func (f PluginSecurityRoleDelete) WithHeader(n map[string]string) func(*PluginSecurityRoleDeleteRequest) {
	return func(r *PluginSecurityRoleDeleteRequest) {
		if r.Header == nil {
			r.Header = make(http.Header, 0)
		}
		for k, v := range n {
			r.Header.Add(k, v)
		}
	}
}

type PluginSecurityRoleUpsertRequest struct {
	Name string

	Body io.Reader

	Header http.Header

	ctx context.Context
}

func (r *PluginSecurityRoleUpsertRequest) Do(ctx context.Context, transport opensearchapi.Transport) (*opensearchapi.Response, error) {
	var (
		path string
		req  *http.Request
		res  *http.Response
		err  error
	)

	path = fmt.Sprintf("/_plugins/_security/roles/%s", r.Name)

	if req, err = newOpenSearchRequest(ctx, http.MethodPut, path, r.Body); err != nil {
		return nil, err
	}

	addOpenSearchRequestHeaders(req, r.Header)

	if res, err = transport.Perform(req); err != nil {
		return nil, err
	}

	return buildOpenSearchAPIResponse(res), nil
}

type PluginSecurityRoleUpsert func(o ...func(request *PluginSecurityRoleUpsertRequest)) (*opensearchapi.Response, error)

func (f PluginSecurityRoleUpsert) WithContext(v context.Context) func(*PluginSecurityRoleUpsertRequest) {
	return func(r *PluginSecurityRoleUpsertRequest) {
		r.ctx = v
	}
}

func (f PluginSecurityRoleUpsert) WithName(v string) func(request *PluginSecurityRoleUpsertRequest) {
	return func(r *PluginSecurityRoleUpsertRequest) {
		r.Name = v
	}
}

func (f PluginSecurityRoleUpsert) WithBody(v io.Reader) func(*PluginSecurityRoleUpsertRequest) {
	return func(r *PluginSecurityRoleUpsertRequest) {
		r.Body = v
	}
}

func (f PluginSecurityRoleUpsert) WithHeader(n map[string]string) func(*PluginSecurityRoleUpsertRequest) {
	return func(r *PluginSecurityRoleUpsertRequest) {
		if r.Header == nil {
			r.Header = make(http.Header, 0)
		}
		for k, v := range n {
			r.Header.Add(k, v)
		}
	}
}
