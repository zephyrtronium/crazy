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
	dist := Uniform1_2{src}
	for {
		x := -math.Log(dist.Next()-1) / normalR
		y := -math.Log(dist.Next() - 1)
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
	0xed5a4424, 0x0, 0xc01e36a7, 0xd9c88f4d, 0xe4b68d43,
	0xeac00a3a, 0xee9243d6, 0xf1344b7a, 0xf3208b87, 0xf4979cba,
	0xf5bec53e, 0xf6ad054c, 0xf771518c, 0xf815ce44, 0xf8a199ce,
	0xf919d8b6, 0xf98259ad, 0xf9ddfda5, 0xfa2efc16, 0xfa771106,
	0xfab79cd9, 0xfaf1bac8, 0xfb26510c, 0xfb561caf, 0xfb81ba60,
	0xfba9ad11, 0xfbce630a, 0xfbf039d4, 0xfc0f8147, 0xfc2c7df4,
	0xfc476b0f, 0xfc607bfb, 0xfc77dd85, 0xfc8db6ee, 0xfca22abb,
	0xfcb55766, 0xfcc757ef, 0xfcd84459, 0xfce8320c, 0xfcf73431,
	0xfd055bf4, 0xfd12b8c7, 0xfd1f5897, 0xfd2b47f6, 0xfd369248,
	0xfd4141de, 0xfd4b601b, 0xfd54f587, 0xfd5e09e7, 0xfd66a455,
	0xfd6ecb4b, 0xfd7684b3, 0xfd7dd5fa, 0xfd84c414, 0xfd8b5389,
	0xfd918882, 0xfd9766ca, 0xfd9cf1db, 0xfda22ce3, 0xfda71ac5,
	0xfdabbe25, 0xfdb01968, 0xfdb42eb9, 0xfdb8000a, 0xfdbb8f1d,
	0xfdbedd7e, 0xfdc1ec8e, 0xfdc4bd7d, 0xfdc75152, 0xfdc9a8e6,
	0xfdcbc4ec, 0xfdcda5eb, 0xfdcf4c40, 0xfdd0b821, 0xfdd1e99b,
	0xfdd2e08e, 0xfdd39cb3, 0xfdd41d95, 0xfdd4628f, 0xfdd46ad1,
	0xfdd43556, 0xfdd3c0e2, 0xfdd30c05, 0xfdd21510, 0xfdd0da11,
	0xfdcf58d4, 0xfdcd8ed3, 0xfdcb7938, 0xfdc914ce, 0xfdc65df9,
	0xfdc350ae, 0xfdbfe85f, 0xfdbc1ff4, 0xfdb7f1b2, 0xfdb35729,
	0xfdae491a, 0xfda8bf5c, 0xfda2b0b8, 0xfd9c12be, 0xfd94d996,
	0xfd8cf7c4, 0xfd845ddd, 0xfd7afa35, 0xfd70b86a, 0xfd6580ea,
	0xfd593840, 0xfd4bbe4f, 0xfd3ced3f, 0xfd2c982d, 0xfd1a8974,
	0xfd068067, 0xfcf02e51, 0xfcd73266, 0xfcbb1434, 0xfc9b3bdb,
	0xfc76e6f4, 0xfc4d185e, 0xfc1c7fea, 0xfbe354bb, 0xfb9f18e4,
	0xfb4c343c, 0xfae541f7, 0xfa61c12e, 0xf9b36957, 0xf8c01e35,
	0xf75217b8, 0xf4e442ec, 0xefacc9cb,
}

var normalW = [128]float32{
	8.6452026077089e-10, 6.34046422095652e-11, 8.44875888499474e-11, 9.93134422002605e-11,
	1.11162158952511e-10, 1.2122468061865e-10, 1.30080659495554e-10, 1.38059943551658e-10,
	1.45369814082302e-10, 1.52149852066078e-10, 1.58498976064358e-10, 1.64490102630483e-10,
	1.70178690604372e-10, 1.75608011063773e-10, 1.80812549748493e-10, 1.85820288170659e-10,
	1.90654282151563e-10, 1.95333784045933e-10, 1.99875059346215e-10, 2.04291993076381e-10,
	2.08596548197381e-10, 2.1279911766965e-10, 2.16908798693063e-10, 2.20933609059526e-10,
	2.24880659810298e-10, 2.28756294469992e-10, 2.32566202404127e-10, 2.36315511921258e-10,
	2.40008867358895e-10, 2.43650493387268e-10, 2.4724424902434e-10, 2.50793673303429e-10,
	2.54302024118731e-10, 2.57772311457325e-10, 2.6120732598292e-10, 2.64609663747967e-10,
	2.67981747663341e-10, 2.71325846238771e-10, 2.74644090015083e-10, 2.7793848603586e-10,
	2.81210930647039e-10, 2.84463220865224e-10, 2.87697064516712e-10, 2.90914089317511e-10,
	2.94115851038608e-10, 2.97303840879281e-10, 3.00479492153477e-10, 3.03644186379486e-10,
	3.06799258850828e-10, 3.09946003755945e-10, 3.13085678905647e-10, 3.16219510119972e-10,
	3.19348695320015e-10, 3.22474408365122e-10, 3.25597802671512e-10, 3.28720014644731e-10,
	3.31842166955318e-10, 3.34965371684511e-10, 3.38090733364741e-10, 3.41219351937947e-10,
	3.44352325653447e-10, 3.47490753926014e-10, 3.50635740174109e-10, 3.53788394657748e-10,
	3.56949837335282e-10, 3.60121200758383e-10, 3.63303633024879e-10, 3.66498300809588e-10,
	3.69706392494124e-10, 3.72929121417756e-10, 3.76167729272771e-10, 3.79423489669496e-10,
	3.82697711898243e-10, 3.85991744917866e-10, 3.89306981603581e-10, 3.92644863290127e-10,
	3.96006884650397e-10, 3.99394598954384e-10, 4.02809623758831e-10, 4.06253647084436e-10,
	4.09728434145039e-10, 4.13235834702097e-10, 4.16777791128172e-10, 4.20356347275444e-10,
	4.23973658259726e-10, 4.27632001287626e-10, 4.31333787674803e-10, 4.35081576227569e-10,
	4.38878088189026e-10, 4.42726223985738e-10, 4.46629082052906e-10, 4.50589980066731e-10,
	4.54612478974482e-10, 4.58700410288226e-10, 4.62857907200945e-10, 4.67089440198375e-10,
	4.71399857982279e-10, 4.7579443469892e-10, 4.80278924690551e-10, 4.84859626271698e-10,
	4.89543456394458e-10, 4.94338038533411e-10, 4.99251806725809e-10, 5.04294129494773e-10,
	5.09475458430131e-10, 5.14807507599408e-10, 5.2030347184908e-10, 5.25978294635501e-10,
	5.31848999595655e-10, 5.37935105081415e-10, 5.44259148029501e-10, 5.50847353905902e-10,
	5.57730504779021e-10, 5.64945080673837e-10, 5.72534785002552e-10, 5.80552621300321e-10,
	5.89063780472024e-10, 5.9814975269177e-10, 6.07914349164025e-10, 6.18492814539511e-10,
	6.30066165029705e-10, 6.42884842209554e-10, 6.57310092483174e-10, 6.73891978109877e-10,
	6.93531765752709e-10, 7.17870159590201e-10, 7.50432951510562e-10, 8.01547396904009e-10,
}

var normalF = [128]float32{
	1.0, 0.963599693155768, 0.936282681708371, 0.913043647992038, 0.892281650802303,
	0.873243048926854, 0.855500607885064, 0.838783605310647, 0.822907211395262,
	0.807738294696121, 0.793177011783859, 0.779146085941703, 0.765584173909236,
	0.752441559185704, 0.739677243683338, 0.727256918354506, 0.715151507420477,
	0.703336099025817, 0.691789143446036, 0.680491841006414, 0.669427667357706,
	0.658582000058654, 0.647941821118551, 0.637495477343145, 0.627232485257815,
	0.617143370826562, 0.607219536632605, 0.597453150951812, 0.587837054441821,
	0.578364681126702, 0.569029991074722, 0.559827412710695, 0.550751793121055,
	0.541798355031724, 0.532962659389988, 0.524240572678993, 0.515628238249872,
	0.507122051081305, 0.498718635476584, 0.490414825289322, 0.482207646334839,
	0.47409430069825, 0.466072152694571, 0.458138716272872, 0.450291643686927,
	0.442528715280247, 0.434847830254662, 0.427246998309562, 0.419724332054038,
	0.412278040107025, 0.404906420811488, 0.397607856498043, 0.390380808241389,
	0.383223811059884, 0.376135469514454, 0.369114453668275, 0.362159495373033,
	0.355269384851547, 0.348442967549872, 0.341679141235014, 0.334976853316971,
	0.328335098376152, 0.321752915879209, 0.315229388068158, 0.308763638009252,
	0.30235482778948, 0.296002156849856, 0.28970486044581, 0.283462208226013,
	0.277273502921898, 0.271138079141025, 0.265055302258162, 0.259024567398711,
	0.253045298509766, 0.247116947514697, 0.241238993547751, 0.235410942265728,
	0.229632325234303, 0.223902699387134, 0.218221646556371, 0.212588773073736,
	0.207003709441874, 0.201466110076203, 0.19597565311811, 0.190532040320914,
	0.185134997010713, 0.179784272124962, 0.174479638332402, 0.169220892238925,
	0.164007854684928, 0.158840371140935, 0.153718312209587, 0.148641574243697,
	0.143610080091933, 0.138623779985851, 0.133682652584648, 0.128786706197104,
	0.123935980203982, 0.119130546708719, 0.114370512449888, 0.109656021015818,
	0.104987255410355, 0.100364441029546, 0.0957878491225782, 0.0912578008276347,
	0.086774671895543, 0.0823388982429574, 0.0779509825146547, 0.0736115018847549,
	0.0693211173941803, 0.0650805852136319, 0.0608907703485664, 0.0567526634815386,
	0.0526674019035032, 0.0486362958602841, 0.0446608622008724, 0.0407428680747906,
	0.0368843887869688, 0.0330878861465052, 0.0293563174402538, 0.0256932919361496,
	0.0221033046161116, 0.0185921027371658, 0.015167298010672, 0.0118394786579823,
	0.00862448441293047, 0.00554899522081647, 0.0026696290839025,
}
