package utils

func PtrBool(b bool) *bool {
	return &b
}

func PtrInt(i int) *int {
	return &i
}

func PtrString(s string) *string {
	return &s
}