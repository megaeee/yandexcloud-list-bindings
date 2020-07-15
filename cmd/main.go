package main

import (
	"context"
	"sync"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/access"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/containerregistry/v1"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/iam/v1"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/resourcemanager/v1"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/serverless/functions/v1"
	ycsdk "github.com/yandex-cloud/go-sdk"
)

const (
	pageSize = 1000
)

func getUserID(ctx context.Context, sdk *ycsdk.SDK, user string) (userID string, err error) {
	response, err := sdk.IAM().UserAccount().Get(ctx, &iam.GetUserAccountRequest{
		UserAccountId: user,
	})
	if err == nil {
		return response.GetId(), err
	}

	response, err = sdk.IAM().YandexPassportUserAccount().GetByLogin(ctx, &iam.GetUserAccountByLoginRequest{
		Login: user,
	})

	return response.GetId(), err
}

type binding struct {
	Binding *access.AccessBinding `json:"binding"`
	Type    string                `json:"resourceType"`
	ID      string                `json:"resourceID"`
}

func listAccessBindingsPerCloud(ctx context.Context, sdk *ycsdk.SDK, wg *sync.WaitGroup, result chan []binding, cloudID, userID string) {
	defer wg.Done()
	var bindings []binding
	token := ""
	for {
		response, err := sdk.ResourceManager().Cloud().ListAccessBindings(ctx, &access.ListAccessBindingsRequest{
			ResourceId: cloudID,
			PageSize:   pageSize,
			PageToken:  token,
		})
		if err != nil {
			return
		}

		for _, k := range response.GetAccessBindings() {
			if k.GetSubject().GetId() == userID {
				bindings = append(bindings, binding{
					Binding: k,
					Type:    "cloud",
					ID:      cloudID,
				})
			}
		}

		token = response.GetNextPageToken()
		if token == "" {
			break
		}
	}

	result <- bindings
}

func listFolders(ctx context.Context, sdk *ycsdk.SDK, cloudID string) (folders []string, err error) {
	token := ""
	for {
		response, err := sdk.ResourceManager().Folder().List(ctx, &resourcemanager.ListFoldersRequest{
			CloudId:   cloudID,
			PageSize:  pageSize,
			PageToken: token,
		})
		if err != nil {
			return nil, err
		}

		for _, k := range response.GetFolders() {
			folders = append(folders, k.GetId())
		}

		token = response.GetNextPageToken()
		if token == "" {
			break
		}
	}

	return
}

func listAccessBindingsPerFolder(ctx context.Context, sdk *ycsdk.SDK, wg *sync.WaitGroup, result chan []binding, folderID, userID string) {
	defer wg.Done()
	var bindings []binding
	token := ""
	for {
		response, err := sdk.ResourceManager().Folder().ListAccessBindings(ctx, &access.ListAccessBindingsRequest{
			ResourceId: folderID,
			PageSize:   pageSize,
			PageToken:  token,
		})
		if err != nil {
			return
		}

		for _, k := range response.GetAccessBindings() {
			if k.GetSubject().GetId() == userID {
				bindings = append(bindings, binding{
					Binding: k,
					Type:    "folder",
					ID:      folderID,
				})
			}
		}

		token = response.GetNextPageToken()
		if token == "" {
			break
		}
	}

	result <- bindings
}

func listFunctions(ctx context.Context, sdk *ycsdk.SDK, folderID string) (funcs []string, err error) {
	token := ""
	for {
		response, err := sdk.Serverless().Functions().Function().List(ctx, &functions.ListFunctionsRequest{
			FolderId:  folderID,
			PageSize:  pageSize,
			PageToken: token,
		})
		if err != nil {
			return nil, err
		}

		for _, k := range response.GetFunctions() {
			funcs = append(funcs, k.GetId())
		}

		token = response.GetNextPageToken()
		if token == "" {
			break
		}
	}

	return
}

func listAccessBindingsPerFunction(ctx context.Context, sdk *ycsdk.SDK, wg *sync.WaitGroup, result chan []binding, functionID, userID string) {
	defer wg.Done()
	var bindings []binding
	token := ""
	for {
		response, err := sdk.Serverless().Functions().Function().ListAccessBindings(ctx, &access.ListAccessBindingsRequest{
			ResourceId: functionID,
			PageSize:   pageSize,
			PageToken:  token,
		})
		if err != nil {
			return
		}

		for _, k := range response.GetAccessBindings() {
			if k.GetSubject().GetId() == userID {
				bindings = append(bindings, binding{
					Binding: k,
					Type:    "function",
					ID:      functionID,
				})
			}
		}

		token = response.GetNextPageToken()
		if token == "" {
			break
		}
	}

	result <- bindings
}

func listRegistries(ctx context.Context, sdk *ycsdk.SDK, folderID string) (registries []string, err error) {
	token := ""
	for {
		response, err := sdk.ContainerRegistry().Registry().List(ctx, &containerregistry.ListRegistriesRequest{
			FolderId:  folderID,
			PageSize:  pageSize,
			PageToken: token,
		})
		if err != nil {
			return nil, err
		}

		for _, k := range response.GetRegistries() {
			registries = append(registries, k.GetId())
		}

		token = response.GetNextPageToken()
		if token == "" {
			break
		}
	}

	return
}

func listAccessBindingsPerRegistry(ctx context.Context, sdk *ycsdk.SDK, wg *sync.WaitGroup, result chan []binding, registryID, userID string) {
	defer wg.Done()
	var bindings []binding
	token := ""
	for {
		response, err := sdk.Serverless().Functions().Function().ListAccessBindings(ctx, &access.ListAccessBindingsRequest{
			ResourceId: registryID,
			PageSize:   pageSize,
			PageToken:  token,
		})
		if err != nil {
			return
		}

		for _, k := range response.GetAccessBindings() {
			if k.GetSubject().GetId() == userID {
				bindings = append(bindings, binding{
					Binding: k,
					Type:    "registry",
					ID:      registryID,
				})
			}
		}

		token = response.GetNextPageToken()
		if token == "" {
			break
		}
	}

	result <- bindings
}

func listRepositories(ctx context.Context, sdk *ycsdk.SDK, folderID string) (repositories []string, err error) {
	token := ""
	for {
		response, err := sdk.ContainerRegistry().Repository().List(ctx, &containerregistry.ListRepositoriesRequest{
			FolderId:  folderID,
			PageSize:  pageSize,
			PageToken: token,
		})
		if err != nil {
			return nil, err
		}

		for _, k := range response.GetRepositories() {
			repositories = append(repositories, k.GetId())
		}

		token = response.GetNextPageToken()
		if token == "" {
			break
		}
	}

	return
}

func listAccessBindingsPerRepository(ctx context.Context, sdk *ycsdk.SDK, wg *sync.WaitGroup, result chan []binding, repositoryID, userID string) {
	defer wg.Done()
	var bindings []binding
	token := ""
	for {
		response, err := sdk.Serverless().Functions().Function().ListAccessBindings(ctx, &access.ListAccessBindingsRequest{
			ResourceId: repositoryID,
			PageSize:   pageSize,
			PageToken:  token,
		})
		if err != nil {
			return
		}

		for _, k := range response.GetAccessBindings() {
			if k.GetSubject().GetId() == userID {
				bindings = append(bindings, binding{
					Binding: k,
					Type:    "repository",
					ID:      repositoryID,
				})
			}
		}

		token = response.GetNextPageToken()
		if token == "" {
			break
		}
	}

	result <- bindings
}

func listServiceAccounts(ctx context.Context, sdk *ycsdk.SDK, folderID string) (serviceAccounts []string, err error) {
	token := ""
	for {
		response, err := sdk.IAM().ServiceAccount().List(ctx, &iam.ListServiceAccountsRequest{
			FolderId:  folderID,
			PageSize:  pageSize,
			PageToken: token,
		})
		if err != nil {
			return nil, err
		}

		for _, k := range response.GetServiceAccounts() {
			serviceAccounts = append(serviceAccounts, k.GetId())
		}

		token = response.GetNextPageToken()
		if token == "" {
			break
		}
	}

	return
}

func listAccessBindingsPerSA(ctx context.Context, sdk *ycsdk.SDK, wg *sync.WaitGroup, result chan []binding, serviceAccountID, userID string) {
	defer wg.Done()
	var bindings []binding
	token := ""
	for {
		response, err := sdk.IAM().ServiceAccount().ListAccessBindings(ctx, &access.ListAccessBindingsRequest{
			ResourceId: serviceAccountID,
			PageSize:   pageSize,
			PageToken:  token,
		})
		if err != nil {
			return
		}

		for _, k := range response.GetAccessBindings() {
			if k.GetSubject().GetId() == userID {
				bindings = append(bindings, binding{
					Binding: k,
					Type:    "serviceAccount",
					ID:      serviceAccountID,
				})
			}
		}

		token = response.GetNextPageToken()
		if token == "" {
			break
		}
	}

	result <- bindings
}

func listAccessBindingsGo(ctx context.Context, sdk *ycsdk.SDK, resourceIDs []string, userID string, callback func(context.Context,
	*ycsdk.SDK,
	*sync.WaitGroup,
	chan []binding,
	string,
	string,
)) (result []binding) {
	ch := make(chan []binding, len(resourceIDs))
	var wg sync.WaitGroup
	for _, resourceID := range resourceIDs {
		wg.Add(1)
		go callback(ctx, sdk, &wg, ch, resourceID, userID)
	}

	wg.Wait()
	close(ch)

	result = <-ch

	return
}

type Request struct {
	CloudID string `json:"cloudId"`
	User    string `json:"user"`
}

type Response struct {
	StatusCode int         `json:"statusCode"`
	Body       interface{} `json:"body"`
}

func ListAccessBindingsPerUser(ctx context.Context, request Request) (*Response, error) {
	sdk, err := ycsdk.Build(ctx, ycsdk.Config{
		Credentials: ycsdk.InstanceServiceAccount(),
	})
	if err != nil {
		return nil, err
	}

	userID, err := getUserID(ctx, sdk, request.User)
	if err != nil {
		return nil, err
	}

	data := listAccessBindingsGo(ctx, sdk, []string{request.CloudID}, userID, listAccessBindingsPerCloud)
	folders, err := listFolders(ctx, sdk, request.CloudID)
	if err != nil {
		return nil, err
	}

	data = append(data, listAccessBindingsGo(ctx, sdk, folders, userID, listAccessBindingsPerFolder)...)

	for _, folder := range folders {
		funcs, err := listFunctions(ctx, sdk, folder)
		if err == nil {
			data = append(data, listAccessBindingsGo(ctx, sdk, funcs, userID, listAccessBindingsPerFunction)...)
		}

		regs, err := listRegistries(ctx, sdk, folder)
		if err == nil {
			data = append(data, listAccessBindingsGo(ctx, sdk, regs, userID, listAccessBindingsPerRegistry)...)
		}

		repos, err := listRepositories(ctx, sdk, folder)
		if err == nil {
			data = append(data, listAccessBindingsGo(ctx, sdk, repos, userID, listAccessBindingsPerRepository)...)
		}

		sas, err := listServiceAccounts(ctx, sdk, folder)
		if err == nil {
			data = append(data, listAccessBindingsGo(ctx, sdk, sas, userID, listAccessBindingsPerSA)...)
		}

	}

	return &Response{
		StatusCode: 200,
		Body:       data,
	}, nil
}
