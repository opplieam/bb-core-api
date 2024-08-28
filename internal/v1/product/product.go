package product

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opplieam/bb-grpc/protogen/go/product"
	"google.golang.org/grpc"
)

type Handler struct {
	ProductService product.ProductServiceClient
}

func NewHandler(conn *grpc.ClientConn) *Handler {
	return &Handler{
		ProductService: product.NewProductServiceClient(conn),
	}
}

func (h *Handler) GetAllProducts(c *gin.Context) {
	// TODO: remove placeholder user_id and use user_id from context instead
	result, err := h.ProductService.GetProductsByUser(c, &product.GetProductsByUserReq{UserId: 1})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}
