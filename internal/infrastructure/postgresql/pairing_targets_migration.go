package infrastructure

import (
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

func migrateAndDropPairingTargets(db *gorm.DB) error {
	if !db.Migrator().HasTable("pairing_targets") {
		return nil
	}

	type extraTarget struct {
		PairingID uint
		TargetID  uint
	}
	var extras []extraTarget
	if err := db.Raw(`
		SELECT pt.pairing_id, pt.target_id
		FROM pairing_targets pt
		INNER JOIN pairings p ON p.id = pt.pairing_id AND p.deleted_at IS NULL
		WHERE pt.deleted_at IS NULL
		  AND pt.kind = 'user'
		  AND pt.target_id <> p.mentee_id
	`).Scan(&extras).Error; err != nil {
		return err
	}

	for _, ex := range extras {
		var source entity.Pairing
		if err := db.Where("id = ?", ex.PairingID).First(&source).Error; err != nil {
			return err
		}
		var existing int64
		if err := db.Model(&entity.Pairing{}).
			Where("mentor_id = ? AND mentee_id = ? AND is_temporary = ?", source.MentorID, ex.TargetID, false).
			Count(&existing).Error; err != nil {
			return err
		}
		if existing > 0 {
			continue
		}
		clone := entity.Pairing{
			MentorID:        source.MentorID,
			MenteeID:        ex.TargetID,
			CreatedByUserID: source.CreatedByUserID,
			IsTemporary:     source.IsTemporary,
		}
		if err := db.Create(&clone).Error; err != nil {
			return err
		}
	}

	return db.Exec("DROP TABLE IF EXISTS pairing_targets CASCADE").Error
}
