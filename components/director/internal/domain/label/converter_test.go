package label_test

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/kyma-incubator/compass/components/director/internal/domain/label"
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConverter_ToEntity(t *testing.T) {
	stringValue := "foo"
	marshalledStringValue, err := json.Marshal(stringValue)
	require.NoError(t, err)

	arrayValue := []interface{}{"foo", "bar"}
	marshalledArrayValue, err := json.Marshal(arrayValue)
	require.NoError(t, err)

	// given
	testCases := []struct {
		Name               string
		Input              model.Label
		Expected           label.Entity
		ExpectedErrMessage string
	}{
		{
			Name:               "All properties given",
			Input:              fixLabelModel("1", arrayValue),
			Expected:           fixLabelEntity("1", marshalledArrayValue),
			ExpectedErrMessage: "",
		},
		{
			Name:               "String value",
			Input:              fixLabelModel("1", stringValue),
			Expected:           fixLabelEntity("1", marshalledStringValue),
			ExpectedErrMessage: "",
		},
		{
			Name: "Empty value",
			Input: model.Label{
				ID:     "2",
				Key:    "foo",
				Tenant: "tenant",
			},
			Expected: label.Entity{
				ID:       "2",
				Key:      "foo",
				TenantID: "tenant",
			},
			ExpectedErrMessage: "",
		},
		{
			Name: "Error",
			Input: model.Label{
				Value: make(chan int),
			},
			Expected:           label.Entity{},
			ExpectedErrMessage: "while marshalling Value: json: unsupported type: chan int",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			conv := label.NewConverter()

			// when
			res, err := conv.ToEntity(testCase.Input)

			if testCase.ExpectedErrMessage != "" {
				require.Error(t, err)
				assert.Equal(t, testCase.ExpectedErrMessage, err.Error())
			}

			// then
			assert.Equal(t, testCase.Expected, res)
		})
	}
}

func TestConverter_FromEntity(t *testing.T) {
	stringValue := "foo"
	marshalledStringValue, err := json.Marshal(stringValue)
	require.NoError(t, err)

	arrayValue := []interface{}{"foo", "bar"}
	marshalledArrayValue, err := json.Marshal(arrayValue)
	require.NoError(t, err)

	// given
	testCases := []struct {
		Name               string
		Input              label.Entity
		Expected           model.Label
		ExpectedErrMessage string
	}{
		{
			Name:               "All properties given",
			Input:              fixLabelEntity("1", marshalledArrayValue),
			Expected:           fixLabelModel("1", arrayValue),
			ExpectedErrMessage: "",
		},
		{
			Name:               "String value",
			Input:              fixLabelEntity("1", marshalledStringValue),
			Expected:           fixLabelModel("1", stringValue),
			ExpectedErrMessage: "",
		},
		{
			Name: "Empty value",
			Input: label.Entity{
				ID:       "2",
				Key:      "foo",
				TenantID: "tenant",
			},
			Expected: model.Label{
				ID:     "2",
				Key:    "foo",
				Tenant: "tenant",
			},
			ExpectedErrMessage: "",
		},
		{
			Name:               "Error",
			Input:              fixLabelEntity("1", []byte("{json")),
			Expected:           model.Label{},
			ExpectedErrMessage: "while unmarshalling Value: invalid character 'j' looking for beginning of object key string",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			conv := label.NewConverter()

			// when
			res, err := conv.FromEntity(testCase.Input)

			if testCase.ExpectedErrMessage != "" {
				require.Error(t, err)
				assert.Equal(t, testCase.ExpectedErrMessage, err.Error())
			}

			// then
			assert.Equal(t, testCase.Expected, res)
		})
	}
}

func fixLabelEntity(id string, value []byte) label.Entity {
	return label.Entity{
		ID:       id,
		TenantID: "tenant",
		AppID:    sql.NullString{},
		RuntimeID: sql.NullString{
			String: "321",
			Valid:  true,
		},
		Key:   "test",
		Value: string(value),
	}
}

func fixLabelModel(id string, value interface{}) model.Label {
	return model.Label{
		ID:         id,
		Tenant:     "tenant",
		Key:        "test",
		ObjectType: model.RuntimeLabelableObject,
		ObjectID:   "321",
		Value:      value,
	}
}
