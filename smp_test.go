package otr4

import (
	"github.com/twstrike/ed448"

	. "gopkg.in/check.v1"
)

func (s *OTR4Suite) Test_SMPSecretGeneration(c *C) {
	aliceFingerprint := hexToBytes("0102030405060708090A0B0C0D0E0F101112" +
		"131415161718191A1B1C1D1E1F202122232425262728292A2B2C2D2E2F30" +
		"3132333435363738393A3B3C3D3E3F40")
	bobFingerprint := hexToBytes("4142434445464748494A4B4C4D4E4F50515253" +
		"5455565758595A5B5C5D5E5F606162636465666768696A6B6C6D6E6F7071" +
		"72737475767778797A7B7C7D7E7F00")
	ssid := hexToBytes("FFF3D1E407346468")
	secret := []byte("user's secret")
	rslt := generateSMPsecret(aliceFingerprint, bobFingerprint, ssid, secret)

	expectedSMPSecret := hexToBytes("c57f90a829917526a94b8ed36f0eea8e676" +
		"190a07f6358682d358bc0471bcba401da479d59926ebd1fdfb371233e319dda35365" +
		"cc141ea5c61dd52dcc0cdcd21")

	c.Assert(rslt, DeepEquals, expectedSMPSecret)
}

func (s *OTR4Suite) Test_GenerateDZKP(c *C) {
	b1 := [56]byte{0x04}
	b2 := [56]byte{0x03}
	b3 := [56]byte{0x01}

	r := ed448.NewScalar(b1[:])
	cc := ed448.NewScalar(b2[:])
	a := ed448.NewScalar(b3[:])

	rslt := generateDZKP(r, cc, a)

	b4 := [56]byte{0x01}
	exp := ed448.NewScalar(b4[:])

	c.Assert(rslt, DeepEquals, exp)
}

func (s *OTR4Suite) Test_GenerateZKP(c *C) {
	b1 := [56]byte{0x04}
	b2 := [56]byte{0x01}

	r := ed448.NewScalar(b1[:])
	a := ed448.NewScalar(b2[:])

	cc, d := generateZKP(r, a, byte(01))

	expC := ed448.NewScalar([]byte{
		0xc4, 0x6b, 0x41, 0x3c, 0xd2, 0x70, 0xc9, 0xcb,
		0x86, 0x68, 0x18, 0x57, 0x67, 0x63, 0x36, 0xf0,
		0xb1, 0x38, 0x07, 0x35, 0x0b, 0x51, 0x90, 0xf5,
		0x0a, 0x1c, 0x21, 0x27, 0xbd, 0xd0, 0x28, 0x0e,
		0xd3, 0x0d, 0x30, 0x4d, 0x5e, 0x72, 0x32, 0x00,
		0xd2, 0x25, 0x94, 0xc0, 0xa7, 0x08, 0x0d, 0xa0,
		0xd2, 0xb5, 0xa4, 0x6e, 0x32, 0xc9, 0x02, 0x24,
	},
	)

	expD := ed448.NewScalar([]byte{
		0x33, 0xd9, 0x16, 0x6f, 0xc0, 0x51, 0xaf, 0x57,
		0xce, 0x26, 0xad, 0x36, 0x0b, 0x5f, 0x36, 0x31,
		0xde, 0xfd, 0xce, 0x79, 0x3e, 0x8a, 0xbe, 0xce,
		0xde, 0x07, 0xa9, 0x55, 0x42, 0x2f, 0xd7, 0xf1,
		0x2c, 0xf2, 0xcf, 0xb2, 0xa1, 0x8d, 0xcd, 0xff,
		0x2d, 0xda, 0x6b, 0x3f, 0x58, 0xf7, 0xf2, 0x5f,
		0x2d, 0x4a, 0x5b, 0x91, 0xcd, 0x36, 0xfd, 0x1b,
	},
	)

	c.Assert(cc, DeepEquals, expC)
	c.Assert(d, DeepEquals, expD)
}

func (s *OTR4Suite) Test_VerifyZKP(c *C) {
	cc := ed448.NewScalar([]byte{
		0xc4, 0x6b, 0x41, 0x3c, 0xd2, 0x70, 0xc9, 0xcb,
		0x86, 0x68, 0x18, 0x57, 0x67, 0x63, 0x36, 0xf0,
		0xb1, 0x38, 0x07, 0x35, 0x0b, 0x51, 0x90, 0xf5,
		0x0a, 0x1c, 0x21, 0x27, 0xbd, 0xd0, 0x28, 0x0e,
		0xd3, 0x0d, 0x30, 0x4d, 0x5e, 0x72, 0x32, 0x00,
		0xd2, 0x25, 0x94, 0xc0, 0xa7, 0x08, 0x0d, 0xa0,
		0xd2, 0xb5, 0xa4, 0x6e, 0x32, 0xc9, 0x02, 0x24,
	},
	)

	dd := ed448.NewScalar([]byte{
		0x33, 0xd9, 0x16, 0x6f, 0xc0, 0x51, 0xaf, 0x57,
		0xce, 0x26, 0xad, 0x36, 0x0b, 0x5f, 0x36, 0x31,
		0xde, 0xfd, 0xce, 0x79, 0x3e, 0x8a, 0xbe, 0xce,
		0xde, 0x07, 0xa9, 0x55, 0x42, 0x2f, 0xd7, 0xf1,
		0x2c, 0xf2, 0xcf, 0xb2, 0xa1, 0x8d, 0xcd, 0xff,
		0x2d, 0xda, 0x6b, 0x3f, 0x58, 0xf7, 0xf2, 0x5f,
		0x2d, 0x4a, 0x5b, 0x91, 0xcd, 0x36, 0xfd, 0x1b,
	},
	)

	gen := ed448.NewPointFromBytes()

	ok := verifyZKP(dd, cc, gen, byte(01))

	c.Assert(ok, Equals, false)
}

func (s *OTR4Suite) Test_VerifyZKP2(c *C) {
	d5 := ed448.NewScalar([]byte{
		0xc4, 0x6b, 0x41, 0x3c, 0xd2, 0x70, 0xc9, 0xcb,
		0x86, 0x68, 0x18, 0x57, 0x67, 0x63, 0x36, 0xf0,
		0xb1, 0x38, 0x07, 0x35, 0x0b, 0x51, 0x90, 0xf5,
		0x0a, 0x1c, 0x21, 0x27, 0xbd, 0xd0, 0x28, 0x0e,
		0xd3, 0x0d, 0x30, 0x4d, 0x5e, 0x72, 0x32, 0x00,
		0xd2, 0x25, 0x94, 0xc0, 0xa7, 0x08, 0x0d, 0xa0,
		0xd2, 0xb5, 0xa4, 0x6e, 0x32, 0xc9, 0x02, 0x24,
	},
	)

	d6 := ed448.NewScalar([]byte{
		0x33, 0xd9, 0x16, 0x6f, 0xc0, 0x51, 0xaf, 0x57,
		0xce, 0x26, 0xad, 0x36, 0x0b, 0x5f, 0x36, 0x31,
		0xde, 0xfd, 0xce, 0x79, 0x3e, 0x8a, 0xbe, 0xce,
		0xde, 0x07, 0xa9, 0x55, 0x42, 0x2f, 0xd7, 0xf1,
		0x2c, 0xf2, 0xcf, 0xb2, 0xa1, 0x8d, 0xcd, 0xff,
		0x2d, 0xda, 0x6b, 0x3f, 0x58, 0xf7, 0xf2, 0x5f,
		0x2d, 0x4a, 0x5b, 0x91, 0xcd, 0x36, 0xfd, 0x1b,
	},
	)

	cp := ed448.NewScalar([]byte{
		0x33, 0xd9, 0x16, 0x6f, 0xc0, 0x51, 0xaf, 0x57,
		0xce, 0x26, 0xad, 0x36, 0x0b, 0x5f, 0x36, 0x31,
		0xde, 0xfd, 0xce, 0x79, 0x3e, 0x8a, 0xbe, 0xce,
		0xde, 0x07, 0xa9, 0x55, 0x42, 0x2f, 0xd7, 0xf1,
		0x2c, 0xf2, 0xcf, 0xb2, 0xa1, 0x8d, 0xcd, 0xff,
		0x2d, 0xda, 0x6b, 0x3f, 0x58, 0xf7, 0xf2, 0x5f,
		0x2d, 0x4a, 0x5b, 0x91, 0xcd, 0x36, 0xfd, 0x1b,
	},
	)

	g2 := ed448.NewPointFromBytes([]byte{
		0xc4, 0x6b, 0x41, 0x3c, 0xd2, 0x70, 0xc9, 0xcb,
		0x86, 0x68, 0x18, 0x57, 0x67, 0x63, 0x36, 0xf0,
		0xb1, 0x38, 0x07, 0x35, 0x0b, 0x51, 0x90, 0xf5,
		0x0a, 0x1c, 0x21, 0x27, 0xbd, 0xd0, 0x28, 0x0e,
		0xd3, 0x0d, 0x30, 0x4d, 0x5e, 0x72, 0x32, 0x00,
		0xd2, 0x25, 0x94, 0xc0, 0xa7, 0x08, 0x0d, 0xa0,
		0xd2, 0xb5, 0xa4, 0x6e, 0x32, 0xc9, 0x02, 0x24,
	},
	)

	g3 := ed448.NewPointFromBytes([]byte{
		0x33, 0xd9, 0x16, 0x6f, 0xc0, 0x51, 0xaf, 0x57,
		0xce, 0x26, 0xad, 0x36, 0x0b, 0x5f, 0x36, 0x31,
		0xde, 0xfd, 0xce, 0x79, 0x3e, 0x8a, 0xbe, 0xce,
		0xde, 0x07, 0xa9, 0x55, 0x42, 0x2f, 0xd7, 0xf1,
		0x2c, 0xf2, 0xcf, 0xb2, 0xa1, 0x8d, 0xcd, 0xff,
		0x2d, 0xda, 0x6b, 0x3f, 0x58, 0xf7, 0xf2, 0x5f,
		0x2d, 0x4a, 0x5b, 0x91, 0xcd, 0x36, 0xfd, 0x1b,
	},
	)

	pb := ed448.NewPointFromBytes([]byte{
		0x33, 0xd9, 0x16, 0x6f, 0xc0, 0x51, 0xaf, 0x57,
		0xce, 0x26, 0xad, 0x36, 0x0b, 0x5f, 0x36, 0x31,
		0xde, 0xfd, 0xce, 0x79, 0x3e, 0x8a, 0xbe, 0xce,
		0xde, 0x07, 0xa9, 0x55, 0x42, 0x2f, 0xd7, 0xf1,
		0x2c, 0xf2, 0xcf, 0xb2, 0xa1, 0x8d, 0xcd, 0xff,
		0x2d, 0xda, 0x6b, 0x3f, 0x58, 0xf7, 0xf2, 0x5f,
		0x2d, 0x4a, 0x5b, 0x91, 0xcd, 0x36, 0xfd, 0x1b,
	},
	)

	qb := ed448.NewPointFromBytes([]byte{
		0x33, 0xd9, 0x16, 0x6f, 0xc0, 0x51, 0xaf, 0x57,
		0xce, 0x26, 0xad, 0x36, 0x0b, 0x5f, 0x36, 0x31,
		0xde, 0xfd, 0xce, 0x79, 0x3e, 0x8a, 0xbe, 0xce,
		0xde, 0x07, 0xa9, 0x55, 0x42, 0x2f, 0xd7, 0xf1,
		0x2c, 0xf2, 0xcf, 0xb2, 0xa1, 0x8d, 0xcd, 0xff,
		0x2d, 0xda, 0x6b, 0x3f, 0x58, 0xf7, 0xf2, 0x5f,
		0x2d, 0x4a, 0x5b, 0x91, 0xcd, 0x36, 0xfd, 0x1b,
	},
	)

	ok := verifyZKP2(g2, g3, pb, qb, d5, d6, cp, byte(1))

	c.Assert(ok, Equals, false)
}

func (s *OTR4Suite) Test_VerifyZKP3(c *C) {
	d5 := ed448.NewScalar([]byte{
		0xc4, 0x6b, 0x41, 0x3c, 0xd2, 0x70, 0xc9, 0xcb,
		0x86, 0x68, 0x18, 0x57, 0x67, 0x63, 0x36, 0xf0,
		0xb1, 0x38, 0x07, 0x35, 0x0b, 0x51, 0x90, 0xf5,
		0x0a, 0x1c, 0x21, 0x27, 0xbd, 0xd0, 0x28, 0x0e,
		0xd3, 0x0d, 0x30, 0x4d, 0x5e, 0x72, 0x32, 0x00,
		0xd2, 0x25, 0x94, 0xc0, 0xa7, 0x08, 0x0d, 0xa0,
		0xd2, 0xb5, 0xa4, 0x6e, 0x32, 0xc9, 0x02, 0x24,
	},
	)

	d6 := ed448.NewScalar([]byte{
		0x33, 0xd9, 0x16, 0x6f, 0xc0, 0x51, 0xaf, 0x57,
		0xce, 0x26, 0xad, 0x36, 0x0b, 0x5f, 0x36, 0x31,
		0xde, 0xfd, 0xce, 0x79, 0x3e, 0x8a, 0xbe, 0xce,
		0xde, 0x07, 0xa9, 0x55, 0x42, 0x2f, 0xd7, 0xf1,
		0x2c, 0xf2, 0xcf, 0xb2, 0xa1, 0x8d, 0xcd, 0xff,
		0x2d, 0xda, 0x6b, 0x3f, 0x58, 0xf7, 0xf2, 0x5f,
		0x2d, 0x4a, 0x5b, 0x91, 0xcd, 0x36, 0xfd, 0x1b,
	},
	)

	cp := ed448.NewScalar([]byte{
		0x33, 0xd9, 0x16, 0x6f, 0xc0, 0x51, 0xaf, 0x57,
		0xce, 0x26, 0xad, 0x36, 0x0b, 0x5f, 0x36, 0x31,
		0xde, 0xfd, 0xce, 0x79, 0x3e, 0x8a, 0xbe, 0xce,
		0xde, 0x07, 0xa9, 0x55, 0x42, 0x2f, 0xd7, 0xf1,
		0x2c, 0xf2, 0xcf, 0xb2, 0xa1, 0x8d, 0xcd, 0xff,
		0x2d, 0xda, 0x6b, 0x3f, 0x58, 0xf7, 0xf2, 0x5f,
		0x2d, 0x4a, 0x5b, 0x91, 0xcd, 0x36, 0xfd, 0x1b,
	},
	)

	g2 := ed448.NewPointFromBytes([]byte{
		0xc4, 0x6b, 0x41, 0x3c, 0xd2, 0x70, 0xc9, 0xcb,
		0x86, 0x68, 0x18, 0x57, 0x67, 0x63, 0x36, 0xf0,
		0xb1, 0x38, 0x07, 0x35, 0x0b, 0x51, 0x90, 0xf5,
		0x0a, 0x1c, 0x21, 0x27, 0xbd, 0xd0, 0x28, 0x0e,
		0xd3, 0x0d, 0x30, 0x4d, 0x5e, 0x72, 0x32, 0x00,
		0xd2, 0x25, 0x94, 0xc0, 0xa7, 0x08, 0x0d, 0xa0,
		0xd2, 0xb5, 0xa4, 0x6e, 0x32, 0xc9, 0x02, 0x24,
	},
	)

	g3 := ed448.NewPointFromBytes([]byte{
		0x33, 0xd9, 0x16, 0x6f, 0xc0, 0x51, 0xaf, 0x57,
		0xce, 0x26, 0xad, 0x36, 0x0b, 0x5f, 0x36, 0x31,
		0xde, 0xfd, 0xce, 0x79, 0x3e, 0x8a, 0xbe, 0xce,
		0xde, 0x07, 0xa9, 0x55, 0x42, 0x2f, 0xd7, 0xf1,
		0x2c, 0xf2, 0xcf, 0xb2, 0xa1, 0x8d, 0xcd, 0xff,
		0x2d, 0xda, 0x6b, 0x3f, 0x58, 0xf7, 0xf2, 0x5f,
		0x2d, 0x4a, 0x5b, 0x91, 0xcd, 0x36, 0xfd, 0x1b,
	},
	)

	pb := ed448.NewPointFromBytes([]byte{
		0x33, 0xd9, 0x16, 0x6f, 0xc0, 0x51, 0xaf, 0x57,
		0xce, 0x26, 0xad, 0x36, 0x0b, 0x5f, 0x36, 0x31,
		0xde, 0xfd, 0xce, 0x79, 0x3e, 0x8a, 0xbe, 0xce,
		0xde, 0x07, 0xa9, 0x55, 0x42, 0x2f, 0xd7, 0xf1,
		0x2c, 0xf2, 0xcf, 0xb2, 0xa1, 0x8d, 0xcd, 0xff,
		0x2d, 0xda, 0x6b, 0x3f, 0x58, 0xf7, 0xf2, 0x5f,
		0x2d, 0x4a, 0x5b, 0x91, 0xcd, 0x36, 0xfd, 0x1b,
	},
	)

	qb := ed448.NewPointFromBytes([]byte{
		0x33, 0xd9, 0x16, 0x6f, 0xc0, 0x51, 0xaf, 0x57,
		0xce, 0x26, 0xad, 0x36, 0x0b, 0x5f, 0x36, 0x31,
		0xde, 0xfd, 0xce, 0x79, 0x3e, 0x8a, 0xbe, 0xce,
		0xde, 0x07, 0xa9, 0x55, 0x42, 0x2f, 0xd7, 0xf1,
		0x2c, 0xf2, 0xcf, 0xb2, 0xa1, 0x8d, 0xcd, 0xff,
		0x2d, 0xda, 0x6b, 0x3f, 0x58, 0xf7, 0xf2, 0x5f,
		0x2d, 0x4a, 0x5b, 0x91, 0xcd, 0x36, 0xfd, 0x1b,
	},
	)

	ok := verifyZKP3(g2, g3, pb, qb, d5, d6, cp, byte(1))

	c.Assert(ok, Equals, false)
}

func (s *OTR4Suite) Test_VerifyZKP4(c *C) {
	d7 := ed448.NewScalar([]byte{
		0xc4, 0x6b, 0x41, 0x3c, 0xd2, 0x70, 0xc9, 0xcb,
		0x86, 0x68, 0x18, 0x57, 0x67, 0x63, 0x36, 0xf0,
		0xb1, 0x38, 0x07, 0x35, 0x0b, 0x51, 0x90, 0xf5,
		0x0a, 0x1c, 0x21, 0x27, 0xbd, 0xd0, 0x28, 0x0e,
		0xd3, 0x0d, 0x30, 0x4d, 0x5e, 0x72, 0x32, 0x00,
		0xd2, 0x25, 0x94, 0xc0, 0xa7, 0x08, 0x0d, 0xa0,
		0xd2, 0xb5, 0xa4, 0x6e, 0x32, 0xc9, 0x02, 0x24,
	},
	)

	cr := ed448.NewScalar([]byte{
		0x33, 0xd9, 0x16, 0x6f, 0xc0, 0x51, 0xaf, 0x57,
		0xce, 0x26, 0xad, 0x36, 0x0b, 0x5f, 0x36, 0x31,
		0xde, 0xfd, 0xce, 0x79, 0x3e, 0x8a, 0xbe, 0xce,
		0xde, 0x07, 0xa9, 0x55, 0x42, 0x2f, 0xd7, 0xf1,
		0x2c, 0xf2, 0xcf, 0xb2, 0xa1, 0x8d, 0xcd, 0xff,
		0x2d, 0xda, 0x6b, 0x3f, 0x58, 0xf7, 0xf2, 0x5f,
		0x2d, 0x4a, 0x5b, 0x91, 0xcd, 0x36, 0xfd, 0x1b,
	},
	)

	g3a := ed448.NewPointFromBytes([]byte{
		0xc4, 0x6b, 0x41, 0x3c, 0xd2, 0x70, 0xc9, 0xcb,
		0x86, 0x68, 0x18, 0x57, 0x67, 0x63, 0x36, 0xf0,
		0xb1, 0x38, 0x07, 0x35, 0x0b, 0x51, 0x90, 0xf5,
		0x0a, 0x1c, 0x21, 0x27, 0xbd, 0xd0, 0x28, 0x0e,
		0xd3, 0x0d, 0x30, 0x4d, 0x5e, 0x72, 0x32, 0x00,
		0xd2, 0x25, 0x94, 0xc0, 0xa7, 0x08, 0x0d, 0xa0,
		0xd2, 0xb5, 0xa4, 0x6e, 0x32, 0xc9, 0x02, 0x24,
	},
	)

	qa := ed448.NewPointFromBytes([]byte{
		0x33, 0xd9, 0x16, 0x6f, 0xc0, 0x51, 0xaf, 0x57,
		0xce, 0x26, 0xad, 0x36, 0x0b, 0x5f, 0x36, 0x31,
		0xde, 0xfd, 0xce, 0x79, 0x3e, 0x8a, 0xbe, 0xce,
		0xde, 0x07, 0xa9, 0x55, 0x42, 0x2f, 0xd7, 0xf1,
		0x2c, 0xf2, 0xcf, 0xb2, 0xa1, 0x8d, 0xcd, 0xff,
		0x2d, 0xda, 0x6b, 0x3f, 0x58, 0xf7, 0xf2, 0x5f,
		0x2d, 0x4a, 0x5b, 0x91, 0xcd, 0x36, 0xfd, 0x1b,
	},
	)

	qb := ed448.NewPointFromBytes([]byte{
		0x33, 0xd9, 0x16, 0x6f, 0xc0, 0x51, 0xaf, 0x57,
		0xce, 0x26, 0xad, 0x36, 0x0b, 0x5f, 0x36, 0x31,
		0xde, 0xfd, 0xce, 0x79, 0x3e, 0x8a, 0xbe, 0xce,
		0xde, 0x07, 0xa9, 0x55, 0x42, 0x2f, 0xd7, 0xf1,
		0x2c, 0xf2, 0xcf, 0xb2, 0xa1, 0x8d, 0xcd, 0xff,
		0x2d, 0xda, 0x6b, 0x3f, 0x58, 0xf7, 0xf2, 0x5f,
		0x2d, 0x4a, 0x5b, 0x91, 0xcd, 0x36, 0xfd, 0x1b,
	},
	)

	ra := ed448.NewPointFromBytes([]byte{
		0x33, 0xd9, 0x16, 0x6f, 0xc0, 0x51, 0xaf, 0x57,
		0xce, 0x26, 0xad, 0x36, 0x0b, 0x5f, 0x36, 0x31,
		0xde, 0xfd, 0xce, 0x79, 0x3e, 0x8a, 0xbe, 0xce,
		0xde, 0x07, 0xa9, 0x55, 0x42, 0x2f, 0xd7, 0xf1,
		0x2c, 0xf2, 0xcf, 0xb2, 0xa1, 0x8d, 0xcd, 0xff,
		0x2d, 0xda, 0x6b, 0x3f, 0x58, 0xf7, 0xf2, 0x5f,
		0x2d, 0x4a, 0x5b, 0x91, 0xcd, 0x36, 0xfd, 0x1b,
	},
	)

	ok := verifyZKP4(g3a, qa, qb, ra, d7, cr, byte(1))

	c.Assert(ok, Equals, false)
}
