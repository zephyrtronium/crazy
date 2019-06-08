package crazy

import "math"

// Normal adapts a Source to produce random numbers under a normal
// distribution.
type Normal struct {
	Source
	z            *Ziggurat
	Mean, StdDev float64
}

// NewNormal creates a normal distribution drawing from the specified source
// with given mean and standard deviation.
func NewNormal(src Source, mean, stddev float64) Normal {
	return Normal{
		Source: src,
		Mean:   mean,
		StdDev: stddev,
		z:      normalZig,
	}
}

// Next generates a normal variate.
func (n Normal) Next() float64 {
	x := n.z.GenNext(n.Source)
	return n.Mean + x*n.StdDev
}

func normalPDF(x float64) float64 {
	return math.Exp(-0.5 * x * x)
}

func normalTail(src Source) float64 {
	dist := Uniform0_1{src}
	for {
		x := -math.Log(dist.Next()) / normalR
		y := -math.Log(dist.Next())
		if y+y >= x*x {
			return x + normalR
		}
	}
}

var normalZig = &Ziggurat{
	PDF:      normalPDF,
	Tail:     normalTail,
	Mirrored: true,
	K:        normalK,
	W:        normalW,
	F:        normalF,
}

const normalR = 3.44261985589665

// The tables produced by ziggurat.py are different from the ones in
// math/rand/normal.go, but testing 15000 variates from three separate runs
// together produced the correct mean, standard deviation, and actually a
// slightly better shape than math/rand.

var normalK = [128]uint32{
	0x76ad2212, 0x0, 0x600f1b53, 0x6ce447a6,
	0x725b46a1, 0x7560051d, 0x774921eb, 0x789a25bd,
	0x799045c3, 0x7a4bce5d, 0x7adf629f, 0x7b5682a6,
	0x7bb8a8c6, 0x7c0ae722, 0x7c50cce7, 0x7c8cec5b,
	0x7cc12cd6, 0x7ceefed2, 0x7d177e0b, 0x7d3b8883,
	0x7d5bce6c, 0x7d78dd64, 0x7d932886, 0x7dab0e57,
	0x7dc0dd30, 0x7dd4d688, 0x7de73185, 0x7df81cea,
	0x7e07c0a3, 0x7e163efa, 0x7e23b587, 0x7e303dfd,
	0x7e3beec2, 0x7e46db77, 0x7e51155d, 0x7e5aabb3,
	0x7e63abf7, 0x7e6c222c, 0x7e741906, 0x7e7b9a18,
	0x7e82adfa, 0x7e895c63, 0x7e8fac4b, 0x7e95a3fb,
	0x7e9b4924, 0x7ea0a0ef, 0x7ea5b00d, 0x7eaa7ac3,
	0x7eaf04f3, 0x7eb3522a, 0x7eb765a5, 0x7ebb4259,
	0x7ebeeafd, 0x7ec2620a, 0x7ec5a9c4, 0x7ec8c441,
	0x7ecbb365, 0x7ece78ed, 0x7ed11671, 0x7ed38d62,
	0x7ed5df12, 0x7ed80cb4, 0x7eda175c, 0x7edc0005,
	0x7eddc78e, 0x7edf6ebf, 0x7ee0f647, 0x7ee25ebe,
	0x7ee3a8a9, 0x7ee4d473, 0x7ee5e276, 0x7ee6d2f5,
	0x7ee7a620, 0x7ee85c10, 0x7ee8f4cd, 0x7ee97047,
	0x7ee9ce59, 0x7eea0eca, 0x7eea3147, 0x7eea3568,
	0x7eea1aab, 0x7ee9e071, 0x7ee98602, 0x7ee90a88,
	0x7ee86d08, 0x7ee7ac6a, 0x7ee6c769, 0x7ee5bc9c,
	0x7ee48a67, 0x7ee32efc, 0x7ee1a857, 0x7edff42f,
	0x7ede0ffa, 0x7edbf8d9, 0x7ed9ab94, 0x7ed7248d,
	0x7ed45fae, 0x7ed1585c, 0x7ece095f, 0x7eca6ccb,
	0x7ec67be2, 0x7ec22eee, 0x7ebd7d1a, 0x7eb85c35,
	0x7eb2c075, 0x7eac9c20, 0x7ea5df27, 0x7e9e769f,
	0x7e964c16, 0x7e8d44ba, 0x7e834033, 0x7e781728,
	0x7e6b9933, 0x7e5d8a1a, 0x7e4d9ded, 0x7e3b737a,
	0x7e268c2f, 0x7e0e3ff5, 0x7df1aa5d, 0x7dcf8c72,
	0x7da61a1e, 0x7d72a0fb, 0x7d30e097, 0x7cd9b4ab,
	0x7c600f1a, 0x7ba90bdc, 0x7a722176, 0x77d664e5,
}

var normalW = [128]float32{
	1.72904052154178e-9, 1.2680928441913e-10, 1.68975177699895e-10, 1.98626884400521e-10,
	2.22324317905022e-10, 2.42449361237299e-10, 2.60161318991109e-10, 2.76119887103316e-10,
	2.90739628164604e-10, 3.04299704132156e-10, 3.16997952128716e-10, 3.28980205260966e-10,
	3.40357381208744e-10, 3.51216022127545e-10, 3.61625099496985e-10, 3.71640576341317e-10,
	3.81308564303125e-10, 3.90667568091865e-10, 3.99750118692429e-10, 4.08583986152762e-10,
	4.17193096394762e-10, 4.255982353393e-10, 4.33817597386126e-10, 4.41867218119051e-10,
	4.49761319620596e-10, 4.57512588939984e-10, 4.65132404808255e-10, 4.72631023842515e-10,
	4.8001773471779e-10, 4.87300986774536e-10, 4.9448849804868e-10, 5.01587346606859e-10,
	5.08604048237462e-10, 5.15544622914649e-10, 5.22414651965841e-10, 5.29219327495934e-10,
	5.35963495326683e-10, 5.42651692477542e-10, 5.49288180030166e-10, 5.5587697207172e-10,
	5.62421861294078e-10, 5.68926441730448e-10, 5.75394129033425e-10, 5.81828178635022e-10,
	5.88231702077216e-10, 5.94607681758562e-10, 6.00958984306954e-10, 6.07288372758972e-10,
	6.13598517701655e-10, 6.1989200751189e-10, 6.26171357811295e-10, 6.32439020239945e-10,
	6.3869739064003e-10, 6.44948816730244e-10, 6.51195605343024e-10, 6.57440029289462e-10,
	6.63684333910635e-10, 6.69930743369023e-10, 6.76181466729481e-10, 6.82438703875893e-10,
	6.88704651306895e-10, 6.94981507852029e-10, 7.01271480348217e-10, 7.07576789315497e-10,
	7.13899674670564e-10, 7.20242401516765e-10, 7.26607266049758e-10, 7.32996601619176e-10,
	7.39412784988247e-10, 7.45858242835513e-10, 7.52335458545542e-10, 7.58846979338992e-10,
	7.65395423796485e-10, 7.71983489835731e-10, 7.78613963207161e-10, 7.85289726580254e-10,
	7.92013769300795e-10, 7.98789197908769e-10, 8.05619247517662e-10, 8.12507294168871e-10,
	8.19456868290078e-10, 8.26471669404194e-10, 8.33555582256345e-10, 8.40712694550887e-10,
	8.47947316519453e-10, 8.55264002575252e-10, 8.62667575349606e-10, 8.70163152455139e-10,
	8.77756176378051e-10, 8.85452447971477e-10, 8.93258164105812e-10, 9.01179960133461e-10,
	9.09224957948965e-10, 9.17400820576452e-10, 9.2571581440189e-10, 9.34178880396749e-10,
	9.42799715964558e-10, 9.5158886939784e-10, 9.60557849381102e-10, 9.69719252543396e-10,
	9.79086912788916e-10, 9.88676077066823e-10, 9.98503613451617e-10, 1.00858825898955e-9,
	1.01895091686026e-9, 1.02961501519882e-9, 1.04060694369816e-9, 1.051956589271e-9,
	1.06369799919131e-9, 1.07587021016283e-9, 1.088518296059e-9, 1.1016947078118e-9,
	1.11546100955804e-9, 1.12989016134767e-9, 1.1450695700051e-9, 1.16110524260064e-9,
	1.17812756094405e-9, 1.19629950538354e-9, 1.21582869832805e-9, 1.23698562907902e-9,
	1.26013233005941e-9, 1.28576968441911e-9, 1.31462018496635e-9, 1.34778395621975e-9,
	1.38706353150542e-9, 1.4357403191804e-9, 1.50086590302112e-9, 1.60309479380802e-9,
}

var normalF = [128]float32{
	1.0, 0.963599693155768, 0.936282681708371, 0.913043647992038,
	0.892281650802303, 0.873243048926854, 0.855500607885064, 0.838783605310647,
	0.822907211395262, 0.807738294696121, 0.793177011783859, 0.779146085941703,
	0.765584173909236, 0.752441559185704, 0.739677243683338, 0.727256918354506,
	0.715151507420477, 0.703336099025817, 0.691789143446036, 0.680491841006414,
	0.669427667357706, 0.658582000058654, 0.647941821118551, 0.637495477343145,
	0.627232485257815, 0.617143370826562, 0.607219536632605, 0.597453150951812,
	0.587837054441821, 0.578364681126702, 0.569029991074722, 0.559827412710695,
	0.550751793121055, 0.541798355031724, 0.532962659389988, 0.524240572678993,
	0.515628238249872, 0.507122051081305, 0.498718635476584, 0.490414825289322,
	0.482207646334839, 0.47409430069825, 0.466072152694571, 0.458138716272872,
	0.450291643686927, 0.442528715280247, 0.434847830254662, 0.427246998309562,
	0.419724332054038, 0.412278040107025, 0.404906420811488, 0.397607856498043,
	0.390380808241389, 0.383223811059884, 0.376135469514454, 0.369114453668275,
	0.362159495373033, 0.355269384851547, 0.348442967549872, 0.341679141235014,
	0.334976853316971, 0.328335098376152, 0.321752915879209, 0.315229388068158,
	0.308763638009252, 0.30235482778948, 0.296002156849856, 0.28970486044581,
	0.283462208226013, 0.277273502921898, 0.271138079141025, 0.265055302258162,
	0.259024567398711, 0.253045298509766, 0.247116947514697, 0.241238993547751,
	0.235410942265728, 0.229632325234303, 0.223902699387134, 0.218221646556371,
	0.212588773073736, 0.207003709441874, 0.201466110076203, 0.19597565311811,
	0.190532040320914, 0.185134997010713, 0.179784272124962, 0.174479638332402,
	0.169220892238925, 0.164007854684928, 0.158840371140935, 0.153718312209587,
	0.148641574243697, 0.143610080091933, 0.138623779985851, 0.133682652584648,
	0.128786706197104, 0.123935980203982, 0.119130546708719, 0.114370512449888,
	0.109656021015818, 0.104987255410355, 0.100364441029546, 0.0957878491225782,
	0.0912578008276347, 0.086774671895543, 0.0823388982429574, 0.0779509825146547,
	0.0736115018847549, 0.0693211173941803, 0.0650805852136319, 0.0608907703485664,
	0.0567526634815386, 0.0526674019035032, 0.0486362958602841, 0.0446608622008724,
	0.0407428680747906, 0.0368843887869688, 0.0330878861465052, 0.0293563174402538,
	0.0256932919361496, 0.0221033046161116, 0.0185921027371658, 0.015167298010672,
	0.0118394786579823, 0.00862448441293047, 0.00554899522081647, 0.0026696290839025,
}
