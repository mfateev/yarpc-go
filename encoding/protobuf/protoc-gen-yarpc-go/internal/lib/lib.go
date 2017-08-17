// Copyright (c) 2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Package lib contains the library code for protoc-gen-yarpc-go.
//
// It is split into a separate package so it can be called by the testing package.
package lib

import (
	"strings"
	"text/template"

	"go.uber.org/yarpc/internal/protoplugin"
)

const tmpl = `{{$packagePath := .GoPackage.Path}}
// Code generated by protoc-gen-yarpc-go
// source: {{.GetName}}
// DO NOT EDIT!

package {{.GoPackage.Name}}
{{if .Services}}
import (
	{{range $i := .Imports}}{{if $i.Standard}}{{$i | printf "%s\n"}}{{end}}{{end}}

	{{range $i := .Imports}}{{if not $i.Standard}}{{$i | printf "%s\n"}}{{end}}{{end}}
){{end}}

{{range $service := .Services}}
// {{$service.GetName}}YARPCClient is the YARPC client-side interface for the {{$service.GetName}} service.
type {{$service.GetName}}YARPCClient interface {
	{{range $method := unaryMethods $service}}{{$method.GetName}}(context.Context, *{{$method.RequestType.GoType $packagePath}}, ...yarpc.CallOption) (*{{$method.ResponseType.GoType $packagePath}}, error)
	{{end}}{{range $method := onewayMethods $service}}{{$method.GetName}}(context.Context, *{{$method.RequestType.GoType $packagePath}}, ...yarpc.CallOption) (yarpc.Ack, error)
	{{end}}{{range $method := clientStreamingMethods $service}}{{$method.GetName}}(context.Context, ...yarpc.CallOption) ({{$service.GetName}}Service{{$method.GetName}}YARPCClient, error)
	{{end}}{{range $method := serverStreamingMethods $service}}{{$method.GetName}}(context.Context, *{{$method.RequestType.GoType $packagePath}}, ...yarpc.CallOption) ({{$service.GetName}}Service{{$method.GetName}}YARPCClient, error)
	{{end}}{{range $method := clientServerStreamingMethods $service}}{{$method.GetName}}(context.Context, ...yarpc.CallOption) ({{$service.GetName}}Service{{$method.GetName}}YARPCClient, error)
	{{end}}
}

{{range $method := clientStreamingMethods $service}}
// {{$service.GetName}}Service{{$method.GetName}}YARPCClient sends {{$method.RequestType.GoType $packagePath}}s and receives the single {{$method.ResponseType.GoType $packagePath}} when sending is done.
type {{$service.GetName}}Service{{$method.GetName}}YARPCClient interface {
	Context() context.Context
	Request() *transport.Request
	Send(*{{$method.RequestType.GoType $packagePath}}) error
	CloseAndRecv() (*{{$method.ResponseType.GoType $packagePath}}, error)
}
{{end}}

{{range $method := serverStreamingMethods $service}}
// {{$service.GetName}}Service{{$method.GetName}}YARPCClient receives {{$method.ResponseType.GoType $packagePath}}s, returning io.EOF when the stream is complete.
type {{$service.GetName}}Service{{$method.GetName}}YARPCClient interface {
	Context() context.Context
	Request() *transport.Request
	Recv() (*{{$method.ResponseType.GoType $packagePath}}, error)
}
{{end}}

{{range $method := clientServerStreamingMethods $service}}
// {{$service.GetName}}Service{{$method.GetName}}YARPCClient sends {{$method.RequestType.GoType $packagePath}}s and receives {{$method.ResponseType.GoType $packagePath}}s, returning io.EOF when the stream is complete.
type {{$service.GetName}}Service{{$method.GetName}}YARPCClient interface {
	Context() context.Context
	Request() *transport.Request
	Send(*{{$method.RequestType.GoType $packagePath}}) error
	Recv() (*{{$method.ResponseType.GoType $packagePath}}, error)
	CloseSend() error
}
{{end}}

// New{{$service.GetName}}YARPCClient builds a new YARPC client for the {{$service.GetName}} service.
func New{{$service.GetName}}YARPCClient(clientConfig transport.ClientConfig, options ...protobuf.ClientOption) {{$service.GetName}}YARPCClient {
	return &_{{$service.GetName}}YARPCCaller{protobuf.NewClient(
		protobuf.ClientParams{
			ServiceName: "{{trimPrefixPeriod $service.FQSN}}",
			ClientConfig: clientConfig,
			Options: options,
		},
	)}
}

// {{$service.GetName}}YARPCServer is the YARPC server-side interface for the {{$service.GetName}} service.
type {{$service.GetName}}YARPCServer interface {
	{{range $method := unaryMethods $service}}{{$method.GetName}}(context.Context, *{{$method.RequestType.GoType $packagePath}}) (*{{$method.ResponseType.GoType $packagePath}}, error)
	{{end}}{{range $method := onewayMethods $service}}{{$method.GetName}}(context.Context, *{{$method.RequestType.GoType $packagePath}}) error
	{{end}}{{range $method := clientStreamingMethods $service}}{{$method.GetName}}({{$service.GetName}}Service{{$method.GetName}}YARPCServer) (*{{$method.ResponseType.GoType $packagePath}}, error)
	{{end}}{{range $method := serverStreamingMethods $service}}{{$method.GetName}}(*{{$method.RequestType.GoType $packagePath}}, {{$service.GetName}}Service{{$method.GetName}}YARPCServer) error
	{{end}}{{range $method := clientServerStreamingMethods $service}}{{$method.GetName}}({{$service.GetName}}Service{{$method.GetName}}YARPCServer) error
	{{end}}
}

{{range $method := clientStreamingMethods $service}}
// {{$service.GetName}}Service{{$method.GetName}}YARPCServer receives {{$method.RequestType.GoType $packagePath}}s.
type {{$service.GetName}}Service{{$method.GetName}}YARPCServer interface {
	Context() context.Context
	Request() *transport.Request
	Recv() (*{{$method.RequestType.GoType $packagePath}}, error)
}
{{end}}

{{range $method := serverStreamingMethods $service}}
// {{$service.GetName}}Service{{$method.GetName}}YARPCServer sends {{$method.ResponseType.GoType $packagePath}}s.
type {{$service.GetName}}Service{{$method.GetName}}YARPCServer interface {
	Context() context.Context
	Request() *transport.Request
	Send(*{{$method.ResponseType.GoType $packagePath}}) error
}
{{end}}

{{range $method := clientServerStreamingMethods $service}}
// {{$service.GetName}}Service{{$method.GetName}}YARPCServer receives {{$method.RequestType.GoType $packagePath}}s and sends {{$method.ResponseType.GoType $packagePath}}.
type {{$service.GetName}}Service{{$method.GetName}}YARPCServer interface {
	Context() context.Context
	Request() *transport.Request
	Recv() (*{{$method.RequestType.GoType $packagePath}}, error)
	Send(*{{$method.ResponseType.GoType $packagePath}}) error
	CloseSend() error
}
{{end}}

// Build{{$service.GetName}}YARPCProcedures prepares an implementation of the {{$service.GetName}} service for YARPC registration.
func Build{{$service.GetName}}YARPCProcedures(server {{$service.GetName}}YARPCServer) []transport.Procedure {
	handler := &_{{$service.GetName}}YARPCHandler{server}
	return protobuf.BuildProcedures(
		protobuf.BuildProceduresParams{
			ServiceName: "{{trimPrefixPeriod $service.FQSN}}",
			UnaryHandlerParams: []protobuf.BuildProceduresUnaryHandlerParams{
			{{range $method := unaryMethods $service}}{
					MethodName: "{{$method.GetName}}",
					Handler: protobuf.NewUnaryHandler(
						protobuf.UnaryHandlerParams{
							Handle: handler.{{$method.GetName}},
							NewRequest: new{{$service.GetName}}Service{{$method.GetName}}YARPCRequest,
						},
					),
				},
			{{end}}
			},
			OnewayHandlerParams: []protobuf.BuildProceduresOnewayHandlerParams{
			{{range $method := onewayMethods $service}}{
					MethodName: "{{$method.GetName}}",
					Handler: protobuf.NewOnewayHandler(
						protobuf.OnewayHandlerParams{
							Handle: handler.{{$method.GetName}},
							NewRequest: new{{$service.GetName}}Service{{$method.GetName}}YARPCRequest,
						},
					),
				},
			{{end}}
			},
		},
	)
}

type _{{$service.GetName}}YARPCCaller struct {
	client protobuf.Client
}

{{range $method := unaryMethods $service}}
func (c *_{{$service.GetName}}YARPCCaller) {{$method.GetName}}(ctx context.Context, request *{{$method.RequestType.GoType $packagePath}}, options ...yarpc.CallOption) (*{{$method.ResponseType.GoType $packagePath}}, error) {
	responseMessage, err := c.client.Call(ctx, "{{$method.GetName}}", request, new{{$service.GetName}}Service{{$method.GetName}}YARPCResponse, options...)
	if responseMessage == nil {
		return nil, err
	}
	response, ok := responseMessage.(*{{$method.ResponseType.GoType $packagePath}})
	if !ok {
		return nil, protobuf.CastError(empty{{$service.GetName}}Service{{$method.GetName}}YARPCResponse, responseMessage)
	}
	return response, err
}
{{end}}
{{range $method := onewayMethods $service}}
func (c *_{{$service.GetName}}YARPCCaller) {{$method.GetName}}(ctx context.Context, request *{{$method.RequestType.GoType $packagePath}}, options ...yarpc.CallOption) (yarpc.Ack, error) {
	return c.client.CallOneway(ctx, "{{$method.GetName}}", request, options...)
}
{{end}}
{{range $method := clientStreamingMethods $service}}
func (c *_{{$service.GetName}}YARPCCaller) {{$method.GetName}}(ctx context.Context, options ...yarpc.CallOption) ({{$service.GetName}}Service{{$method.GetName}}YARPCClient, error) {
	// TODO
	return &_{{$service.GetName}}Service{{$method.GetName}}YARPCClient{}, nil
}
{{end}}
{{range $method := serverStreamingMethods $service}}
func (c *_{{$service.GetName}}YARPCCaller) {{$method.GetName}}(ctx context.Context, request *{{$method.RequestType.GoType $packagePath}}, options ...yarpc.CallOption) ({{$service.GetName}}Service{{$method.GetName}}YARPCClient, error) {
	// TODO
	return &_{{$service.GetName}}Service{{$method.GetName}}YARPCClient{}, nil
}
{{end}}
{{range $method := clientServerStreamingMethods $service}}
func (c *_{{$service.GetName}}YARPCCaller) {{$method.GetName}}(ctx context.Context, options ...yarpc.CallOption) ({{$service.GetName}}Service{{$method.GetName}}YARPCClient, error) {
	// TODO
	return &_{{$service.GetName}}Service{{$method.GetName}}YARPCClient{}, nil
}
{{end}}

type _{{$service.GetName}}YARPCHandler struct {
	server {{$service.GetName}}YARPCServer
}

{{range $method := unaryMethods $service}}
func (h *_{{$service.GetName}}YARPCHandler) {{$method.GetName}}(ctx context.Context, requestMessage proto.Message) (proto.Message, error) {
	var request *{{$method.RequestType.GoType $packagePath}}
	var ok bool
	if requestMessage != nil {
		request, ok = requestMessage.(*{{$method.RequestType.GoType $packagePath}})
		if !ok {
			return nil, protobuf.CastError(empty{{$service.GetName}}Service{{$method.GetName}}YARPCRequest, requestMessage)
		}
	}
	response, err := h.server.{{$method.GetName}}(ctx, request)
	if response == nil {
		return nil, err
	}
	return response, err
}
{{end}}
{{range $method := onewayMethods $service}}
func (h *_{{$service.GetName}}YARPCHandler) {{$method.GetName}}(ctx context.Context, requestMessage proto.Message) error {
	var request *{{$method.RequestType.GoType $packagePath}}
	var ok bool
	if requestMessage != nil {
		request, ok = requestMessage.(*{{$method.RequestType.GoType $packagePath}})
		if !ok {
			return protobuf.CastError(empty{{$service.GetName}}Service{{$method.GetName}}YARPCRequest, requestMessage)
		}
	}
	return h.server.{{$method.GetName}}(ctx, request)
}
{{end}}
{{range $method := clientStreamingMethods $service}}
func (h *_{{$service.GetName}}YARPCHandler) {{$method.GetName}}(serverStream transport.ServerStream) (*{{$method.ResponseType.GoType $packagePath}}, error) {
	return h.server.{{$method.GetName}}(&_{{$service.GetName}}Service{{$method.GetName}}YARPCServer{serverStream: serverStream})
}
{{end}}
{{range $method := serverStreamingMethods $service}}
func (h *_{{$service.GetName}}YARPCHandler) {{$method.GetName}}(requestMessage proto.Message , serverStream transport.ServerStream) error {
	var request *{{$method.RequestType.GoType $packagePath}}
	var ok bool
	if requestMessage != nil {
		request, ok = requestMessage.(*{{$method.RequestType.GoType $packagePath}})
		if !ok {
			return protobuf.CastError(empty{{$service.GetName}}Service{{$method.GetName}}YARPCRequest, requestMessage)
		}
	}
	return h.server.{{$method.GetName}}(request, &_{{$service.GetName}}Service{{$method.GetName}}YARPCServer{serverStream: serverStream})
}
{{end}}
{{range $method := clientServerStreamingMethods $service}}
func (h *_{{$service.GetName}}YARPCHandler) {{$method.GetName}}(serverStream transport.ServerStream) error {
	return h.server.{{$method.GetName}}(&_{{$service.GetName}}Service{{$method.GetName}}YARPCServer{serverStream: serverStream})
}
{{end}}

{{range $method := clientStreamingMethods $service}}
type _{{$service.GetName}}Service{{$method.GetName}}YARPCClient struct {
	// TODO
}

func (c *_{{$service.GetName}}Service{{$method.GetName}}YARPCClient) Context() context.Context {
	// TODO
	return nil
}

func (c *_{{$service.GetName}}Service{{$method.GetName}}YARPCClient) Request() *transport.Request {
	// TODO
	return nil
}

func (c *_{{$service.GetName}}Service{{$method.GetName}}YARPCClient) Send(request *{{$method.RequestType.GoType $packagePath}}) error {
	// TODO
	return nil
}

func (c *_{{$service.GetName}}Service{{$method.GetName}}YARPCClient) CloseAndRecv() (*{{$method.ResponseType.GoType $packagePath}}, error) {
	// TODO
	return nil, nil
}
{{end}}

{{range $method := serverStreamingMethods $service}}
type _{{$service.GetName}}Service{{$method.GetName}}YARPCClient struct {
	//TODO
}

func (c *_{{$service.GetName}}Service{{$method.GetName}}YARPCClient) Context() context.Context {
	//TODO
	return nil
}

func (c *_{{$service.GetName}}Service{{$method.GetName}}YARPCClient) Request() *transport.Request {
	//TODO
	return nil
}

func (c *_{{$service.GetName}}Service{{$method.GetName}}YARPCClient) Recv() (*{{$method.ResponseType.GoType $packagePath}}, error) {
	//TODO
	return nil, nil
}
{{end}}

{{range $method := clientServerStreamingMethods $service}}
type _{{$service.GetName}}Service{{$method.GetName}}YARPCClient struct {
	//TODO
}

func (c *_{{$service.GetName}}Service{{$method.GetName}}YARPCClient) Context() context.Context {
	//TODO
	return nil
}

func (c *_{{$service.GetName}}Service{{$method.GetName}}YARPCClient) Request() *transport.Request {
	//TODO
	return nil
}

func (c *_{{$service.GetName}}Service{{$method.GetName}}YARPCClient) Send(request *{{$method.RequestType.GoType $packagePath}}) error {
	//TODO
	return nil
}

func (c *_{{$service.GetName}}Service{{$method.GetName}}YARPCClient) Recv() (*{{$method.ResponseType.GoType $packagePath}}, error) {
	//TODO
	return nil, nil
}

func (c *_{{$service.GetName}}Service{{$method.GetName}}YARPCClient) CloseSend() error {
	return nil
}
{{end}}

{{range $method := clientStreamingMethods $service}}
type _{{$service.GetName}}Service{{$method.GetName}}YARPCServer struct {
	serverStream transport.ServerStream
}

func (s *_{{$service.GetName}}Service{{$method.GetName}}YARPCServer) Context() context.Context {
	// TODO
	return nil
}

func (s *_{{$service.GetName}}Service{{$method.GetName}}YARPCServer) Request() *transport.Request {
	// TODO
	return nil
}

func (s *_{{$service.GetName}}Service{{$method.GetName}}YARPCServer) Recv() (*{{$method.RequestType.GoType $packagePath}}, error) {
	// TODO
	return nil, nil
}
{{end}}

{{range $method := serverStreamingMethods $service}}
type _{{$service.GetName}}Service{{$method.GetName}}YARPCServer struct {
	serverStream transport.ServerStream
}

func (s *_{{$service.GetName}}Service{{$method.GetName}}YARPCServer) Context() context.Context {
	// TODO
	return nil
}

func (s *_{{$service.GetName}}Service{{$method.GetName}}YARPCServer) Request() *transport.Request {
	// TODO
	return nil
}

func (s *_{{$service.GetName}}Service{{$method.GetName}}YARPCServer) Send(response *{{$method.ResponseType.GoType $packagePath}}) error {
	// TODO
	return nil
}
{{end}}

{{range $method := clientServerStreamingMethods $service}}
type _{{$service.GetName}}Service{{$method.GetName}}YARPCServer struct {
	serverStream transport.ServerStream
}

func (s *_{{$service.GetName}}Service{{$method.GetName}}YARPCServer) Context() context.Context {
	// TODO
	return nil
}

func (s *_{{$service.GetName}}Service{{$method.GetName}}YARPCServer) Request() *transport.Request {
	// TODO
	return nil
}

func (s *_{{$service.GetName}}Service{{$method.GetName}}YARPCServer) Recv() (*{{$method.RequestType.GoType $packagePath}}, error) {
	// TODO
	return nil, nil
}

func (s *_{{$service.GetName}}Service{{$method.GetName}}YARPCServer) Send(response *{{$method.ResponseType.GoType $packagePath}}) error {
	// TODO
	return nil
}

func (s *_{{$service.GetName}}Service{{$method.GetName}}YARPCServer) CloseSend() error {
	// TODO
	return nil
}
{{end}}

{{range $method := $service.Methods}}
func new{{$service.GetName}}Service{{$method.GetName}}YARPCRequest() proto.Message {
	return &{{$method.RequestType.GoType $packagePath}}{}
}

func new{{$service.GetName}}Service{{$method.GetName}}YARPCResponse() proto.Message {
	return &{{$method.ResponseType.GoType $packagePath}}{}
}
{{end}}
var (
{{range $method := $service.Methods}}
	empty{{$service.GetName}}Service{{$method.GetName}}YARPCRequest = &{{$method.RequestType.GoType $packagePath}}{}
	empty{{$service.GetName}}Service{{$method.GetName}}YARPCResponse = &{{$method.ResponseType.GoType $packagePath}}{}{{end}}
)
{{end}}
{{if .Services}}func init() { {{range $service := .Services}}
	yarpc.RegisterClientBuilder(
		func(clientConfig transport.ClientConfig, structField reflect.StructField) {{$service.GetName}}YARPCClient {
			return New{{$service.GetName}}YARPCClient(clientConfig, protobuf.ClientBuilderOptions(clientConfig, structField)...)
		},
	){{end}}
}{{end}}
`

// Runner is the Runner used for protoc-gen-yarpc-go.
var Runner = protoplugin.NewRunner(
	template.Must(template.New("tmpl").Funcs(
		template.FuncMap{
			"unaryMethods":                 unaryMethods,
			"onewayMethods":                onewayMethods,
			"clientStreamingMethods":       clientStreamingMethods,
			"serverStreamingMethods":       serverStreamingMethods,
			"clientServerStreamingMethods": clientServerStreamingMethods,
			"trimPrefixPeriod":             trimPrefixPeriod,
		}).Parse(tmpl)),
	checkTemplateInfo,
	[]string{
		"context",
		"reflect",
		"github.com/gogo/protobuf/proto",
		"go.uber.org/yarpc",
		"go.uber.org/yarpc/api/transport",
		"go.uber.org/yarpc/encoding/protobuf",
	},
	"pb.yarpc.go",
)

func checkTemplateInfo(templateInfo *protoplugin.TemplateInfo) error {
	return nil
}

func unaryMethods(service *protoplugin.Service) ([]*protoplugin.Method, error) {
	methods := make([]*protoplugin.Method, 0, len(service.Methods))
	for _, method := range service.Methods {
		if !method.GetClientStreaming() && !method.GetServerStreaming() && method.ResponseType.FQMN() != ".uber.yarpc.Oneway" {
			methods = append(methods, method)
		}
	}
	return methods, nil
}

func onewayMethods(service *protoplugin.Service) ([]*protoplugin.Method, error) {
	methods := make([]*protoplugin.Method, 0, len(service.Methods))
	for _, method := range service.Methods {
		if !method.GetClientStreaming() && !method.GetServerStreaming() && method.ResponseType.FQMN() == ".uber.yarpc.Oneway" {
			methods = append(methods, method)
		}
	}
	return methods, nil
}

func clientStreamingMethods(service *protoplugin.Service) ([]*protoplugin.Method, error) {
	methods := make([]*protoplugin.Method, 0, len(service.Methods))
	for _, method := range service.Methods {
		if method.GetClientStreaming() && !method.GetServerStreaming() {
			methods = append(methods, method)
		}
	}
	return methods, nil
}

func serverStreamingMethods(service *protoplugin.Service) ([]*protoplugin.Method, error) {
	methods := make([]*protoplugin.Method, 0, len(service.Methods))
	for _, method := range service.Methods {
		if !method.GetClientStreaming() && method.GetServerStreaming() {
			methods = append(methods, method)
		}
	}
	return methods, nil
}

func clientServerStreamingMethods(service *protoplugin.Service) ([]*protoplugin.Method, error) {
	methods := make([]*protoplugin.Method, 0, len(service.Methods))
	for _, method := range service.Methods {
		if method.GetClientStreaming() && method.GetServerStreaming() {
			methods = append(methods, method)
		}
	}
	return methods, nil
}

func trimPrefixPeriod(s string) string {
	return strings.TrimPrefix(s, ".")
}
