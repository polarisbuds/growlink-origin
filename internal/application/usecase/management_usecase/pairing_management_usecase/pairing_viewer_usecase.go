package pairing_mng_usecase

import repository_interface "growlink/internal/domain/interface/repository"

type PairingViewer struct {
	pairingRepo repository_interface.PairingRepository
}

func NewPairingViewer(pairingRepo repository_interface.PairingRepository) *PairingViewer {
	return &PairingViewer{
		pairingRepo: pairingRepo,
	}
}
