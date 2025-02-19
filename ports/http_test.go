package ports

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	app "github.com/alvarezjulia/meisterwerk-catalog/internal/application"
	"github.com/alvarezjulia/meisterwerk-catalog/internal/application/command/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHTTPServer_CreateProduct(t *testing.T) {
	tests := []struct {
		name         string
		input        map[string]interface{}
		setupMock    func(*mocks.MockCreateProduct)
		expectedCode int
	}{
		{
			name: "successful creation",
			input: map[string]interface{}{
				"name":        "Test Product",
				"description": "Test Description",
				"price":       99.99,
			},
			setupMock: func(m *mocks.MockCreateProduct) {
				m.EXPECT().
					Handle(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "bad request - invalid input (400)",
			input: map[string]interface{}{
				"product_name": "",
			},
			setupMock:    func(m *mocks.MockCreateProduct) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "bad request - malformed JSON (400)",
			input: map[string]interface{}{
				"product_price": "not_a_number",
			},
			setupMock:    func(m *mocks.MockCreateProduct) {},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCmd := mocks.NewMockCreateProduct(ctrl)
			tt.setupMock(mockCmd)

			app := &app.App{
				Commands: &app.Commands{
					CreateProduct: mockCmd,
				},
			}

			server := NewHTTPServer(app)

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/product", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			server.CreateProduct(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}
