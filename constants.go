package otr4

import (
	"github.com/twstrike/ed448"
)

const (
	fieldBytes  = 56
	keySigBytes = 144
	symKeyBytes = 32
	sigBytes    = 112
	dsaSigBytes = 40
	// PublicKeySize is the size, in bytes, of public keys.
	publicKeySize = 57
	// PrivateKeySize is the size, in bytes, of private keys.
	privateKeySize = 57
	// SignatureSize is the size, in bytes, of signatures generated and verified.
	signatureSize = 114
	mask          = 0x80
)

var (
	g2 = ed448.NewPoint(
		[16]uint32{
			0x0cf14237, 0x0ac97f43, 0x0a9543bc, 0x0dc98db8,
			0x0bcca6a6, 0x07874a17, 0x021af78f, 0x0fffa763,
			0x0cf2ac0b, 0x074f2a89, 0x0f89f88d, 0x0356a31e,
			0x09f61e5a, 0x00c01083, 0x0c84b7a5, 0x00bf3b5c,
		},
		[16]uint32{
			0x00c9a64c, 0x06b775bc, 0x026148bb, 0x0ee0c3e1,
			0x0303aa98, 0x04fad09b, 0x0efaf59d, 0x03008555,
			0x072a0bf6, 0x023bc0fa, 0x0c52ee5b, 0x0f0f61f9,
			0x05cf8d7f, 0x0b8b7f38, 0x018a4398, 0x06a9849a,
		},
		[16]uint32{
			0x014e2fce, 0x0198c24c, 0x0b74b290, 0x0080f748,
			0x0fb60b6e, 0x08ab2f53, 0x06c32b60, 0x06979188,
			0x0e87a66d, 0x087ecac7, 0x0f354ebd, 0x035faebf,
			0x0e30d07f, 0x0c96f513, 0x0fab82ed, 0x0da28e58,
		},
		[16]uint32{
			0x0702239a, 0x05c67537, 0x0ce76a54, 0x0fae388e,
			0x034bcae9, 0x06b5fe3d, 0x0d3c37ae, 0x09cac77d,
			0x0761a657, 0x0a02246f, 0x06490757, 0x09448b04,
			0x05281bbe, 0x0e0bd3d4, 0x0abc5ecb, 0x07c655f9,
		},
	)
)
