package dbworker

import (
    "encoding/json"
)

func ParseJson(events [][]byte) ([]map[string]interface{}, error) {
    var data []map[string]interface{}
    for _, event := range events {
        var parsed map[string]interface{}
        err := json.Unmarshal(event, &parsed)
        if err != nil {
            return nil, err
        }
        data = append(data, parsed)
    }
    return data, nil
}