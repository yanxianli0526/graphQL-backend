package handlers

import (
	config "graphql-go-template/envconfig"
	"graphql-go-template/internal/auth"
	orm "graphql-go-template/internal/database"
	gql "graphql-go-template/internal/gql/generated"
	"graphql-go-template/internal/gql/resolvers"
	objectStorage "graphql-go-template/pkg/gcp"

	"github.com/99designs/gqlgen/handler"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GraphqlHandler defines the GQLGen GraphQL server handler
func GraphqlHandler(orm *orm.GormDatabase, store objectStorage.ObjectStorage, config config.AuthEndpoint, photoService config.PhotoService, logger *zap.Logger) gin.HandlerFunc {
	h := auth.Middleware(orm,
		handler.GraphQL(gql.NewExecutableSchema(resolvers.NewRootResolvers(orm, store, config, photoService, logger))),
	)

	return func(c *gin.Context) {

		h.ServeHTTP(c.Writer, c.Request)
	}
}

// PlaygroundHandler Defines the Playground handler to expose our playground
func PlaygroundHandler(path string) gin.HandlerFunc {
	h := handler.Playground("Go GraphQL Server", path)
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
