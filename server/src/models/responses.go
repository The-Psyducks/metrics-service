package models

type MethodDistribution struct {
	EmailAndPassword int `json:"email_and_password" db:"email"`
	Federated        int `json:"federated" db:"federated"`
}

type LoginSummaryMetrics struct {
	TotalLogins        int                `json:"total_logins" db:"total_logins"`
	SuccessfulLogins   int                `json:"successful_logins" db:"succesfull_logins"`
	FailedLogins       int                `json:"failed_logins" db:"failed_logins"`
	MethodDistribution MethodDistribution `json:"method_distribution"`
	FederatedProviders map[string]int     `json:"federated_providers" db:"federated_providers"`
}
