package ws

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/malikkhoiri/csms/internal/domain"
)

var upgrader = websocket.Upgrader{
	Subprotocols: []string{"ocpp1.6"},
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const CallResult uint16 = 3

type OCPPHandler struct {
	chargePointService domain.ChargePointService
	transactionService domain.TransactionService
	userService        domain.UserService
	connectorService   domain.ConnectorService
	idTagService       domain.IDTagService
}

func NewOCPPHandler(
	chargePointService domain.ChargePointService,
	transactionService domain.TransactionService,
	userService domain.UserService,
	connectorService domain.ConnectorService,
	idTagService domain.IDTagService,
) *OCPPHandler {
	return &OCPPHandler{
		chargePointService: chargePointService,
		transactionService: transactionService,
		userService:        userService,
		connectorService:   connectorService,
		idTagService:       idTagService,
	}
}

func (h *OCPPHandler) HandleWebSocket(c *gin.Context) {
	cpCode := strings.TrimPrefix(c.Param("cpID"), "/")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	log.Printf("Connected CP: %s", cpCode)
	log.Println("Subprotocol:", conn.Subprotocol())

	h.handleOCPPMessages(conn, cpCode)
}

func (h *OCPPHandler) handleOCPPMessages(conn *websocket.Conn, cpCode string) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			return
		}

		h.processOCPPMessage(conn, msg, cpCode)
	}
}

func (h *OCPPHandler) processOCPPMessage(conn *websocket.Conn, msg []byte, cpCode string) {
	var ocppMsg []interface{}
	if err := json.Unmarshal(msg, &ocppMsg); err != nil {
		log.Println("Invalid JSON:", err)
		return
	}

	// Expecting a CALL message [2, messageId, action, payload]
	if len(ocppMsg) >= 4 && int(ocppMsg[0].(float64)) == 2 {
		messageID := ocppMsg[1].(string)
		action := ocppMsg[2].(string)
		payload := ocppMsg[3].(map[string]interface{})

		log.Printf("OCPP CALL received: ID=%s, Action=%s", messageID, action)

		h.handleOCPPAction(conn, action, messageID, payload, cpCode)
	}
}

func (h *OCPPHandler) handleOCPPAction(conn *websocket.Conn, action, messageID string, payload map[string]interface{}, cpCode string) {
	switch action {
	case "BootNotification":
		h.handleBootNotification(conn, messageID, payload, cpCode)
	case "Heartbeat":
		h.handleHeartbeat(conn, messageID, cpCode)
	case "Authorize":
		h.handleAuthorize(conn, messageID, payload, cpCode)
	case "StartTransaction":
		h.handleStartTransaction(conn, messageID, payload, cpCode)
	case "StopTransaction":
		h.handleStopTransaction(conn, messageID, payload, cpCode)
	case "StatusNotification":
		h.handleStatusNotification(conn, messageID, payload, cpCode)
	case "MeterValues":
		h.handleMeterValues(conn, messageID, payload, cpCode)
	default:
		log.Printf("Unhandled OCPP action: %s", action)
	}
}

func (h *OCPPHandler) handleBootNotification(conn *websocket.Conn, messageID string, payload map[string]interface{}, cpCode string) {
	ctx := context.Background()

	request := &domain.BootNotificationRequest{
		ChargePointVendor:       getString(payload, "chargePointVendor"),
		ChargePointModel:        getString(payload, "chargePointModel"),
		ChargePointSerialNumber: getString(payload, "chargePointSerialNumber"),
		ChargeBoxSerialNumber:   getString(payload, "chargeBoxSerialNumber"),
		FirmwareVersion:         getString(payload, "firmwareVersion"),
		Iccid:                   getString(payload, "iccid"),
		Imsi:                    getString(payload, "imsi"),
		MeterType:               getString(payload, "meterType"),
		MeterSerialNumber:       getString(payload, "meterSerialNumber"),
	}

	response, err := h.chargePointService.RegisterChargePoint(ctx, request, cpCode)
	if err != nil {
		log.Printf("Error registering charge point: %v", err)
		return
	}

	ocppResponse := []interface{}{
		CallResult,
		messageID,
		response,
	}
	replyBytes, _ := json.Marshal(ocppResponse)
	conn.WriteMessage(websocket.TextMessage, replyBytes)
	log.Println("Sent BootNotification response")
}

func (h *OCPPHandler) handleHeartbeat(conn *websocket.Conn, messageID string, cpCode string) {
	ctx := context.Background()

	chargePoint, err := h.chargePointService.GetChargePointByCode(ctx, cpCode)
	if err != nil {
		log.Printf("Charge point not found for heartbeat: %s", cpCode)
		return
	}

	if err := h.chargePointService.UpdateHeartbeat(ctx, chargePoint.ID); err != nil {
		log.Printf("Error updating heartbeat: %v", err)
	}

	response := map[string]interface{}{
		"currentTime": time.Now().UTC().Format(time.RFC3339),
	}
	ocppResponse := []interface{}{
		CallResult,
		messageID,
		response,
	}
	replyBytes, _ := json.Marshal(ocppResponse)
	conn.WriteMessage(websocket.TextMessage, replyBytes)
	log.Println("Sent Heartbeat response")
}

func (h *OCPPHandler) handleAuthorize(conn *websocket.Conn, messageID string, payload map[string]interface{}, cpCode string) {
	ctx := context.Background()

	request := &domain.AuthorizeRequest{
		IDTag: getString(payload, "idTag"),
	}

	response, err := h.idTagService.Authorize(ctx, request)
	if err != nil {
		log.Printf("Error authorizing user: %v", err)
		return
	}

	ocppResponse := []interface{}{
		CallResult,
		messageID,
		response,
	}
	replyBytes, _ := json.Marshal(ocppResponse)
	conn.WriteMessage(websocket.TextMessage, replyBytes)
	log.Println("Sent Authorize response")
}

func (h *OCPPHandler) handleStartTransaction(conn *websocket.Conn, messageID string, payload map[string]interface{}, cpCode string) {
	ctx := context.Background()

	chargePoint, err := h.chargePointService.GetChargePointByCode(ctx, cpCode)
	if err != nil {
		log.Printf("Charge point not found for start transaction: %s", cpCode)
		return
	}

	meterStart := int(payload["meterStart"].(float64))
	request := &domain.StartTransactionRequest{
		ConnectorId: int(payload["connectorId"].(float64)),
		IDTag:       getString(payload, "idTag"),
		MeterStart:  meterStart,
	}

	response, err := h.transactionService.StartTransaction(ctx, request, chargePoint.ID)
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return
	}

	ocppResponse := []interface{}{
		CallResult,
		messageID,
		response,
	}
	replyBytes, _ := json.Marshal(ocppResponse)
	conn.WriteMessage(websocket.TextMessage, replyBytes)
	log.Println("Sent StartTransaction response")
}

func (h *OCPPHandler) handleStopTransaction(conn *websocket.Conn, messageID string, payload map[string]interface{}, cpCode string) {
	ctx := context.Background()

	chargePoint, err := h.chargePointService.GetChargePointByCode(ctx, cpCode)
	if err != nil {
		log.Printf("Charge point not found for stop transaction: %s", cpCode)
		return
	}

	request := &domain.StopTransactionRequest{
		TransactionId: int(payload["transactionId"].(float64)),
		IDTag:         getString(payload, "idTag"),
		MeterStop:     int(payload["meterStop"].(float64)),
		Reason:        getString(payload, "reason"),
	}

	response, err := h.transactionService.StopTransaction(ctx, request, chargePoint.ID)
	if err != nil {
		log.Printf("Error stopping transaction: %v", err)
		return
	}

	ocppResponse := []interface{}{
		CallResult,
		messageID,
		response,
	}
	replyBytes, _ := json.Marshal(ocppResponse)
	conn.WriteMessage(websocket.TextMessage, replyBytes)
	log.Println("Sent StopTransaction response")
}

func (h *OCPPHandler) handleStatusNotification(conn *websocket.Conn, messageID string, payload map[string]interface{}, cpCode string) {
	ctx := context.Background()

	chargePoint, err := h.chargePointService.GetChargePointByCode(ctx, cpCode)
	if err != nil {
		log.Printf("Charge point not found for status notification: %s", cpCode)
		return
	}

	request := &domain.StatusNotificationRequest{
		ConnectorId:     int(payload["connectorId"].(float64)),
		Status:          getString(payload, "status"),
		ErrorCode:       getString(payload, "errorCode"),
		Info:            getString(payload, "info"),
		VendorId:        getString(payload, "vendorId"),
		VendorErrorCode: getString(payload, "vendorErrorCode"),
	}

	if err := h.connectorService.UpdateConnectorStatus(ctx, request, chargePoint.ID); err != nil {
		log.Printf("Error updating connector status: %v", err)
	}

	ocppResponse := []interface{}{
		CallResult,
		messageID,
		map[string]interface{}{},
	}
	replyBytes, _ := json.Marshal(ocppResponse)
	conn.WriteMessage(websocket.TextMessage, replyBytes)
	log.Println("Sent StatusNotification response")
}

func (h *OCPPHandler) handleMeterValues(conn *websocket.Conn, messageID string, payload map[string]interface{}, cpCode string) {
	ctx := context.Background()

	chargePoint, err := h.chargePointService.GetChargePointByCode(ctx, cpCode)
	if err != nil {
		log.Printf("Charge point not found for meter values: %s", cpCode)
		return
	}

	request := &domain.MeterValuesRequest{
		ConnectorId: int(payload["connectorId"].(float64)),
		MeterValue:  []domain.MeterValue{},
	}

	if transactionId, exists := payload["transactionId"]; exists && transactionId != nil {
		tid := int(transactionId.(float64))
		request.TransactionId = &tid
	}

	if err := h.transactionService.UpdateMeterValues(ctx, request, chargePoint.ID); err != nil {
		log.Printf("Error updating meter values: %v", err)
	}

	ocppResponse := []interface{}{
		CallResult,
		messageID,
		map[string]interface{}{},
	}
	replyBytes, _ := json.Marshal(ocppResponse)
	conn.WriteMessage(websocket.TextMessage, replyBytes)
	log.Println("Sent MeterValues response")
}

func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}
