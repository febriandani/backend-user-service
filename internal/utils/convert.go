package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/febriandani/backend-user-service/internal/utils/constant/general"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func GetInt(x string) int {
	i, err := strconv.Atoi(x)
	if err != nil {
		fmt.Println("utils -> GentInt : error, ", err)
		fmt.Println("Can't convert into Integer")
		fmt.Println("Please re-check .env, you recently input")
		fmt.Println(x)
	}

	return i
}

func GetFloat(x string) float32 {
	i, err := strconv.ParseFloat(x, 32)
	if err != nil {
		fmt.Println("utils -> GetFloat : error, ", err)
		fmt.Println("Can't convert into float32")
		fmt.Println("Please re-check .env, you recently input")
		fmt.Println(x)
	}

	return float32(i)
}

func ToFormatTime(datetime string) (string, error) {
	t, err := time.Parse(general.DBTimeLayout, datetime)
	if err != nil {
		return datetime, err
	}

	tString := t.Format(general.ResponseTimeLayout)

	return tString, nil
}

func GetTimeString() string {
	t := time.Now().UTC()
	tString := t.Format(general.DBTimeLayout)

	return tString
}

func StrToInt(data string) (int, error) {
	return strconv.Atoi(data)
}

func StrToInt64(data string) (int64, error) {
	return strconv.ParseInt(data, 10, 64)
}

func StrToFloat64(data string) (float64, error) {
	return strconv.ParseFloat(data, 64)
}

func Int64sJoin(data []int64) string {
	s, err := json.Marshal(data)
	if err != nil {
		return ""
	}

	return strings.Trim(string(s), "[]")
}

func GetDataFromKey(data, key string) (string, error) {
	if data == "" || key == "" {
		return "", errors.New("data/key cannot be empty")
	}

	result, err := GetDecrypt([]byte(key), data)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", result), nil
}

func GetKeyData(data, key string) (string, error) {
	if data == "" || key == "" {
		return "", errors.New("data/key cannot be empty")
	}

	result, err := GetEncrypt([]byte(key), fmt.Sprintf("%v", data))
	if err != nil {
		return "", err
	}

	return result, nil
}

func StructToString(data interface{}) string {
	result, err := json.Marshal(data)
	if err != nil {
		return ""
	}

	return string(result)
}

func FloatToRupiah(price float64) string {
	p := message.NewPrinter(language.Indonesian)
	moneyString := p.Sprintf("Rp %.2f", price)

	return moneyString
}

func ConvertIDs(ids string) ([]int64, error) {
	var result []int64
	idString := strings.Split(ids, ",")

	for _, val := range idString {
		id, err := StrToInt64(val)
		if err != nil {
			return result, err
		}

		result = append(result, id)
	}

	return result, nil
}

func ArrInt64Join(ids []int64, separator string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(ids), " ", separator, -1), "[]")
}

func StrToArrInt64(data, separator string) ([]int64, error) {
	splitData := strings.Split(data, separator)

	var result []int64
	for _, dt := range splitData {
		convData, err := StrToInt64(dt)
		if err != nil {
			return result, err
		}

		result = append(result, convData)
	}

	return result, nil
}

func StrToArrMapInt64(data, separator string) (map[int64]int64, error) {
	splitData := strings.Split(data, separator)

	result := make(map[int64]int64)
	for _, dt := range splitData {
		convData, err := StrToInt64(dt)
		if err != nil {
			return result, err
		}

		result[convData] = convData
	}

	return result, nil
}

func StrToArrMapString(data, separator string) (map[string]string, error) {
	splitData := strings.Split(data, separator)

	result := make(map[string]string)
	for _, dt := range splitData {
		result[dt] = dt
	}

	return result, nil
}

func FormatPhoneNumber(phone string) string {
	// Remove any leading spaces or plus sign
	phone = strings.TrimLeft(phone, " +")

	// Check if the phone number starts with "62"
	if strings.HasPrefix(phone, "62") {
		return phone
	}

	// Check if the phone number starts with "+62"
	if strings.HasPrefix(phone, "+62") {
		// Remove the "+" character and return the number
		return strings.TrimPrefix(phone, "+")
	}

	if strings.Contains(phone, "0") {
		// Replace the "0" with "62" as the prefix and return the number
		return fmt.Sprintf("62%s", strings.Replace(phone, "0", "", 1))
	}

	// If the phone number doesn't start with "62" or "+62", add "62" as the prefix
	return fmt.Sprintf("62%s", phone)
}
