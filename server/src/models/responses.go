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

type RegistrationSummaryMetrics struct {
	TotalRegistrations        int                `json:"total_registrations" db:"total_registrations"`
	SuccessfulRegistrations   int                `json:"successful_registrations" db:"succesfull_registrations"`
	FailedRegistrations       int                `json:"failed_registrations" db:"failed_registrations"`
	AverageRegistrationTimeMs float64            `json:"average_registration_time" db:"average_registration_time"`
	MethodDistribution        MethodDistribution `json:"method_distribution"`
	FederatedProviders        map[string]int     `json:"federated_providers" db:"federated_providers"`
}

type LocationMetric struct {
	Country string `json:"country" db:"country"`
	Amount  int    `json:"amount" db:"amount"`
}

type LocationMetrics struct {
	Locations []LocationMetric `json:"locations"`
}
