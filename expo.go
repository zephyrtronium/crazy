package crazy

import "math"

// Exponential adapts a Source to generate random numbers under an exponential
// distribution.
type Exponential struct {
	Source
	z    *Ziggurat
	Rate float64
}

// NewExponential creates an exponential distribution drawing from the
// specified source with given rate parameter.
func NewExponential(src Source, rate float64) Exponential {
	return Exponential{
		Source: src,
		Rate:   rate,
		z:      expoZig,
	}
}

// Next generates an exponential variate.
func (e Exponential) Next() float64 {
	x := e.z.GenNext(e.Source)
	return x / e.Rate
}

func expoPDF(x float64) float64 {
	return math.Exp(-x)
}

func expoTail(src Source) float64 {
	return expoR - math.Log(Uniform0_1{src}.Next())
}

var expoZig = &Ziggurat{
	PDF:      expoPDF,
	Tail:     expoTail,
	Mirrored: false,
	K:        expoK,
	W:        expoW,
	F:        expoF,
}

const expoR = 6.89831511661564

// This time our tables are very different from math/rand/exp.go because that
// uses 256 segments for the ziggurat.

var expoK = [128]uint32{
	0x6fcb445e, 0x0, 0x4d737e37, 0x6139aec0,
	0x69f3d12b, 0x6ed02c13, 0x71e564e8, 0x74054185,
	0x75923ce5, 0x76c08125, 0x77ae33e4, 0x786df18f,
	0x790bd3b9, 0x79900e0b, 0x7a00617c, 0x7a60f586,
	0x7ab4dd05, 0x7afe6a88, 0x7b3f679d, 0x7b793a18,
	0x7bacfdcf, 0x7bdb96c4, 0x7c05be39, 0x7c2c0c3d,
	0x7c4efeca, 0x7c6eff21, 0x7c8c65df, 0x7ca77e24,
	0x7cc0880e, 0x7cd7baa9, 0x7ced4578, 0x7d0151b8,
	0x7d140358, 0x7d2579d5, 0x7d35d0db, 0x7d4520d9,
	0x7d537f74, 0x7d60ffe8, 0x7d6db35d, 0x7d79a929,
	0x7d84ef0d, 0x7d8f9165, 0x7d999b57, 0x7da316f0,
	0x7dac0d4c, 0x7db486aa, 0x7dbc8a89, 0x7dc41fb8,
	0x7dcb4c6d, 0x7dd2164e, 0x7dd88285, 0x7dde95ca,
	0x7de45469, 0x7de9c253, 0x7deee321, 0x7df3ba1c,
	0x7df84a45, 0x7dfc9659, 0x7e00a0d8, 0x7e046c07,
	0x7e07f9f5, 0x7e0b4c7f, 0x7e0e6553, 0x7e1145ef,
	0x7e13efab, 0x7e1663b3, 0x7e18a30e, 0x7e1aae9b,
	0x7e1c8717, 0x7e1e2d1a, 0x7e1fa118, 0x7e20e363,
	0x7e21f42c, 0x7e22d37d, 0x7e238141, 0x7e23fd3c,
	0x7e24470f, 0x7e245e34, 0x7e2441fd, 0x7e23f192,
	0x7e236bf1, 0x7e22afe8, 0x7e21bc12, 0x7e208ed7,
	0x7e1f2664, 0x7e1d80a7, 0x7e1b9b4b, 0x7e1973ae,
	0x7e1706dd, 0x7e145188, 0x7e114ffb, 0x7e0dfe0e,
	0x7e0a571a, 0x7e0655ea, 0x7e01f4a5, 0x7dfd2cb6,
	0x7df7f6ba, 0x7df24a56, 0x7dec1e19, 0x7de5674f,
	0x7dde19c8, 0x7dd62798, 0x7dcd80ca, 0x7dc412ff,
	0x7db9c8f1, 0x7dae89e5, 0x7da238ea, 0x7d94b3ed,
	0x7d85d286, 0x7d756471, 0x7d632f86, 0x7d4eed05,
	0x7d3845fe, 0x7d1ece56, 0x7d01fdd5, 0x7ce12650,
	0x7cbb653d, 0x7c8f8e2b, 0x7c5c0985, 0x7c1e9f69,
	0x7bd41ef5, 0x7b77c2ba, 0x7b020e4b, 0x7a6683fe,
	0x798e86fc, 0x784c37de, 0x7631178f, 0x71d70ee8,
}

var expoW = [128]float32{
	3.67793958476542e-9, 4.25332775818666e-11, 7.02928220559249e-11, 9.25425698651012e-11,
	1.11799711698583e-10, 1.29139548761729e-10, 1.4513110148658e-10, 1.60116324464922e-10,
	1.74318900065836e-10, 1.87894306851329e-10, 2.00955673138069e-10, 2.13588308305464e-10,
	2.25858420922826e-10, 2.37818627491997e-10, 2.49511583359182e-10, 2.609724613393e-10,
	2.7223069460054e-10, 2.83311233675164e-10, 2.94235473180608e-10, 3.0502194829816e-10,
	3.15686867174807e-10, 3.26244524099819e-10, 3.36707624532228e-10, 3.4708754393575e-10,
	3.5739453620911e-10, 3.67637903245841e-10, 3.77826134172222e-10, 3.87967020683321e-10,
	3.98067753356986e-10, 4.08135002696457e-10, 4.18174987814139e-10, 4.28193535039946e-10,
	4.38196128260386e-10, 4.4818795242902e-10, 4.58173931406425e-10, 4.68158761067606e-10,
	4.7814693844189e-10, 4.88142787513488e-10, 4.98150482201979e-10, 5.08174066954612e-10,
	5.18217475311942e-10, 5.28284546751281e-10, 5.38379042066019e-10, 5.48504657500921e-10,
	5.58665037832372e-10, 5.68863788556926e-10, 5.7910448733044e-10, 5.8939069478261e-10,
	5.99725964817428e-10, 6.10113854498201e-10, 6.20557933606133e-10, 6.31061793953567e-10,
	6.41629058526655e-10, 6.52263390527231e-10, 6.62968502379842e-10, 6.73748164767088e-10,
	6.8460621575462e-10, 6.95546570066147e-10, 7.06573228568673e-10, 7.17690288028809e-10,
	7.28901951202435e-10, 7.40212537322142e-10, 7.51626493049885e-10, 7.63148403966071e-10,
	7.74783006670982e-10, 7.86535201580078e-10, 7.98410066501379e-10, 8.10412871090942e-10,
	8.22549092291536e-10, 8.34824430870155e-10, 8.47244829182181e-10, 8.59816490304013e-10,
	8.72545898692189e-10, 8.85439842545622e-10, 8.98505438069087e-10, 9.11750155860903e-10,
	9.25181849676424e-10, 9.38808787852181e-10, 9.52639687714081e-10, 9.66683753337867e-10,
	9.80950717082385e-10, 9.9545088537727e-10, 1.01019518931839e-9, 1.02519524070859e-9,
	1.04046339428088e-9, 1.05601281695863e-9, 1.07185756514753e-9, 1.08801267122058e-9,
	1.10449424055701e-9, 1.12131956073616e-9, 1.13850722477669e-9, 1.15607727066324e-9,
	1.17405133983018e-9, 1.19245285779677e-9, 1.21130724079318e-9, 1.23064213301574e-9,
	1.25048768014499e-9, 1.27087684600675e-9, 1.29184578082892e-9, 1.31343425154456e-9,
	1.33568614714767e-9, 1.35865007540522e-9, 1.38238007151703e-9, 1.40693644494364e-9,
	1.43238679808051e-9, 1.45880726044375e-9, 1.48628399555461e-9, 1.51491505624071e-9,
	1.54481268980057e-9, 1.57610623071043e-9, 1.60894577037379e-9, 1.6435068688001e-9,
	1.67999668483269e-9, 1.71866207061178e-9, 1.75980043772427e-9, 1.80377461835613e-9,
	1.85103362490333e-9, 1.90214236118167e-9, 1.95782535678183e-9, 2.01903329886769e-9,
	2.08704828855707e-9, 2.16365844610263e-9, 2.25146504460128e-9, 2.35446456298163e-9,
	2.47926561773638e-9, 2.63800493724911e-9, 2.85692179785168e-9, 3.21227829745768e-9,
}

var expoF = [128]float32{
	1.0, 0.912707777475125, 0.859888382509951,
	0.819768204987318, 0.786558599390136, 0.757808116931513,
	0.732225562307742, 0.70903726877305, 0.687738233900834,
	0.667978059185555, 0.649502221835904, 0.632119133965389,
	0.615680409581244, 0.600068409921909, 0.585188041317763,
	0.570961159629425, 0.557322637577984, 0.544217529519877,
	0.531598981894556, 0.51942666330534, 0.507665564832616,
	0.496285069354161, 0.485258219764561, 0.474561136575439,
	0.464172549299398, 0.454073415617565, 0.444246609064009,
	0.434676660760544, 0.425349544207953, 0.41625249468547,
	0.407373856699979, 0.398702954344945, 0.390229980505289,
	0.381945901669019, 0.37384237574385, 0.365911680774217,
	0.358146652844753, 0.350540631765755, 0.343087413382875,
	0.335781207551399, 0.328616600975456, 0.321588524242527,
	0.314692222489879, 0.307923229226805, 0.301277342908602,
	0.294750605918, 0.288339285659537, 0.282039857514072,
	0.275848989435651, 0.269763528002485, 0.26378048575885,
	0.257897029705952, 0.252110470817996, 0.246418254475166,
	0.240817951718609, 0.235307251243941, 0.229883952059758,
	0.22454595674618, 0.219291265255944, 0.21411796920705,
	0.209024246621641, 0.204008357070756, 0.199068637188979,
	0.194203496526831, 0.189411413712133, 0.18469093289458,
	0.180040660450423, 0.175459261926495, 0.170945459204952,
	0.166498027871972, 0.162115794775345, 0.157797635757395,
	0.153542473551082, 0.149349275828338, 0.145217053390856,
	0.141144858494566, 0.137131783300012, 0.133176958441711,
	0.129279551710424, 0.125438766843016, 0.121653842415384,
	0.117924050834575, 0.114248697426989, 0.110627119620188,
	0.107058686216575, 0.103542796757883, 0.100078880980158,
	0.0966663983596901, 0.0933048377511747, 0.0899937171202717,
	0.0867325833737408, 0.0835210122914216, 0.0803586085655968,
	0.0772450059547157, 0.0741798675601352, 0.0711628862364988,
	0.0681937851487071, 0.0652723184912142, 0.0623982723887404,
	0.0595714660015745, 0.0567917528636362, 0.0540590224876553,
	0.0513732022795311, 0.048734259813628, 0.0461422055330636,
	0.043597095954805, 0.0410990374797804, 0.0386481909348909,
	0.0362447770091108, 0.0338890827931727, 0.0315814696966008,
	0.0293223831044655, 0.0271123642604352, 0.0249520650399684,
	0.0228422665356619, 0.0207839027613827, 0.0187780913696083,
	0.0168261742012834, 0.0149297719923083, 0.0130908601063722,
	0.0113118766723302, 0.00959588294047461, 0.00794681255393166,
	0.00636988317909038, 0.00487233313651832, 0.0034648966761225,
	0.00216530760995328, 0.00100948486124341,
}
