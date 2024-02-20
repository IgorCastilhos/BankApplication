package mail

import (
    "github.com/IgorCastilhos/BankApplication/utils"
    "github.com/stretchr/testify/require"
    "testing"
)

func TestSendEmailWithGmail(t *testing.T) {
    if testing.Short() {
        t.Skip()
    }
    
    config, err := utils.LoadConfig("..")
    require.NoError(t, err)
    
    sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
    subject := "A test email"
    content := `
    <h1>Hello World<h1>
    <p>This is a test message from <a href="https://instagram.com/igor_paprocki_dev/">Igor Castilhos</a></p>
    `
    
    to := []string{"igorcastilhos2020@gmail.com"}
    attachFiles := []string{"../README.md"}
    
    err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
    require.NoError(t, err)
}
