package product

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opplieam/bb-grpc/protogen/go/product"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type Handler struct {
	ProductService product.ProductServiceClient
	Tracer         trace.Tracer
}

func NewHandler(conn *grpc.ClientConn) *Handler {
	return &Handler{
		ProductService: product.NewProductServiceClient(conn),
		Tracer:         otel.GetTracerProvider().Tracer("Product"),
	}
}

func (h *Handler) GetAllProducts(c *gin.Context) {
	// TODO: remove placeholder user_id and use user_id from context instead
	ctx, span := h.Tracer.Start(c.Request.Context(), "GetAllProducts")
	defer span.End()

	userID := 1
	span.SetAttributes(attribute.Int("user_id", userID))
	result, err := h.ProductService.GetProductsByUser(ctx, &product.GetProductsByUserReq{UserId: uint32(userID)})
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}
