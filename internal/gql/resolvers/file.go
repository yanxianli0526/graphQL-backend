package resolvers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	orm "graphql-go-template/internal/database"
	"net/http"
	"os"
	"strconv"
	"time"

	gqlmodels "graphql-go-template/internal/gql/models"

	"graphql-go-template/internal/models"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PrintRequestData struct {
	Ct   string `json:"ct"`
	Name string `json:"name"`
}

type PrintRequest struct {
	Data []PrintRequestData `json:"data"`
}

type PrintResponseData struct {
	SignedUrl  string `json:"signedUrl"`
	PublicLink string `json:"publicLink"`
	FullName   string `json:"fullName"`
}

type PrintResponse struct {
	Data []PrintResponseData `json:"data"`
}

// Mutations
func (r *mutationResolver) CreateFile(ctx context.Context, fileName string) (*gqlmodels.UploadFileResponse, error) {
	reqContext := graphql.GetOperationContext(ctx)
	apiStartTime := time.Now()
	userIdStr := fmt.Sprintf("%v", ctx.Value(UserIdCtxKey))
	organizationIdStr := fmt.Sprintf("%v", ctx.Value(OrganizationId))
	organizationId, err := uuid.Parse(organizationIdStr)
	if err != nil {
		r.Logger.Warn("CreateFile uuid.Parse(organizationIdStr)", zap.Error(err), zap.String("originalUrl", "createFile"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return nil, err
	}

	_, err = os.Create(fileName)
	if err != nil {
		r.Logger.Error("CreateFile os.Create", zap.Error(err))
		return nil, err
	}

	client := &http.Client{}
	endPoint := r.PhotoServiceProtocol + "://" + r.PhotoServiceHost + ":" + r.PhotoServicePort + "/signedurl"

	var printRequestDataArray []PrintRequestData

	printRequestData := PrintRequestData{
		Ct:   "image/jpeg",
		Name: "inventory-tool/" + organizationIdStr + "/" + fileName,
	}
	printRequestDataArray = append(printRequestDataArray, printRequestData)
	observation := PrintRequest{
		Data: printRequestDataArray,
	}

	bodyJSON, err := json.Marshal(observation)
	if err != nil {
		r.Logger.Error("CreateFile json.Marshal(observation)", zap.Error(err), zap.String("originalUrl", "createFile"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return nil, nil
	}
	res, err := client.Post(endPoint, "application/json", bytes.NewBuffer(bodyJSON))
	if err != nil {
		r.Logger.Error("CreateFile client.Post", zap.Error(err), zap.String("originalUrl", "createFile"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return nil, fmt.Errorf("client.Post failed: %v", err)
	}

	if code := res.StatusCode; code < 200 || code > 299 {
		r.Logger.Error("CreateFile res.StatusCode; code < 200 || code > 299", zap.Error(err), zap.String("originalUrl", "createFile"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return nil, fmt.Errorf("status code failed: %v", err)
	}

	printResponse := &PrintResponse{}
	if err = json.NewDecoder(res.Body).Decode(printResponse); err != nil {
		r.Logger.Error("CreateFile json.NewDecoder(res.Body).Decode(printResponse)", zap.Error(err), zap.String("originalUrl", "createFile"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return nil, fmt.Errorf("decode json failed: %w", err)
	}

	file := models.File{
		ID:             uuid.New(),
		FileName:       fileName,
		Url:            printResponse.Data[0].PublicLink,
		OrganizationId: organizationId,
	}

	err = orm.CreateFile(r.ORM.DB, &file)
	if err != nil {
		r.Logger.Error("CreateFile orm.CreateFile", zap.Error(err), zap.String("originalUrl", "createFile"),
			zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))
		return nil, err
	}

	uploadFileResponse := gqlmodels.UploadFileResponse{
		SignedURL:  printResponse.Data[0].SignedUrl,
		PublicLink: printResponse.Data[0].PublicLink,
		FullName:   printResponse.Data[0].FullName,
	}

	apiRunTime, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(time.Now().UnixMicro()-apiStartTime.UnixMicro())/1000000), 64)
	r.Logger.Info("createFile run time", zap.Float64("duration_sec", apiRunTime), zap.String("fieldName", "createFile"),
		zap.String("message", "REQ END"), zap.String("userId", userIdStr), zap.Int64("timestamp", time.Now().Unix()), zap.Any("requestBody", reqContext.Variables))

	return &uploadFileResponse, nil
}

// fileField resolvers
type fileResolver struct{ *Resolver }

func (r *fileResolver) ID(ctx context.Context, obj *models.File) (string, error) {
	return obj.ID.String(), nil
}
