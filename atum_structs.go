package libhac

type CNMT struct {
	Path                          string
	Type                          string
	ID                            string
	Version                       string
	RequiredSystemVersion         string
	RequiredDownloadSystemVersion string
	Digest                        string
	MasterKeyRevision             string
	ContentEntries                []ContentEntry
}

type ContentEntry struct {
	Hash string
	ID   string
	Size string
	Type string
}
