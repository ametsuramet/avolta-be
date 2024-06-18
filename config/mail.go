package config

type MailConfig struct {
	SmtpHost        string
	SmtpPort        int
	SmtpUsername    string
	SmtpPassword    string
	SmtpSenderEmail string
}
