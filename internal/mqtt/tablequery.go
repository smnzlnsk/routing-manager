package mqtt

import "github.com/smnzlnsk/routing-manager/internal/models"

// HandleTableQueryMessage handles messages from table query topics with variables
// Topic pattern: tablequery/{requesterId}/{routingPolicy}
func handleTableQueryMessage(request models.TableQueryMessageRequest) (models.TableQueryMessageResponse, error) {
	return models.TableQueryMessageResponse{
		RequesterId:   request.RequesterId,
		RoutingPolicy: request.RoutingPolicy,
		Payload:       []byte("{\"result\": \"success\"}"),
	}, nil
}
