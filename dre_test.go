package otr4

import (
	"crypto/rand"
	"errors"
	"io"
	"testing"

	"github.com/twstrike/ed448"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type DualReceiverEncryptionSuite struct{}

var _ = Suite(&DualReceiverEncryptionSuite{})

type fixedRandReader struct {
	data []byte
	at   int
}

func fixedRand(data []byte) io.Reader {
	return &fixedRandReader{data, 0}
}

func (r *fixedRandReader) Read(p []byte) (n int, err error) {
	if r.at < len(r.data) {
		n = copy(p, r.data[r.at:])
		r.at += fieldBytes
		return
	}
	return 0, io.EOF
}

var randData = []byte{
	0xea, 0x25, 0xbc, 0x1d, 0x8d, 0x18, 0x2f, 0xe2,
	0x33, 0xe4, 0xd1, 0x8c, 0x58, 0xac, 0x3a, 0x75,
	0x32, 0x7a, 0xb0, 0x91, 0xcf, 0x85, 0x81, 0xf8,
	0x2c, 0xc5, 0xf3, 0x55, 0x3d, 0x32, 0x2d, 0x3e,
	0x8c, 0x0b, 0xf5, 0xfb, 0x6a, 0x11, 0xf9, 0x5b,
	0x35, 0xcc, 0x4b, 0xda, 0x05, 0x8a, 0x3c, 0x66,
	0x17, 0x65, 0x14, 0xf4, 0x04, 0x40, 0xd3, 0x02,

	0xad, 0x6a, 0xc8, 0x1c, 0xb1, 0x91, 0xef, 0xcb,
	0xfa, 0x5d, 0x81, 0xbe, 0xf3, 0xa2, 0x01, 0x0f,
	0xb2, 0x32, 0x97, 0x2c, 0x2c, 0x46, 0xb6, 0xd2,
	0x55, 0x19, 0xad, 0xb4, 0x57, 0x74, 0xe6, 0x61,
	0xb0, 0xe2, 0x1a, 0x22, 0x91, 0x82, 0xe4, 0x9a,
	0xf7, 0xff, 0x82, 0x5b, 0x4f, 0xeb, 0x05, 0x2b,
	0x0d, 0xcf, 0x78, 0xe0, 0x02, 0x53, 0x6d, 0x33,

	0xd7, 0xfb, 0x29, 0xa2, 0xd1, 0x3d, 0x98, 0xec,
	0x18, 0xe8, 0x91, 0x07, 0x22, 0x13, 0xbf, 0x76,
	0x5a, 0x2f, 0x93, 0xa6, 0xd4, 0xcf, 0xe7, 0x77,
	0x99, 0x73, 0xbb, 0x20, 0x4e, 0x09, 0xb1, 0xc0,
	0x68, 0x93, 0xf1, 0xf2, 0x35, 0x28, 0x5b, 0xa3,
	0x89, 0x9d, 0x76, 0x75, 0x46, 0x5b, 0xe9, 0xa1,
	0xff, 0x1c, 0x2e, 0x8e, 0xbf, 0x5e, 0x22, 0x15,

	0x2c, 0x80, 0x3f, 0x45, 0x89, 0x58, 0x2d, 0x7e,
	0xda, 0x26, 0x01, 0xa7, 0x78, 0xa6, 0xa6, 0x61,
	0xa3, 0x56, 0x5c, 0x0b, 0xe0, 0xae, 0x09, 0xea,
	0x56, 0xce, 0x09, 0x17, 0xc0, 0x98, 0xd5, 0x40,
	0xe8, 0x9b, 0x18, 0xdf, 0xb7, 0x78, 0x44, 0xc3,
	0xc9, 0x66, 0x62, 0x85, 0xae, 0xa6, 0xd7, 0x10,
	0x7b, 0x14, 0xf0, 0x40, 0xd7, 0x72, 0x92, 0x2b,

	0x02, 0x35, 0xbd, 0x52, 0x34, 0x56, 0x9b, 0x58,
	0x2c, 0x39, 0xaf, 0x2e, 0x92, 0xca, 0x6c, 0x0a,
	0x81, 0x22, 0x88, 0x38, 0xd3, 0xdd, 0x17, 0x25,
	0x27, 0xc9, 0x2d, 0xf6, 0x4d, 0xa1, 0xf2, 0x9c,
	0xbd, 0x08, 0xf4, 0xa0, 0x91, 0x08, 0x79, 0xf6,
	0x8a, 0x78, 0x3c, 0xf0, 0xac, 0x2d, 0x97, 0x03,
	0x54, 0xe3, 0xc6, 0x22, 0xb9, 0xf4, 0x55, 0x3a,
}

var (
	testByteSlice = []byte{
		0xad, 0xd0, 0x35, 0x07, 0x1d, 0x09, 0x6c, 0xd3,
		0xdd, 0xf8, 0x96, 0x59, 0x39, 0x1c, 0x29, 0xa2,
		0x1e, 0x49, 0x34, 0xae, 0xc1, 0x79, 0x0e, 0x85,
		0x1c, 0x06, 0x73, 0xf2, 0xdd, 0x5d, 0x39, 0x71,
		0xf5, 0x70, 0x71, 0x4d, 0x5c, 0xca, 0x18, 0x02,
		0xaf, 0xa3, 0x85, 0x1b, 0x8a, 0x53, 0x39, 0xb7,
		0xa2, 0x33, 0x1b, 0x8a, 0x53, 0x39, 0xb7, 0xa2,
		0x33, 0x2a, 0xf4, 0xf7, 0xb6, 0x26, 0x37, 0x3e,
		0xb7, 0xd5, 0x9a, 0x1b, 0x3c, 0xf2, 0xfd, 0x63,
	}

	testScalar = ed448.NewDecafScalar([]byte{
		0xea, 0x25, 0xbc, 0x1d, 0x8d, 0x18, 0x2f, 0xe2,
		0x33, 0xe4, 0xd1, 0x8c, 0x58, 0xac, 0x3a, 0x75,
		0x32, 0x7a, 0xb0, 0x91, 0xcf, 0x85, 0x81, 0xf8,
		0x2c, 0xc5, 0xf3, 0x55, 0x3d, 0x32, 0x2d, 0x3e,
		0x8c, 0x0b, 0xf5, 0xfb, 0x6a, 0x11, 0xf9, 0x5b,
		0x35, 0xcc, 0x4b, 0xda, 0x05, 0x8a, 0x3c, 0x66,
		0x17, 0x65, 0x14, 0xf4, 0x04, 0x40, 0xd3, 0x02,
	})

	testPubA = ed448.NewPoint(
		[16]uint32{
			0xc08ab5a, 0x4acfddb, 0x5e4c04e, 0x3acf9a3,
			0xd087d8e, 0x647ecbf, 0x194d04f, 0x6bb47af,
			0xc2d95e1, 0x476921c, 0x80f4539, 0xd214e2c,
			0x5f67871, 0x1d6d92b, 0x62e91f4, 0x4524530,
		},
		[16]uint32{
			0xad39e04, 0xeceeefe, 0xacb980d, 0x7b3f1b3,
			0xbfa096c, 0x5b4ee92, 0x5ffb07c, 0xb292193,
			0x32cc0b8, 0x06a3f72, 0xec1eca1, 0x76d0a95,
			0x332251a, 0xa5e5b8d, 0x5b9464b, 0x112b7f5,
		},
		[16]uint32{
			0xd1a2892, 0xaf2db51, 0x7cdb68d, 0x058577b,
			0x0c536eb, 0x7886a30, 0xcca0a19, 0xc93688a,
			0x7e4e010, 0x0736d82, 0xddebbf1, 0x6c3b266,
			0xbd4ccfd, 0xa05d988, 0x45e999e, 0x7870395,
		},
		[16]uint32{
			0x8b3997c, 0x8e9e14e, 0xfc86491, 0x1065ed6,
			0x62efb08, 0x9a65e9c, 0x809bf21, 0x58dd268,
			0x03823b6, 0xba0056d, 0x7434725, 0xd78d8fb,
			0xe1ec6de, 0x2c117ab, 0xf162f56, 0x2191794,
		})

	testPubB = ed448.NewPoint(
		[16]uint32{
			0x7e96333, 0x3080266, 0x2352576, 0x123149b,
			0x87ced35, 0x618b702, 0x77f3cbc, 0x878027d,
			0x646923e, 0x93fc854, 0xa2fc6bf, 0x18c235e,
			0x2c88ef0, 0x33b7029, 0x6c3c6ad, 0x72c71ad,
		},
		[16]uint32{
			0xe1bf670, 0x67a513a, 0xcbb873d, 0x6c0c93c,
			0xc80bc3c, 0x3448b80, 0x7f97f93, 0x25af4bb,
			0xc357d13, 0xd941bfe, 0x333efc9, 0x69a3d5e,
			0x0799c74, 0x966eb5e, 0x7492557, 0x674e842,
		},
		[16]uint32{
			0x418e567, 0x363aeec, 0xfec2818, 0x09f689c,
			0x5ca8f0d, 0xa9d6d7f, 0xa012988, 0x353eb85,
			0x63d46bf, 0xbcb8016, 0x23aaeec, 0x0fcc526,
			0x91ef39e, 0x7f94504, 0x8d785d7, 0xd961794,
		},
		[16]uint32{
			0x9b58ad7, 0x5c1d549, 0xab23cf8, 0x51fbb7b,
			0x0e22d61, 0x91e74c8, 0x8925914, 0xf83f519,
			0x57511c8, 0x2ddeb8b, 0x7764f1e, 0x19e0110,
			0xd446741, 0xa1bfa3d, 0x7876150, 0x93279d1,
		})

	testPubC = ed448.NewPoint(
		[16]uint32{
			0xbaaead5, 0x685c976, 0xb0ca061, 0xba86cd7,
			0xea519fa, 0x57e4ab4, 0xc02f5fc, 0xb82204c,
			0xd78542c, 0x74a01e5, 0x9328d21, 0x3871d1a,
			0x5dcaa1c, 0x0156612, 0x7ea2255, 0xd9d787d,
		},
		[16]uint32{
			0x1a1f403, 0xbfea69f, 0x4a6d633, 0xa60a88d,
			0x4e36fdc, 0x8028db6, 0xc2fe5ba, 0xc58f5c6,
			0x85375d6, 0xbef9d4f, 0x2953a6a, 0xa779d7a,
			0xb729468, 0x9b47792, 0x0ac10fe, 0x434eff1,
		},
		[16]uint32{
			0x803ebe1, 0x18470ae, 0xb80caad, 0x3b777f9,
			0x67a6139, 0x6d9aa39, 0x787cf2e, 0x7c11d72,
			0x0374caf, 0x02cd168, 0x3d6858b, 0xdac0675,
			0xcaae54b, 0xffe21b4, 0x4bb4af8, 0x55a2cce,
		},
		[16]uint32{
			0x384109e, 0x40695a8, 0xa822b8e, 0x6026944,
			0x8e9ae46, 0xaad36b0, 0xa10ec79, 0x505a46c,
			0x4e7c598, 0x0b9daf8, 0xe22fb37, 0xa2aeb13,
			0x126a250, 0x874aa07, 0x0fe3b33, 0x1be9c1d,
		})

	testSigma = []byte{
		0x31, 0x02, 0x64, 0x13, 0xe8, 0xd2, 0xcf, 0x73,
		0x04, 0x1b, 0x8e, 0x53, 0x76, 0x28, 0xf8, 0xad,
		0x1d, 0xc2, 0x1e, 0x28, 0xf7, 0xbd, 0x21, 0x2a,
		0x94, 0x08, 0xcb, 0x0e, 0x25, 0x49, 0x27, 0x8f,
		0x4a, 0x6d, 0x60, 0xcd, 0x04, 0x0b, 0x5c, 0xe0,
		0xd8, 0x17, 0x26, 0x47, 0x9d, 0x72, 0xf5, 0x9f,
		0xe2, 0x0c, 0x90, 0xa9, 0x5e, 0x78, 0xc4, 0x31,
		0xdf, 0xe9, 0x45, 0x2f, 0xd9, 0x2f, 0x02, 0x73,
		0x70, 0xe4, 0x0a, 0xbf, 0xa0, 0x3a, 0x82, 0xd4,
		0x08, 0x8a, 0xe4, 0x01, 0xc9, 0x53, 0xd2, 0x82,
		0x62, 0x54, 0x1e, 0xa3, 0x35, 0x96, 0x8e, 0xff,
		0xb0, 0x0f, 0x70, 0x9a, 0xc7, 0x9c, 0x99, 0xbd,
		0x3b, 0x27, 0x7a, 0xa5, 0x9a, 0x18, 0x16, 0x5e,
		0x12, 0xe4, 0xf7, 0xf6, 0x2f, 0x17, 0xa4, 0x02,
		0xad, 0x6a, 0xc8, 0x1c, 0xb1, 0x91, 0xef, 0xcb,
		0xfa, 0x5d, 0x81, 0xbe, 0xf3, 0xa2, 0x01, 0x0f,
		0xb2, 0x32, 0x97, 0x2c, 0x2c, 0x46, 0xb6, 0xd2,
		0x55, 0x19, 0xad, 0xb4, 0x57, 0x74, 0xe6, 0x61,
		0xb0, 0xe2, 0x1a, 0x22, 0x91, 0x82, 0xe4, 0x9a,
		0xf7, 0xff, 0x82, 0x5b, 0x4f, 0xeb, 0x05, 0x2b,
		0x0d, 0xcf, 0x78, 0xe0, 0x02, 0x53, 0x6d, 0x33,
		0x2c, 0x80, 0x3f, 0x45, 0x89, 0x58, 0x2d, 0x7e,
		0xda, 0x26, 0x01, 0xa7, 0x78, 0xa6, 0xa6, 0x61,
		0xa3, 0x56, 0x5c, 0x0b, 0xe0, 0xae, 0x09, 0xea,
		0x56, 0xce, 0x09, 0x17, 0xc0, 0x98, 0xd5, 0x40,
		0xe8, 0x9b, 0x18, 0xdf, 0xb7, 0x78, 0x44, 0xc3,
		0xc9, 0x66, 0x62, 0x85, 0xae, 0xa6, 0xd7, 0x10,
		0x7b, 0x14, 0xf0, 0x40, 0xd7, 0x72, 0x92, 0x2b,
		0xd7, 0xfb, 0x29, 0xa2, 0xd1, 0x3d, 0x98, 0xec,
		0x18, 0xe8, 0x91, 0x07, 0x22, 0x13, 0xbf, 0x76,
		0x5a, 0x2f, 0x93, 0xa6, 0xd4, 0xcf, 0xe7, 0x77,
		0x99, 0x73, 0xbb, 0x20, 0x4e, 0x09, 0xb1, 0xc0,
		0x68, 0x93, 0xf1, 0xf2, 0x35, 0x28, 0x5b, 0xa3,
		0x89, 0x9d, 0x76, 0x75, 0x46, 0x5b, 0xe9, 0xa1,
		0xff, 0x1c, 0x2e, 0x8e, 0xbf, 0x5e, 0x22, 0x15,
		0x02, 0x35, 0xbd, 0x52, 0x34, 0x56, 0x9b, 0x58,
		0x2c, 0x39, 0xaf, 0x2e, 0x92, 0xca, 0x6c, 0x0a,
		0x81, 0x22, 0x88, 0x38, 0xd3, 0xdd, 0x17, 0x25,
		0x27, 0xc9, 0x2d, 0xf6, 0x4d, 0xa1, 0xf2, 0x9c,
		0xbd, 0x08, 0xf4, 0xa0, 0x91, 0x08, 0x79, 0xf6,
		0x8a, 0x78, 0x3c, 0xf0, 0xac, 0x2d, 0x97, 0x03,
		0x54, 0xe3, 0xc6, 0x22, 0xb9, 0xf4, 0x55, 0x3a,
	}
)

func (s *DualReceiverEncryptionSuite) Test_HashToScalar(c *C) {
	scalar := hashToScalar(testByteSlice)

	exp := ed448.NewDecafScalar([]byte{
		0x1e, 0xda, 0x47, 0xce, 0x5a, 0x2a, 0x10, 0xdb,
		0x67, 0x8a, 0x38, 0x2c, 0xe2, 0x70, 0x2f, 0xea,
		0x92, 0x8d, 0x6a, 0x4c, 0x11, 0x27, 0xfd, 0x7c,
		0xb0, 0x6f, 0x1a, 0x0b, 0x71, 0x82, 0x6b, 0x90,
		0xe3, 0x6b, 0xdd, 0x7d, 0x17, 0xab, 0xfd, 0x9e,
		0xad, 0xf2, 0x04, 0x0d, 0x97, 0x19, 0x46, 0x09,
		0x3d, 0xb3, 0xa3, 0x67, 0xca, 0x01, 0x8d, 0x95,
	})

	c.Assert(scalar, DeepEquals, exp)
}

func (s *DualReceiverEncryptionSuite) Test_Concat(c *C) {
	empty := []byte{}
	bytes := []byte{0x04, 0x2a, 0xf3, 0xcc, 0x69, 0xbb, 0xa1, 0x50}

	exp := []byte{
		0x04, 0x2a, 0xf3, 0xcc, 0x69, 0xbb, 0xa1, 0x50,
		0xea, 0x25, 0xbc, 0x1d, 0x8d, 0x18, 0x2f, 0xe2,
		0x33, 0xe4, 0xd1, 0x8c, 0x58, 0xac, 0x3a, 0x75,
		0x32, 0x7a, 0xb0, 0x91, 0xcf, 0x85, 0x81, 0xf8,
		0x2c, 0xc5, 0xf3, 0x55, 0x3d, 0x32, 0x2d, 0x3e,
		0x8c, 0x0b, 0xf5, 0xfb, 0x6a, 0x11, 0xf9, 0x5b,
		0x35, 0xcc, 0x4b, 0xda, 0x05, 0x8a, 0x3c, 0x66,
		0x17, 0x65, 0x14, 0xf4, 0x04, 0x40, 0xd3, 0x02,
		0x06, 0xea, 0x48, 0xc4, 0x23, 0x28, 0xe1, 0x99,
		0x08, 0xa5, 0x88, 0x8f, 0xad, 0x7f, 0x39, 0xdf,
		0x56, 0xa3, 0xaa, 0x4d, 0x59, 0x66, 0xec, 0xd5,
		0x6c, 0x38, 0x02, 0x8c, 0x80, 0x96, 0xd2, 0xd4,
		0x54, 0x24, 0x76, 0x70, 0xda, 0x99, 0xc5, 0xd6,
		0x81, 0x40, 0x49, 0xcd, 0x76, 0xb1, 0x05, 0xc4,
		0xa8, 0x42, 0x17, 0x09, 0x51, 0xc2, 0xa9, 0x2e,
	}

	c.Assert(func() { concat() }, Panics, "missing concat arguments")
	c.Assert(func() { concat(bytes) }, Panics, "missing concat arguments")
	c.Assert(func() { concat("not a valid input", bytes) }, Panics, "not a valid input")
	c.Assert(concat(empty, bytes, testScalar, testPubA), DeepEquals, exp)
}

func (s *DualReceiverEncryptionSuite) Test_Auth(c *C) {

	secA := ed448.NewDecafScalar([]byte{
		0x71, 0x7b, 0x24, 0xd5, 0xd4, 0x98, 0x0c, 0xfe,
		0xce, 0x60, 0xe7, 0x97, 0x84, 0xf4, 0x1c, 0x72,
		0x01, 0x07, 0xb8, 0x24, 0xa8, 0x43, 0x0e, 0x81,
		0x25, 0xca, 0xb4, 0xa0, 0xda, 0xf5, 0xfa, 0xf6,
		0x0c, 0x90, 0x99, 0x7f, 0x1e, 0xed, 0x83, 0xde,
		0xbe, 0xe7, 0xef, 0x8e, 0xea, 0xeb, 0xc8, 0x5d,
		0x67, 0x5b, 0x3b, 0x04, 0x55, 0x0a, 0x36, 0x2f,
	})
	message := []byte("our message")

	sigma, err := auth(fixedRand(randData), testPubA, testPubB, testPubC, secA, message)
	c.Assert(sigma, DeepEquals, testSigma)
	c.Assert(err, IsNil)

	r := []byte{0x71, 0x7b, 0x24, 0xd5, 0xd4, 0x98, 0x0c, 0xfe}

	out, err := auth(fixedRand(r), testPubA, testPubB, testPubC, secA, message)
	c.Assert(err, DeepEquals, errors.New("unexpected EOF: not enough bytes"))
	c.Assert(out, IsNil)
}

func (s *DualReceiverEncryptionSuite) Test_Verify(c *C) {
	message := []byte("our message")

	expC := true
	b := verify(testPubA, testPubB, testPubC, testSigma, message)

	c.Assert(b, Equals, expC)
}

// functional test - move?
func (s *DualReceiverEncryptionSuite) Test_VerifyAndAuth(c *C) {
	secA := ed448.NewDecafScalar([]byte{
		0x71, 0x7b, 0x24, 0xd5, 0xd4, 0x98, 0x0c, 0xfe,
		0xce, 0x60, 0xe7, 0x97, 0x84, 0xf4, 0x1c, 0x72,
		0x01, 0x07, 0xb8, 0x24, 0xa8, 0x43, 0x0e, 0x81,
		0x25, 0xca, 0xb4, 0xa0, 0xda, 0xf5, 0xfa, 0xf6,
		0x0c, 0x90, 0x99, 0x7f, 0x1e, 0xed, 0x83, 0xde,
		0xbe, 0xe7, 0xef, 0x8e, 0xea, 0xeb, 0xc8, 0x5d,
		0x67, 0x5b, 0x3b, 0x04, 0x55, 0x0a, 0x36, 0x2f,
	})
	message := []byte("hello, I am a message")
	sigma, _ := auth(rand.Reader, testPubA, testPubB, testPubC, secA, message)
	ver := verify(testPubA, testPubB, testPubC, sigma, message)
	c.Assert(ver, Equals, true)

	fakeMessage := []byte("fake message")
	sigma, _ = auth(rand.Reader, testPubA, testPubB, testPubC, secA, message)
	ver = verify(testPubA, testPubB, testPubC, sigma, fakeMessage)
	c.Assert(ver, Equals, false)

	sigma, _ = auth(rand.Reader, testPubA, testPubB, testPubC, secA, message)
	ver = verify(testPubB, testPubB, testPubC, sigma, message)
	c.Assert(ver, Equals, false)

	sigma, _ = auth(rand.Reader, testPubA, testPubB, testPubC, secA, message)
	ver = verify(testPubA, testPubA, testPubC, sigma, message)
	c.Assert(ver, Equals, false)

	sigma, _ = auth(rand.Reader, testPubA, testPubB, testPubC, secA, message)
	ver = verify(testPubA, testPubB, testPubB, sigma, message)
	c.Assert(ver, Equals, false)

	sigma, _ = auth(rand.Reader, testPubA, testPubB, testPubC, secA, message)
	ver = verify(testPubA, testPubB, testPubC, testSigma, message)
	c.Assert(ver, Equals, false)
}
