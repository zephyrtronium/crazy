package crazy

import (
	"encoding/binary"
	"io"
)

const (
	lfgJ = 273
	lfgK = 607
)

// LFG implements a lagged Fibonacci generator. Numbers are produced under the
// recurrence S[n] = f(S[n-j], S[n-k]) (mod m), 0 < j < k. For the sake of
// speed, this implementation defines f(x, y) = x + y, j = 273, k = 607,
// m = 2**64. This yields a period of (2**k - 1) * m/2 = 2**670 - 63.
type LFG struct {
	f, t int
	s    [lfgK]uint64
}

// NewLFG produces an unseeded LFG. Call either lfg.Seed[IV]() or lfg.Restore()
// prior to use.
func NewLFG() *LFG {
	return &LFG{f: lfgK - lfgJ}
}

// Seed calls SeedInt64(lfg, x). This serves to satisfy the rand.Source
// interface.
func (lfg *LFG) Seed(x int64) {
	SeedInt64(lfg, x)
}

// SeedIV initializes the generator using all bits of iv, which may be of any
// size or nil.
func (lfg *LFG) SeedIV(iv []byte) {
	lfg.f, lfg.t = lfgK-lfgJ, 0
	copy(lfg.s[:], lfgS0[:])
	for i, v := range iv {
		lfg.s[(i+lfgK-lfgJ)%lfgK] += lfgS0[(int(v)+256*i)%lfgK]
	}
	for i := 0; i < lfgK*8; i++ {
		lfg.Uint64()
	}
}

// Uint64 produces a 64-bit pseudo-random value. This primarily serves to
// satisfy the rand.Source64 interface, but it also provides direct access to
// the algorithm's values, which can simplify usage in some scenarios.
func (lfg *LFG) Uint64() uint64 {
	x := lfg.s[lfg.f] + lfg.s[lfg.t]
	lfg.f++
	lfg.t++
	if lfg.f >= lfgK {
		lfg.f = 0
	} else if lfg.t >= lfgK {
		lfg.t = 0
	}
	lfg.s[lfg.f] = x
	return x
}

// Read fills p with random bytes generated 64 bits at a time, discarding
// unused bytes. n will always be len(p) and err will always be nil.
func (lfg *LFG) Read(p []byte) (n int, err error) {
	n = len(p)
	for len(p) >= 8 {
		binary.LittleEndian.PutUint64(p, lfg.Uint64())
		p = p[8:]
	}
	b := [8]byte{}
	binary.LittleEndian.PutUint64(b[:], lfg.Uint64())
	copy(p, b[:])
	return n, nil
}

// Int63 generates an integer in the interval [0, 2**63 - 1]. This serves to
// satisfy the rand.Source interface.
func (lfg *LFG) Int63() int64 {
	return int64(lfg.Uint64() >> 1)
}

// Save serializes the current state of the LFG. Values produced by an LFG that
// has Restore()d this state are guaranteed to match those produced by this
// exact generator. n should always be K*8 = 4856 bytes.
func (lfg *LFG) Save(into io.Writer) (n int, err error) {
	p := [lfgK * 8]byte{}
	// We avoid having to save f and t by rotating the LFG's ring buffer such
	// that t is in the first element.
	for i, v := range lfg.s[lfg.t:] {
		binary.LittleEndian.PutUint64(p[i<<3:], v)
	}
	for i, v := range lfg.s[:lfg.t] {
		binary.LittleEndian.PutUint64(p[(lfgK-lfg.t+i)<<3:], v)
	}
	return into.Write(p[:])
}

// Restore loads a Save()d LFG state. This reads K*8 = 4856 bytes as the last K
// values of the LFG.
func (lfg *LFG) Restore(from io.Reader) (n int, err error) {
	p := [lfgK * 8]byte{}
	if n, err = from.Read(p[:]); n < len(p) {
		return n, err
	}
	for i := range lfg.s {
		lfg.s[i] = binary.LittleEndian.Uint64(p[i<<3:])
	}
	lfg.f, lfg.t = lfgK-lfgJ, 0
	return n, nil
}

// Initial seed state. This is the state of the algorithm after 2**34
// iterations initialized with 0, 1, ... 606.
var lfgS0 = [lfgK]uint64{
	0x9fc957b19b39ea28, 0x1d77b33b087c7c87, 0xb8f31fe09b7b9f67, 0xc11b2756053f5ed5,
	0xbfc776182b975380, 0x4cb56bd60674d700, 0x1882170f226de11e, 0x2705fc53fd3b28e7,
	0xf20b931247459d6f, 0x733fb29c7f3229a0, 0x61a8f4206c3a5f73, 0x8e50e2aea8f06622,
	0x635f24c31c569a23, 0xa5b2da002839e24e, 0xb039b035e547147d, 0x15f7ad2204d2495a,
	0xb177eec6dded1818, 0xe7fa8e6c40f53589, 0x47d072d7bec738e7, 0xe2e3774704873b56,
	0xcea4a59486d7ec32, 0x1ae2551c01e9f1f6, 0xdb40d4dda91cde84, 0x295fb1b2624daab5,
	0x803b7ff58462704e, 0x176670a1b1db3a8d, 0x80266c91f2e3c45f, 0x52fe74cedbe1edfc,
	0x3f4c1e19da2a4a92, 0xaed0c8e095947d6, 0xeb2dd04e113963ec, 0x93f202b34490259b,
	0xa5925f794f246b56, 0x840314097e343ed, 0x206770de1c02b2ca, 0x27fdf62316dda9dd,
	0xac10c477a0cefee9, 0x523daa0265c618df, 0x61868e49ca95a36c, 0xacddcc949e93a692,
	0xe110851f629797d9, 0xb10974c87cfc4f5a, 0x643f658a687956d2, 0x2a8030e6c5321407,
	0x57c904fdeb25c23f, 0x6ae864f1f87a390, 0xa870c87ce6fa690, 0x3f2bfd1ea185ae39,
	0x3df4a4542f692ac3, 0xc7009e262b4e5698, 0x7ae5615daa0f45bd, 0xf4e12bca0c0d1159,
	0xbfd8a7885ba3102c, 0x88c919669614e7b7, 0xc688d152cd90600b, 0x210d9977907e7a4e,
	0x5d342a526bb387b, 0xb1ec3ccf4b6cbfb, 0xe5b5e49cf8427990, 0xb4654d7326ae9239,
	0x32c7c74ea64ca72a, 0xe00f2fec37a2a884, 0x940a01a08e1256cd, 0x2b1c2f5ffd4e64b9,
	0x9ae80d9c802052e8, 0xad6aa0ea615955a9, 0xf35316a207941e65, 0xc05a609503702a33,
	0x3843fd6131c11b8d, 0x49d237e0df2f68ad, 0x790dd3313cf2b2fc, 0x96575250bf93dadd,
	0x34a2a0c335883ffe, 0xb1845b648af86c23, 0xc7e5642b95efc503, 0x52130fbbd6e943fc,
	0x52c1f3cfb8d39031, 0x47b183534fca934c, 0xca24ca79ba67934, 0xd6fd920cb3e4f9ab,
	0x72e959925b6eb6d7, 0x161ee79cb63dda2c, 0xc720366c5f0dc0fb, 0xd5ea4a062e9779fb,
	0xf293af0cb192b32c, 0x57ca10d6bd9b6a5f, 0xde6d926cc72dbf5d, 0xdab65a659b3a51da,
	0xb1032a876b98072e, 0x1d64393dad0326a9, 0x28664fa4ffd40baa, 0x4b8e4183a3ab4515,
	0xd8b5e8e3eae1a5cf, 0xae2c33858b690977, 0x5100505724fda789, 0x93c9cb74b8d1ae34,
	0xb7b75304fa37f515, 0x114e4e2a3fcc8e91, 0xb407b29c62f17f97, 0x8d528944c490d5f3,
	0x2d7d809fbfb6a662, 0x8b070d6adef4f72a, 0x21b8e4f03f10cc7d, 0xec32db0eadb4100a,
	0x46bb53d48d88cd9f, 0x31f200e5f25f928, 0x81979422537ec984, 0x135eae7a452d2c32,
	0x3c522841c769078c, 0xfc8db71932b9b1db, 0xa988a0c7cf63692c, 0xe6a260764c632af3,
	0x9ee8db21f607cbeb, 0x89a2889805fc13fc, 0x4fa864a7d23a9b7f, 0x1721720fa9034d99,
	0x516c69827b5d10d4, 0x5de87c1fbdeaaa88, 0xfdb17da692992fc6, 0x4b4e2bfde4386a40,
	0x1dffe044c5dcfa6e, 0xd7059fbf35e3f8b9, 0xa5a437e61317e65f, 0xb19942338d7e887d,
	0xb96dccd0e146381c, 0xff50ab546fbb5e7d, 0x350bd6c46f60b29c, 0x738a67ee376a4985,
	0xe33c910f42d3bd7, 0x3fe664963c1c347f, 0xa5326b005588d415, 0x48b71ddea385882e,
	0x99def931306fd95c, 0xf5e58769a4fe7fcf, 0x4b2260c0dda3860, 0x25b7d3b1f60fc804,
	0xf6e8320383b8b0f0, 0xb984687833b6169b, 0x6317b57af465e4ac, 0x39ba3125ffcfa491,
	0xa16ffdd538184d6a, 0x2ee5b2fadb08e197, 0x36d8eecaad6064b8, 0xb3b33eec85955a1c,
	0x5f7ca001b9df5f6a, 0xf2b3165c916a7287, 0xfdf40ac754dd6da6, 0x699f8888072b68c0,
	0x83adbf5c40ee2d12, 0x2ded4b7ccdbe2f7e, 0x14c87e51aaba8fcc, 0x86064b22c87247a0,
	0xe9e1ef95467d4adf, 0xc4afde761a47fc81, 0x1aadfe580ec58766, 0x7fdb3ebe5ff1cf3b,
	0x9c8f9510138265c5, 0x1e1ee3e1d32d9471, 0x633855bddedf9698, 0x3f877e33facae68a,
	0x44c645f1bf8b442b, 0xd7511ca8e8b48826, 0x51389314623e8d96, 0xe9c67e07eba1a6eb,
	0x8ad9506a0420e39b, 0xb846cf752661d2e6, 0x19cd9dfcd186fdb0, 0x60def2729b04214,
	0x2cec8de7eb095f54, 0x10519406770e202a, 0x9a722505a7341f8b, 0x5a382de1c374df36,
	0x25ce96885c159164, 0x177f6fd70996c10c, 0x93e0663cb767d1ad, 0x6e4e30a2b5acb588,
	0x821b997a8cb0f4f4, 0xf47e0799399d91a9, 0x156f9a4e99622a8b, 0x5accee7fd9c06577,
	0xbff911d12f0d7c, 0xd02189ca594fd1ec, 0xeac8c9b2b81b64d6, 0x4bae086ee642ebe8,
	0x76fbf54608c118a6, 0xf6f04288d020541c, 0x80b2cdd54fb28140, 0x78bd0fd5adf04655,
	0x28e4561ca5974e12, 0x77433374a9a63d40, 0x3cef5744d27c8b24, 0x5d93cf01243f1cdd,
	0x76a91791c291d667, 0xbf9e2a62e01d5e11, 0xe9083036b97d8016, 0x83feeadb36aed21a,
	0xc53e300a010bc099, 0x1bb82a62f50d9141, 0x35b81b1ccff88804, 0xd41dbe1010d7a25d,
	0x36404dfc876e608b, 0x6787ec6fc7043c03, 0x30fca03a0820f19, 0xb4476557dc2b3964,
	0x998204db50c40964, 0xafc9aa433cf6dcbf, 0xd480fe9afd03e341, 0xfe5f19466c14d2f9,
	0x966758033c3606ec, 0x70777f2ff69e9516, 0xd04379f1c235bfad, 0x8c1f3f829ef5b45d,
	0x8cb7a49ca45080f9, 0x4f2b162743062209, 0xf6e223ffd21c0a60, 0xd0a7def8a7bec412,
	0x99dd9a7ec09803fa, 0x1e98f5c9849639b9, 0x6ac421330a9ba1df, 0x238e60ad40f729d9,
	0x9f416ae20e060a2, 0xa130e2e8366a79ea, 0x76ef99997c5965d4, 0xd1207fd2f209a0a,
	0x3d4c881ede57c76e, 0x4ac92ecb0c89bdaf, 0xc73f23ce05053419, 0x2f77a219d0786d5d,
	0x525dea5b43435aaa, 0xc913d58d222baffc, 0x7e53eff3a456d469, 0xc4ee89f29c3a95d1,
	0xe541f55189a13e7b, 0x1a747057f660f620, 0x5ccc34fcc643b847, 0x99b758d355145de3,
	0xe718e817d7c1daa0, 0x850429b5fe002677, 0x2a9a0ef83a4439, 0xc81982bb1dd2b3b7,
	0x37833e76da98ab2f, 0x577d9b216f4c3bda, 0x6077b5cf7effa0ff, 0x39a881ea1593b6bd,
	0x6a17d46ea0e051a2, 0x9dffb34161eb5399, 0xf8669ac6e37276ff, 0x4e3d6cb7e15f173e,
	0x362f5f4fb594ea74, 0x66e6b7c646d46e39, 0xe3ada59287a617a, 0x4c92d50a5c9fad2b,
	0x2406577acf994e5c, 0x3ecbf3b63d7eead6, 0x35417e59d9e72f08, 0x35b1c3bf54414279,
	0x8bb1887d0bb67a10, 0x2d608dd01c2118ce, 0x17a9a78bbb64f5fe, 0x55ceae793e038e45,
	0x80e41e56d4b47d19, 0xb9df69a432aee805, 0xb5863923de1d8893, 0x46567f55496aad98,
	0x4c0009b9de087afc, 0x621f1cd4aa9a359b, 0x2f906164d6f3276f, 0x4b0ea8ed7a8cb1b3,
	0xe01aab3c249417c9, 0xb766876e191ea425, 0xa7d5df71bcc36bba, 0x9aab787fc11dab6,
	0xa4a82b955d554d9f, 0x7dae5b896d42925f, 0x9b7b6ca592ff22e0, 0x828077569c3bf6e,
	0xfeac4ec22657f4ab, 0x8cedf5bddadd8380, 0xcbccab391bf90a1e, 0xe83e544dacd47c9,
	0xcb0596be4a0a7488, 0x81341f8a37ec8c31, 0xee694183ed0835d3, 0x2ca7ee8e3cb606af,
	0xd50e421473663401, 0x4253b53d0be3482b, 0xa86d635bd0d322f, 0x65bdfcec1f8b34dd,
	0x9b8041a4d91acebe, 0x36829fa563081d71, 0x5fd5e46b7dd2035e, 0x9b13046f45c0026f,
	0xebc12e4d8250b0dc, 0x4c3daf877b1205c4, 0xc05e7fc1a732ec8e, 0x4e1edcd4b930cc31,
	0x56dbce432214c599, 0x972af0ac2d78ca3f, 0x68bffbf0410889d2, 0xd2d8083ce8fe299d,
	0xec4da94afe485c96, 0xcba0ee742f2efd44, 0xe040c3c007e01c16, 0xa8c432653356c1af,
	0x11a05cc60a9445bb, 0x62add1c748bed897, 0x18273f9d841f6edd, 0x7968544fadaf713,
	0x8412ce5489f1550c, 0xa62ce58ac4f719f6, 0xf48e44764cf8a8d, 0x4b573e4ad3fe0326,
	0x3432b88ac403f147, 0xcff8efa91a64b781, 0xb335f0c1eb7d0778, 0xc640cb5c5cb8bd35,
	0x2d48d41725f3ae38, 0xaee581513461e151, 0x3d88638aee80300, 0x34a4f096d31607a9,
	0xfec8a7358de37c8a, 0x890bf9d1fbe52bd5, 0xb3e4c3377ec0ef25, 0x79fbca6c61fdcb9c,
	0xcaf77bbe4f95fed3, 0xc8f071de3a71d78b, 0x3dbfb7ec377b7854, 0x5a84c824c2ee1a43,
	0xe4c5a92d963cbe2d, 0x54b8127cdfb9380, 0xda9720d0038bad95, 0xceaf68d62e6c18a9,
	0x7e6279db7f9e14f1, 0xad47689d9156015a, 0xb3fad1b4566fae49, 0x97122dbf6f3c0dec,
	0x6fcbde3c82d1ee2f, 0x1282934de13902c1, 0x45e875b7a63ac8bc, 0xcd0749f2fbdc0bce,
	0x77e99ccc2e50f15a, 0x118e3a7fad6e4d20, 0x2f3b9b505dc34a4f, 0x1d497f1f82a127e1,
	0x9e4b4e7275f46521, 0x7ce1baa155702c25, 0x166108c70af758e0, 0x8a2dab9040f97ef9,
	0xaee413e1ea4c35, 0xf4ef8f8396f7031b, 0xc4f0c9544bdbe5e8, 0xca5b4565183e8077,
	0x9bebc785a789bd2c, 0xa3358e0a5acf2355, 0xb1014ecfa8cfe6cf, 0xeca1399cf89b900,
	0x1ca9650682fb3931, 0x653661ca0c08b733, 0x86a38196099254fe, 0xfc48c7f8d40c927d,
	0xd64cd021d05db554, 0x6c610eb6416b1f7b, 0xb02166752d0e501, 0x2327f1dea3d7396b,
	0x8d27a760473660ba, 0xd5764aa1a08763a8, 0xa2d41cd199949e12, 0x42c97b1d93d406ab,
	0x23ed8790416646e1, 0x5996fb254594997c, 0xa2b964722324f106, 0xd94ad6a8619f565c,
	0xa02af75afb25d06f, 0x5d898ccb1f3e50c8, 0x96b1d785601bd553, 0xca79f61e6ea3438d,
	0x5a8878c5dfd4bd95, 0xbc63cc39d19d2b89, 0x7e787413f458d05c, 0x91c71a57f1ae62ae,
	0x28f379c7823bdb5a, 0xc03b8ed76b50aa4f, 0xacfae9ae9ca9b751, 0x3d19bfae7cffc1c7,
	0xb8467aaba9a4a0f8, 0xeab9ad760ff44811, 0x4921f26fbcda0252, 0xfc56dbdb0884b3d1,
	0x98b142b19e1aad2e, 0xd651dc84cba5407, 0xd4a4a7b15d1c683, 0xfd2bf5454a00124e,
	0xb2bccf16d0bf330e, 0xcca7d7463f4ff896, 0x2884c50b3b0dcf01, 0x1eb2d85af375e0cf,
	0x5642bef413bf1368, 0xd3f99a3da44776ff, 0x163d379b0c2fc159, 0xebc65a7a8e2e66a1,
	0x450fe933db4dae4d, 0x89ad09ef1f71dfb0, 0x1f43ddfad9621d5f, 0x90c9bf9550d8392f,
	0x1f32452dd542782a, 0x3a6fb5709d1da3c, 0xfbbb3f4bb77a8a48, 0x86d78cb5e8662759,
	0xeb5763055dda1eb5, 0x17def649f2478c4e, 0x98cd7bbd3493488d, 0x8eca0042609feef6,
	0x15407d56f2083986, 0x30900e1fa519ef34, 0x7f32077cf977b099, 0xafe0b49d27bcf3ca,
	0xdc9f8824f584ed72, 0x292863c65d682299, 0x3782f7f99b51be5a, 0x7254861470b4938c,
	0x3df0a657ec184291, 0x9948eb8957b499ce, 0x85a2cc128874eb9d, 0xcb6052271a45e9e1,
	0xa0086856d19bbc17, 0x658c9a12fd5fc510, 0x224637b2da27acf8, 0x71c1ce4dd404d17f,
	0x72071fafaf694eaf, 0x3b71214bab27d221, 0x673bde9ac9fd9f38, 0x281642a74ae3397c,
	0x5bbf9405e60b2c70, 0x1da8fdc0f161d229, 0xeb17308ae5ce0e19, 0x822d037b1a6d9396,
	0x45ea377a812076c, 0x5bb31fd90aa85123, 0xeb4e4f2cabf278f4, 0x304dd79647995421,
	0x7bf24dd40ff2f73, 0x89ba9420300e097f, 0x2f824915500d9aa8, 0xf6745bbc9100ae07,
	0xfdd328e496764507, 0x8525d0e68894708b, 0x1667148d0f0bcc33, 0xe5e91bbb5549f6e8,
	0xf7256641d12b89d5, 0x46cdcae9a3662554, 0x5def3cf94c6972ed, 0xaf5ba67bc7c683c1,
	0xd44229b85b12e49, 0xaf5a042184a5e0f, 0x5643cc3ffc82c84f, 0x7443ac84c25fc2bd,
	0x4b494c43f843bb76, 0xf0ed842a0b5ba1d5, 0xa286c65d98da2a52, 0x5bf4932e7a20626e,
	0x5b453e82e9dbc0eb, 0x90511547593c7387, 0x3db7d3590a6bd0c, 0x120f464684d3f8e3,
	0x51f5aadcc0f02d62, 0xf72815dd16790177, 0x3fdf33bbb9fe89a5, 0xd9be2cecea6e6301,
	0xcfa3b4568f6ce2d0, 0xd455da629d471b30, 0xfa0dae149356e334, 0xf0f5e018170f9424,
	0xaa7a48904ac5aabf, 0xd91fe0b3f2b8f6b, 0x474c2f313efb33fc, 0xe8bc2d0677138166,
	0x17a1e001521c62fd, 0x4e7acecbff7cc7b5, 0x22e0db8851221d1, 0x61aaadba3ef1813b,
	0x545dc416d05bf3c2, 0x5251cede25396168, 0xbbf157662c64ca28, 0x3f9f16c26d52f73a,
	0x6d8c623f3b1126b8, 0x8254e090e5cbb684, 0x85b2bb3ae3dfe24, 0xf23d1b48f4bb4903,
	0xb6ecf9bf0f034584, 0xd19af8171dc8ccea, 0x517636d57dba9c25, 0xee05cbe5913d01ea,
	0xc24afc7646a965b, 0x6f5d0585434a2cf3, 0xaee483b93e15137d, 0xf3aac9aafda057a8,
	0xcafbe653e654dfce, 0x1c34796848936d64, 0x5faf7703435144f, 0x90d447da3855f7ea,
	0x491b174f5eb7cad0, 0x62e8b54c303ec880, 0x68f6a47359ef0a94, 0x95e3325b44f869e8,
	0xa634c661bc068a12, 0x40a6eb67633aa99d, 0x9adf194926af88d3, 0xc0adafd182c51a37,
	0xd82d1fa88c5bdb43, 0x6c0d85e543c3acf0, 0xda5bb687f9706278, 0x5c7750028621576c,
	0x50f5579bbfbee915, 0x6664f1ea592113a0, 0xc131e06a32e17917, 0xc1f1d97c04108693,
	0x921363465d60587f, 0x7cdc2cf9157bbd55, 0xc88a3567fbbea93d, 0x3f862aae047fc1e3,
	0x36766d36d4a015ff, 0xb7293b0c2452973f, 0x2fe64ae1d242dd94, 0x58caa0fe77da2ba6,
	0xd00dd473218068e6, 0xcfd2bb7f3fcf40a, 0x6a90fab9183c10e7, 0xe13a124adacde74e,
	0xa0d83cadbaeb455f, 0x89e06ce47468c575, 0xddf57bfab17978f, 0xd31d87c9ac235828,
	0xeed5b22ca130e969, 0x248dcd497129716d, 0xf8ab8b59820113ca, 0x2eebd956096f7455,
	0x9673c5c5d073b058, 0x99838fc970f5bf71, 0x4dcaf5214d20f8d5, 0xe74cf9fc9de50239,
	0x9716a43fdadbdef8, 0x6b97a2dad7dfc239, 0x69f6bc2143f49532, 0x5e1424802a9c1e,
	0x70d5935476c93134, 0x41190d4638fef0e1, 0xcd384cc8d7f4a53e, 0x59eff1657c452637,
	0xa91b078cbf4b4840, 0x9ffd2b8c916752a0, 0x70a50a85392616b2, 0xa82a503f9be1aac,
	0x291b9acd7e545465, 0x93dfbc0088eff644, 0xb76e1cadc9e7201d, 0xc162335beac780bf,
	0x629316442131faa9, 0xd982afdd9d8b607d, 0xe694b7daccabfa87, 0x23e13ff9ab03c1f5,
	0x6eaa6ec4b78d7fa4, 0x35e99292bc92b3bd, 0x656134ac8d0b211a, 0xb7bf1f07d04e7bc4,
	0x80d2f494f27a2bc0, 0xff26e48896d10029, 0xc4156e7b330b95fa, 0xa95763ccbcacd475,
	0xc3cbd424b30dca95, 0x20980921795182dc, 0xba4f61f4ce65e0bf, 0xa1684a0ca627bb5f,
	0x266c73c2a427e1d6, 0x26970dd19c62260f, 0xeeb0908cba34d9c6, 0x2633cf0394cd84f5,
	0x7db16a250419c0cf, 0xde291ff4831961ce, 0x17d1a1de98ad188b, 0x81e9764d398d6a2d,
	0x1fe9298e9b78bdc6, 0x184fc4557eeb34c5, 0x668d310d604a4c03, 0x9cbc905d15df3677,
	0x3a348235cb3a4b0, 0x11de227c852e062a, 0x5e70f786e1cdb355, 0x82774f01b16701b1,
	0xc14342b7eee5ec87, 0xf684c111c8cd1b8f, 0x2c3684d11d0e5e08, 0xb7e80d4e28c4d818,
	0xe5489b1e44e5f0e6, 0xfcf242aa004ae6e4, 0x52c0f1233e4e7529, 0xd3a50f7a1302f242,
	0x8d84791e45b1da47, 0x430ab24223cf62da, 0x896131976d3a1072, 0xd5613b514b428b6e,
	0x37805825f5dcc109, 0x6710b98acccfe878, 0xb21f6278475c9a2b, 0x923a0db46bf0b1f4,
	0x49a09522850f5619, 0xf176749441d2c1d3, 0xfb212c1c3de49c89,
}
