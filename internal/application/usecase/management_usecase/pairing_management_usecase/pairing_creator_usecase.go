package pairing_mng_usecase

import repository_interface "growlink/internal/domain/interface/repository"

type PairingCreator struct {
	pairingRepo repository_interface.PairingRepository
}

func NewPairingCreator(pairingRepo repository_interface.PairingRepository) *PairingCreator {
	return &PairingCreator{
		pairingRepo: pairingRepo,
	}
}
