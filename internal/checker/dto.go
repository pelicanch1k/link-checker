package checker

type CheckLinksInput struct {
	URLs []string
}

type CheckLinksOutput struct {
	TaskID   int
	Links    map[string]string
}

type CheckLinksByIDsInput struct {
	LinksList []int
}

type CheckLinksByIDsOutput struct {
	PDFData []byte
}