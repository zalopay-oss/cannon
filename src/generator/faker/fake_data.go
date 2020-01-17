package faker

import (
	"encoding/json"
	"github.com/jhump/protoreflect/desc"
	"github.com/tranndc/benchmark/generator/parser"
	"github.com/tranndc/benchmark/generator/random"
	"math/rand"
	"strings"
)

func FakeData(call, inputMessage string, fileDesc *desc.FileDescriptor) (string, error) {
	_ = random.SetRandomNumberBoundaries(1,100)
	_ = random.SetRandomStringLength(30)
	_ = random.SetRandomMapAndSliceSize(5)

	files := map[string]*desc.FileDescriptor{}
	files[fileDesc.GetName()] = fileDesc

	msgList := fileDesc.GetMessageTypes()
	svc, _, _  := parser.ParseServiceMethod(call)
	pos := strings.LastIndex(svc, ".")
	serviceClass := svc[:pos]

	var jsonStr string = ""
	for _, msg := range(msgList) {
		if msg.GetName() == inputMessage {
			structObj := make(map[string]interface{})
			result := FakeDataForMessage(&structObj, serviceClass, msg, fileDesc)
			empData, err := json.Marshal(result)
			if err != nil {
				return jsonStr, err
			}
			jsonStr = string(empData)
			break
		}
	}
	return jsonStr, nil
}

func FakeDataForMessage(objRef *map[string]interface{}, serviceClass string, msg *desc.MessageDescriptor, fileDesc *desc.FileDescriptor) map[string]interface{} {
	obj := *objRef
	fields := msg.GetFields()

	// Get one field in OneOf
	oneOfFields := getOneOfNamesList(msg)

	// Get NestedMessages
	nestedMsgList := msg.GetNestedMessageTypes()
	nestedMsgNames := make(map[string]int)
	if nestedMsgList != nil && len(nestedMsgList) > 0 {
		for _, nst := range(nestedMsgList) {
			nestedMsgNames[nst.GetName()] = 1
		}
	}

	for _, f := range(fields) {
		msgType := f.GetType().String()
		label := f.GetLabel().String()
		if !checkInOneOfList(f.GetName(), oneOfFields) && f.GetOneOf() != nil {
			continue
		}
		if msgType == "TYPE_ENUM" {
			// fmt.Println(f.GetEnumType().GetName())
			obj[f.GetName()] = getValueForEnum(f, serviceClass, fileDesc)
		} else if label == "LABEL_REPEATED" && msgType == "TYPE_MESSAGE" {
			nestedMsg := f.GetMessageType()
			// Check if it is a NestedMessageType
			if _, ok := nestedMsgNames[nestedMsg.GetName()]; ok {
				obj[f.GetName()] = getValueNestedFields(f, serviceClass, fileDesc)
			} else {
				// If not a map then this field is a list/slice/array
				obj[f.GetName()] = getValueForList(msgType, f, serviceClass, fileDesc)
			}
		} else if label == "LABEL_REPEATED" {
			obj[f.GetName()] = getValueForList(msgType, f, serviceClass, fileDesc)
		} else {
			obj[f.GetName()] = getValueNormalFields(f, msg, serviceClass, fileDesc)
		}
	}
	return *objRef
}

func getValueNestedFields(field *desc.FieldDescriptor, serviceClass string, fileDesc *desc.FileDescriptor) interface{} {
	msg := field.GetMessageType()

	if msg.GetMessageOptions().GetMapEntry() {
		mapLen := random.RandomSliceAndMapSize()
		keyType := field.GetMapKeyType()
		valType := field.GetMapValueType()
		mapObj := make(map[string]interface{})

		for i:=0; i < mapLen; i++ {
			keyVal := getValueNormalFields(keyType, msg, serviceClass, fileDesc)
			valVal := getValueNormalFields(valType, msg, serviceClass, fileDesc)
			mapObj[keyVal.(string)] = valVal
		}
		return mapObj
	} else {
		objList := getValueForList(field.GetType().String(), field, serviceClass, fileDesc)
		return objList
	}
	return nil
}

func getValueNormalFields(
	field *desc.FieldDescriptor,
	msg *desc.MessageDescriptor,
	serviceClass string,
	fileDesc *desc.FileDescriptor) interface{} {

	nestedMsgNames := make(map[string]int)
	if msg != nil {
		nestedMsgList := msg.GetNestedMessageTypes()
		if nestedMsgList != nil && len(nestedMsgList) > 0 {
			for _, nst := range(nestedMsgList) {
				nestedMsgNames[nst.GetName()] = 1
			}
		}
	}
	msgType := field.GetType().String()

	fieldType := field.GetType().String()
	label := field.GetLabel().String()

	if msgType == "TYPE_ENUM" {
		enumVal := getValueForEnum(field, serviceClass, fileDesc)
		return enumVal
	} else if label == "LABEL_REPEATED" && msgType == "TYPE_MESSAGE" {
		nestedMsg := field.GetMessageType()
		// Check if it is a NestedMessageType
		if _, ok := nestedMsgNames[nestedMsg.GetName()]; ok {
			return getValueNestedFields(field, serviceClass, fileDesc)
		} else {
			// If not a map then this field is a list/slice/array
			return getValueForList(msgType, field, serviceClass, fileDesc)
		}
	} else if label == "LABEL_REPEATED" {
		objList := getValueForList(fieldType, field, serviceClass, fileDesc)
		// fmt.Println(objList)
		return objList
	} else if fieldType == "TYPE_MESSAGE" {
		// childObj := make(map[string]interface{})
		m := fileDesc.FindMessage(serviceClass + "." + field.GetMessageType().GetName())
		fieldList := m.GetFields()
		if m != nil {
			msgValue := make(map[string]interface{})
			for _, field := range(fieldList) {
				// Get one field in OneOf
				oneOfFields := getOneOfNamesList(msg)
				if !checkInOneOfList(field.GetName(), oneOfFields) && field.GetOneOf() != nil{
					continue
				}
				msgValue[field.GetName()] = getValueNormalFields(field, m, serviceClass, fileDesc)
			}
			return msgValue
		}
	} else if label == "LABEL_OPTIONAL" {
		value := assignValueScalaType(fieldType)
		return value
	}
	return nil
}

func getValueForList(msgType string, f *desc.FieldDescriptor, serviceClass string, fileDesc *desc.FileDescriptor) []interface{} {
	listType := f.GetType().String()
	numObjs := random.RandomSliceAndMapSize()
	objList := make([]interface{}, numObjs)
	if listType == "TYPE_MESSAGE" {
		m := fileDesc.FindMessage(serviceClass + "." + f.GetMessageType().GetName())
		if m != nil {
			fieldList := m.GetFields()
			for i := 0; i < numObjs; i++ {
				eleObj := make(map[string]interface{})
				for _, field := range(fieldList) {
					oneOfFields := getOneOfNamesList(m)
					if !checkInOneOfList(field.GetName(), oneOfFields) && field.GetOneOf() != nil {
						continue
					}
					element := getValueNormalFields(field, m, serviceClass, fileDesc)
					eleObj[field.GetName()] = element
				}
				objList[i] = eleObj
			}
		}
	} else {
		for i:=0; i < numObjs; i++ {
			element := assignValueScalaType(msgType)
			objList[i] = element
		}
	}
	return objList
}

func getValueForEnum(field *desc.FieldDescriptor, serviceClass string, fileDesc *desc.FileDescriptor) int {
	enumMsg := fileDesc.FindEnum(serviceClass + "." + field.GetEnumType().GetName())
	enumLen := len(enumMsg.GetValues())
	bound := random.NumberBoundary{Start: 1, End: enumLen}
	enumVal := random.RandomIntegerWithBoundary(bound)
	return enumVal
}

func checkInOneOfList(fieldName string, oneOfList []string) bool {
	for _, oneOfName := range(oneOfList) {
		if oneOfName == fieldName {
			return true
		}
	}
	return false
}

func getOneOfNamesList(msg *desc.MessageDescriptor) []string {
	oneOfs := msg.GetOneOfs()
	var oneOfFields []string
	for _, oneOf := range(oneOfs) {
		choices := oneOf.GetChoices()
		choiceIdx := random.RandomIntegerWithBoundary(random.NumberBoundary{Start: 0, End: len(choices)})
		oneOfFields = append(oneOfFields, choices[choiceIdx].GetName())
	}
	return oneOfFields
}

func assignValueScalaType(msgType string) interface{} {
	var value interface{}
	switch msgType {
	case "TYPE_STRING":
		value = random.RandomString()
	case "TYPE_INT32":
		value = int32(random.RandomInteger())
	case "TYPE_INT64":
		value = int64(random.RandomInteger())
	case "TYPE_UINT32":
		value = uint32(random.RandomInteger())
	case "TYPE_UINT64":
		value = uint64(random.RandomInteger())
	case "TYPE_SINT32":
		value = int32(random.RandomInteger())
	case "TYPE_SINT64":
		value = int64(random.RandomInteger())
	case "TYPE_FIXED32":
		value = uint32(random.RandomInteger())
	case "TYPE_FIXED64":
		value = uint64(random.RandomInteger())
	case "TYPE_SFIXED32":
		value = int32(random.RandomInteger())
	case "TYPE_SFIXED64":
		value = int64(random.RandomInteger())
	case "TYPE_FLOAT":
		value = rand.Float64()
	case "TYPE_BOOL":
		value = rand.Intn(2) > 0
	default:
		value = nil
	}
	return value
}