package libhac

import "encoding/xml"

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

type CNMTXML struct {
	XMLName                       xml.Name          `xml:"ContentMeta"`
	Type                          string            `xml:"Type"`
	ID                            string            `xml:"Id"`
	Version                       string            `xml:"Version"`
	RequiredDownloadSystemVersion string            `xml:"RequiredDownloadSystemVersion"`
	ContentEntries                []ContentEntryXML `xml:"Content"`
	Digest                        string            `xml:"Digest"`
	KeyGenerationMin              string            `xml:"KeyGenerationMin"`
	RequiredSystemVersion         string            `xml:"RequiredSystemVersion"`
	PatchID                       string            `xml:"PatchId"`
}

type ContentEntryXML struct {
	Type          string `xml:"Type"`
	ID            string `xml:"Id"`
	Size          string `xml:"Size"`
	Hash          string `xml:"Hash"`
	KeyGeneration string `xml:"KeyGeneration"`
}
