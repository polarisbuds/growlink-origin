package pairing_mng_usecase

import repository_interface "growlink/internal/domain/interface/repository"

type PairingUpdater struct {
	pairingRepo repository_interface.PairingRepository
}

func NewPairingUpdater(pairingRepo repository_interface.PairingRepository) *PairingUpdater {
	return &PairingUpdater{
		pairingRepo: pairingRepo,
	}
}
