package main

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/bagaking/ankibuild/apkg"
)

// updateCardRuntime Updates the runtime information of the configuration with the IDs of the created note and card.
func updateCardRuntime(note *apkg.Note, card *apkg.Card, cardConf QnACard) (*QnACard, error) {
	if cardConf.Question != note.Front() {
		return nil, fmt.Errorf("write card runtime back failed, cardConf= %v", cardConf)
	}
	cardConf.Runtime = CardRuntime{
		NoteID: note.ID,
		CardID: card.ID,
	}
	return &cardConf, nil
}

// writeRuntimeBack writes the runtime information back into the original .apkg.toml configuration file.
// writeRuntimeBack writes the runtime information back into the original .apkg.toml configuration file with a backup.
func writeRuntimeBack(conf *KnowledgePage, filePath string) error {
	// Generate the backup file name with timestamp
	timeStamp := time.Now().Format("20060102150405")
	backupFilePath := filePath + "." + timeStamp + ".bak"

	// Read original file
	original, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Write to the backup file
	if err = os.WriteFile(backupFilePath, original, 0644); err != nil {
		return err
	}

	// Encode the updated configuration
	var buf bytes.Buffer
	if err = toml.NewEncoder(&buf).Encode(conf); err != nil {
		return err
	}

	// Write back to the original file
	return os.WriteFile(filePath, buf.Bytes(), 0644)
}
