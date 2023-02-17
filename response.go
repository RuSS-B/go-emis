package emis

type EmisResponse struct {
	Status            int    `json:"status"`
	StatusDesc        string `json:"statusDesc"`
	Ts                string `json:"ts"`
	ServiceURL        string `json:"serviceURL"`
	ProcessedEntries  int    `json:"processedEntries"`
	DataSupplierID    string `json:"dataSupplierID"`
	DataSupplierLogin string `json:"dataSupplierLogin"`
}
