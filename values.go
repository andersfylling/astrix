package astrix

func DefaultValue(t primitiveType) (result string, success bool) {
	switch t {
	case TypeInt, TypeInt8, TypeInt16, TypeInt32, TypeInt64:
		result = "0"
	case TypeUint, TypeUint8, TypeUint16, TypeUint32, TypeUint64:
		result = "0"
	case TypeFloat32, TypeFloat64:
		result = "0.0"
	case TypeString:
		result = ""
	case TypeBool:
		result = "false"
	case TypeSlice, TypeMap:
		result = "nil"
	case TypeInterface, TypePointer:
		result = "nil"
	default:
		return "", false
	}
	// TODO: struct

	return result, true
}
