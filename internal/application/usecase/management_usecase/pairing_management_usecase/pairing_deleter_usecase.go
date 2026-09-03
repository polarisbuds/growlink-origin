package pairing_mng_usecase

import repository_interface "growlink/internal/domain/interface/repository"

type PairingDeleter struct {
	pairingRepo repository_interface.PairingRepository
}

func NewPairingDeleter(pairingRepo repository_interface.PairingRepository) *PairingDeleter {
	return &PairingDeleter{
		pairingRepo: pairingRepo,
	}
}
