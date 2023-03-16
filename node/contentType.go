/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package node

import "sync"

var contentType sync.Map

func init() {
	contentType.Store(".323", "text/h323")
	contentType.Store(".323", "text/h323")
	contentType.Store(".3gp", "video/3gpp")
	contentType.Store(".aab", "application/x-authoware-bin")
	contentType.Store(".aam", "application/x-authoware-map")
	contentType.Store(".aas", "application/x-authoware-seg")
	contentType.Store(".acx", "application/internet-property-stream")
	contentType.Store(".ai", "application/postscript")
	contentType.Store(".aif", "audio/x-aiff")
	contentType.Store(".aifc", "audio/x-aiff")
	contentType.Store(".aiff", "audio/x-aiff")
	contentType.Store(".als", "audio/X-Alpha5")
	contentType.Store(".amc", "application/x-mpeg")
	contentType.Store(".ani", "application/octet-stream")
	contentType.Store(".apk", "application/vnd.android.package-archive")
	contentType.Store(".asc", "text/plain")
	contentType.Store(".asd", "application/astound")
	contentType.Store(".asf", "video/x-ms-asf")
	contentType.Store(".asn", "application/astound")
	contentType.Store(".asp", "application/x-asap")
	contentType.Store(".asr", "video/x-ms-asf")
	contentType.Store(".asx", "video/x-ms-asf")
	contentType.Store(".au", "audio/basic")
	contentType.Store(".avb", "application/octet-stream")
	contentType.Store(".avi", "video/x-msvideo")
	contentType.Store(".awb", "audio/amr-wb")
	contentType.Store(".axs", "application/olescript")
	contentType.Store(".bas", "text/plain")
	contentType.Store(".bcpio", "application/x-bcpio")
	contentType.Store(".bin", "application/octet-stream")
	contentType.Store(".bld", "application/bld")
	contentType.Store(".bld2", "application/bld2")
	contentType.Store(".bmp", "image/bmp")
	contentType.Store(".bpk", "application/octet-stream")
	contentType.Store(".bz2", "application/x-bzip2")
	contentType.Store(".c", "text/plain")
	contentType.Store(".cal", "image/x-cals")
	contentType.Store(".cat", "application/vnd.ms-pkiseccat")
	contentType.Store(".ccn", "application/x-cnc")
	contentType.Store(".cco", "application/x-cocoa")
	contentType.Store(".cdf", "application/x-cdf")
	contentType.Store(".cer", "application/x-x509-ca-cert")
	contentType.Store(".cgi", "magnus-internal/cgi")
	contentType.Store(".chat", "application/x-chat")
	contentType.Store(".class", "application/octet-stream")
	contentType.Store(".clp", "application/x-msclip")
	contentType.Store(".cmx", "image/x-cmx")
	contentType.Store(".co", "application/x-cult3d-object")
	contentType.Store(".cod", "image/cis-cod")
	contentType.Store(".conf", "text/plain")
	contentType.Store(".cpio", "application/x-cpio")
	contentType.Store(".cpp", "text/plain")
	contentType.Store(".cpt", "application/mac-compactpro")
	contentType.Store(".crd", "application/x-mscardfile")
	contentType.Store(".crl", "application/pkix-crl")
	contentType.Store(".crt", "application/x-x509-ca-cert")
	contentType.Store(".csh", "application/x-csh")
	contentType.Store(".csm", "chemical/x-csml")
	contentType.Store(".csml", "chemical/x-csml")
	contentType.Store(".css", "text/css")
	contentType.Store(".cur", "application/octet-stream")
	contentType.Store(".dcm", "x-lml/x-evm")
	contentType.Store(".dcr", "application/x-director")
	contentType.Store(".dcx", "image/x-dcx")
	contentType.Store(".der", "application/x-x509-ca-cert")
	contentType.Store(".dhtml", "text/html")
	contentType.Store(".dir", "application/x-director")
	contentType.Store(".dll", "application/x-msdownload")
	contentType.Store(".dmg", "application/octet-stream")
	contentType.Store(".dms", "application/octet-stream")
	contentType.Store(".doc", "application/msword")
	contentType.Store(".docx", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	contentType.Store(".dot", "application/msword")
	contentType.Store(".dvi", "application/x-dvi")
	contentType.Store(".dwf", "drawing/x-dwf")
	contentType.Store(".dwg", "application/x-autocad")
	contentType.Store(".dxf", "application/x-autocad")
	contentType.Store(".dxr", "application/x-director")
	contentType.Store(".ebk", "application/x-expandedbook")
	contentType.Store(".emb", "chemical/x-embl-dl-nucleotide")
	contentType.Store(".embl", "chemical/x-embl-dl-nucleotide")
	contentType.Store(".eps", "application/postscript")
	contentType.Store(".epub", "application/epub+zip")
	contentType.Store(".eri", "image/x-eri")
	contentType.Store(".es", "audio/echospeech")
	contentType.Store(".esl", "audio/echospeech")
	contentType.Store(".etc", "application/x-earthtime")
	contentType.Store(".etx", "text/x-setext")
	contentType.Store(".evm", "x-lml/x-evm")
	contentType.Store(".evy", "application/envoy")
	contentType.Store(".exe", "application/octet-stream")
	contentType.Store(".fh4", "image/x-freehand")
	contentType.Store(".fh5", "image/x-freehand")
	contentType.Store(".fhc", "image/x-freehand")
	contentType.Store(".fif", "application/fractals")
	contentType.Store(".flr", "x-world/x-vrml")
	contentType.Store(".flv", "flv-application/octet-stream")
	contentType.Store(".fm", "application/x-maker")
	contentType.Store(".fpx", "image/x-fpx")
	contentType.Store(".fvi", "video/isivideo")
	contentType.Store(".gau", "chemical/x-gaussian-input")
	contentType.Store(".gca", "application/x-gca-compressed")
	contentType.Store(".gdb", "x-lml/x-gdb")
	contentType.Store(".gif", "image/gif")
	contentType.Store(".gps", "application/x-gps")
	contentType.Store(".gtar", "application/x-gtar")
	contentType.Store(".gz", "application/x-gzip")
	contentType.Store(".h", "text/plain")
	contentType.Store(".hdf", "application/x-hdf")
	contentType.Store(".hdm", "text/x-hdml")
	contentType.Store(".hdml", "text/x-hdml")
	contentType.Store(".hlp", "application/winhlp")
	contentType.Store(".hqx", "application/mac-binhex40")
	contentType.Store(".hta", "application/hta")
	contentType.Store(".htc", "text/x-component")
	contentType.Store(".htm", "text/html")
	contentType.Store(".html", "text/html")
	contentType.Store(".hts", "text/html")
	contentType.Store(".htt", "text/webviewhtml")
	contentType.Store(".ice", "x-conference/x-cooltalk")
	contentType.Store(".ico", "image/x-icon")
	contentType.Store(".ief", "image/ief")
	contentType.Store(".ifm", "image/gif")
	contentType.Store(".ifs", "image/ifs")
	contentType.Store(".iii", "application/x-iphone")
	contentType.Store(".imy", "audio/melody")
	contentType.Store(".ins", "application/x-internet-signup")
	contentType.Store(".ips", "application/x-ipscript")
	contentType.Store(".ipx", "application/x-ipix")
	contentType.Store(".isp", "application/x-internet-signup")
	contentType.Store(".it", "audio/x-mod")
	contentType.Store(".itz", "audio/x-mod")
	contentType.Store(".ivr", "i-world/i-vrml")
	contentType.Store(".j2k", "image/j2k")
	contentType.Store(".jad", "text/vnd.sun.j2me.app-descriptor")
	contentType.Store(".jam", "application/x-jam")
	contentType.Store(".jar", "application/java-archive")
	contentType.Store(".java", "text/plain")
	contentType.Store(".jfif", "image/pipeg")
	contentType.Store(".jnlp", "application/x-java-jnlp-file")
	contentType.Store(".jpe", "image/jpeg")
	contentType.Store(".jpeg", "image/jpeg")
	contentType.Store(".jpg", "image/jpeg")
	contentType.Store(".jpz", "image/jpeg")
	contentType.Store(".js", "application/x-javascript")
	contentType.Store(".jwc", "application/jwc")
	contentType.Store(".kjx", "application/x-kjx")
	contentType.Store(".lak", "x-lml/x-lak")
	contentType.Store(".latex", "application/x-latex")
	contentType.Store(".lcc", "application/fastman")
	contentType.Store(".lcl", "application/x-digitalloca")
	contentType.Store(".lcr", "application/x-digitalloca")
	contentType.Store(".lgh", "application/lgh")
	contentType.Store(".lha", "application/octet-stream")
	contentType.Store(".lml", "x-lml/x-lml")
	contentType.Store(".lmlpack", "x-lml/x-lmlpack")
	contentType.Store(".log", "text/plain")
	contentType.Store(".lsf", "video/x-la-asf")
	contentType.Store(".lsx", "video/x-la-asf")
	contentType.Store(".lzh", "application/octet-stream")
	contentType.Store(".m13", "application/x-msmediaview")
	contentType.Store(".m14", "application/x-msmediaview")
	contentType.Store(".m15", "audio/x-mod")
	contentType.Store(".m3u", "audio/x-mpegurl")
	contentType.Store(".m3url", "audio/x-mpegurl")
	contentType.Store(".m4a", "audio/mp4a-latm")
	contentType.Store(".m4b", "audio/mp4a-latm")
	contentType.Store(".m4p", "audio/mp4a-latm")
	contentType.Store(".m4u", "video/vnd.mpegurl")
	contentType.Store(".m4v", "video/x-m4v")
	contentType.Store(".ma1", "audio/ma1")
	contentType.Store(".ma2", "audio/ma2")
	contentType.Store(".ma3", "audio/ma3")
	contentType.Store(".ma5", "audio/ma5")
	contentType.Store(".man", "application/x-troff-man")
	contentType.Store(".map", "magnus-internal/imagemap")
	contentType.Store(".mbd", "application/mbedlet")
	contentType.Store(".mct", "application/x-mascot")
	contentType.Store(".mdb", "application/x-msaccess")
	contentType.Store(".mdz", "audio/x-mod")
	contentType.Store(".me", "application/x-troff-me")
	contentType.Store(".mel", "text/x-vmel")
	contentType.Store(".mht", "message/rfc822")
	contentType.Store(".mhtml", "message/rfc822")
	contentType.Store(".mi", "application/x-mif")
	contentType.Store(".mid", "audio/mid")
	contentType.Store(".midi", "audio/midi")
	contentType.Store(".mif", "application/x-mif")
	contentType.Store(".mil", "image/x-cals")
	contentType.Store(".mio", "audio/x-mio")
	contentType.Store(".mmf", "application/x-skt-lbs")
	contentType.Store(".mng", "video/x-mng")
	contentType.Store(".mny", "application/x-msmoney")
	contentType.Store(".moc", "application/x-mocha")
	contentType.Store(".mocha", "application/x-mocha")
	contentType.Store(".mod", "audio/x-mod")
	contentType.Store(".mof", "application/x-yumekara")
	contentType.Store(".mol", "chemical/x-mdl-molfile")
	contentType.Store(".mop", "chemical/x-mopac-input")
	contentType.Store(".mov", "video/quicktime")
	contentType.Store(".movie", "video/x-sgi-movie")
	contentType.Store(".mp2", "video/mpeg")
	contentType.Store(".mp3", "audio/mpeg")
	contentType.Store(".mp4", "video/mp4")
	contentType.Store(".mpa", "video/mpeg")
	contentType.Store(".mpc", "application/vnd.mpohun.certificate")
	contentType.Store(".mpe", "video/mpeg")
	contentType.Store(".mpeg", "video/mpeg")
	contentType.Store(".mpg", "video/mpeg")
	contentType.Store(".mpg4", "video/mp4")
	contentType.Store(".mpga", "audio/mpeg")
	contentType.Store(".mpn", "application/vnd.mophun.application")
	contentType.Store(".mpp", "application/vnd.ms-project")
	contentType.Store(".mps", "application/x-mapserver")
	contentType.Store(".mpv2", "video/mpeg")
	contentType.Store(".mrl", "text/x-mrml")
	contentType.Store(".mrm", "application/x-mrm")
	contentType.Store(".ms", "application/x-troff-ms")
	contentType.Store(".msg", "application/vnd.ms-outlook")
	contentType.Store(".mts", "application/metastream")
	contentType.Store(".mtx", "application/metastream")
	contentType.Store(".mtz", "application/metastream")
	contentType.Store(".mvb", "application/x-msmediaview")
	contentType.Store(".mzv", "application/metastream")
	contentType.Store(".nar", "application/zip")
	contentType.Store(".nbmp", "image/nbmp")
	contentType.Store(".nc", "application/x-netcdf")
	contentType.Store(".ndb", "x-lml/x-ndb")
	contentType.Store(".ndwn", "application/ndwn")
	contentType.Store(".nif", "application/x-nif")
	contentType.Store(".nmz", "application/x-scream")
	contentType.Store(".nokia-op-logo", "image/vnd.nok-oplogo-color")
	contentType.Store(".npx", "application/x-netfpx")
	contentType.Store(".nsnd", "audio/nsnd")
	contentType.Store(".nva", "application/x-neva1")
	contentType.Store(".nws", "message/rfc822")
	contentType.Store(".oda", "application/oda")
	contentType.Store(".ogg", "audio/ogg")
	contentType.Store(".oom", "application/x-AtlasMate-Plugin")
	contentType.Store(".p10", "application/pkcs10")
	contentType.Store(".p12", "application/x-pkcs12")
	contentType.Store(".p7b", "application/x-pkcs7-certificates")
	contentType.Store(".p7c", "application/x-pkcs7-mime")
	contentType.Store(".p7m", "application/x-pkcs7-mime")
	contentType.Store(".p7r", "application/x-pkcs7-certreqresp")
	contentType.Store(".p7s", "application/x-pkcs7-signature")
	contentType.Store(".pac", "audio/x-pac")
	contentType.Store(".pae", "audio/x-epac")
	contentType.Store(".pan", "application/x-pan")
	contentType.Store(".pbm", "image/x-portable-bitmap")
	contentType.Store(".pcx", "image/x-pcx")
	contentType.Store(".pda", "image/x-pda")
	contentType.Store(".pdb", "chemical/x-pdb")
	contentType.Store(".pdf", "application/pdf")
	contentType.Store(".pfr", "application/font-tdpfr")
	contentType.Store(".pfx", "application/x-pkcs12")
	contentType.Store(".pgm", "image/x-portable-graymap")
	contentType.Store(".pict", "image/x-pict")
	contentType.Store(".pko", "application/ynd.ms-pkipko")
	contentType.Store(".pm", "application/x-perl")
	contentType.Store(".pma", "application/x-perfmon")
	contentType.Store(".pmc", "application/x-perfmon")
	contentType.Store(".pmd", "application/x-pmd")
	contentType.Store(".pml", "application/x-perfmon")
	contentType.Store(".pmr", "application/x-perfmon")
	contentType.Store(".pmw", "application/x-perfmon")
	contentType.Store(".png", "image/png")
	contentType.Store(".pnm", "image/x-portable-anymap")
	contentType.Store(".pnz", "image/png")
	contentType.Store(".pot", "application/vnd.ms-powerpoint")
	contentType.Store(".ppm", "image/x-portable-pixmap")
	contentType.Store(".pps", "application/vnd.ms-powerpoint")
	contentType.Store(".ppt", "application/vnd.ms-powerpoint")
	contentType.Store(".pptx", "application/vnd.openxmlformats-officedocument.presentationml.presentation")
	contentType.Store(".pqf", "application/x-cprplayer")
	contentType.Store(".pqi", "application/cprplayer")
	contentType.Store(".prc", "application/x-prc")
	contentType.Store(".prf", "application/pics-rules")
	contentType.Store(".prop", "text/plain")
	contentType.Store(".proxy", "application/x-ns-proxy-autoconfig")
	contentType.Store(".ps", "application/postscript")
	contentType.Store(".ptlk", "application/listenup")
	contentType.Store(".pub", "application/x-mspublisher")
	contentType.Store(".pvx", "video/x-pv-pvx")
	contentType.Store(".qcp", "audio/vnd.qcelp")
	contentType.Store(".qt", "video/quicktime")
	contentType.Store(".qti", "image/x-quicktime")
	contentType.Store(".qtif", "image/x-quicktime")
	contentType.Store(".r3t", "text/vnd.rn-realtext3d")
	contentType.Store(".ra", "audio/x-pn-realaudio")
	contentType.Store(".ram", "audio/x-pn-realaudio")
	contentType.Store(".rar", "application/octet-stream")
	contentType.Store(".ras", "image/x-cmu-raster")
	contentType.Store(".rc", "text/plain")
	contentType.Store(".rdf", "application/rdf+xml")
	contentType.Store(".rf", "image/vnd.rn-realflash")
	contentType.Store(".rgb", "image/x-rgb")
	contentType.Store(".rlf", "application/x-richlink")
	contentType.Store(".rm", "audio/x-pn-realaudio")
	contentType.Store(".rmf", "audio/x-rmf")
	contentType.Store(".rmi", "audio/mid")
	contentType.Store(".rmm", "audio/x-pn-realaudio")
	contentType.Store(".rmvb", "audio/x-pn-realaudio")
	contentType.Store(".rnx", "application/vnd.rn-realplayer")
	contentType.Store(".roff", "application/x-troff")
	contentType.Store(".rp", "image/vnd.rn-realpix")
	contentType.Store(".rpm", "audio/x-pn-realaudio-plugin")
	contentType.Store(".rt", "text/vnd.rn-realtext")
	contentType.Store(".rte", "x-lml/x-gps")
	contentType.Store(".rtf", "application/rtf")
	contentType.Store(".rtg", "application/metastream")
	contentType.Store(".rtx", "text/richtext")
	contentType.Store(".rv", "video/vnd.rn-realvideo")
	contentType.Store(".rwc", "application/x-rogerwilco")
	contentType.Store(".s3m", "audio/x-mod")
	contentType.Store(".s3z", "audio/x-mod")
	contentType.Store(".sca", "application/x-supercard")
	contentType.Store(".scd", "application/x-msschedule")
	contentType.Store(".sct", "text/scriptlet")
	contentType.Store(".sdf", "application/e-score")
	contentType.Store(".sea", "application/x-stuffit")
	contentType.Store(".setpay", "application/set-payment-initiation")
	contentType.Store(".setreg", "application/set-registration-initiation")
	contentType.Store(".sgm", "text/x-sgml")
	contentType.Store(".sgml", "text/x-sgml")
	contentType.Store(".sh", "application/x-sh")
	contentType.Store(".shar", "application/x-shar")
	contentType.Store(".shtml", "magnus-internal/parsed-html")
	contentType.Store(".shw", "application/presentations")
	contentType.Store(".si6", "image/si6")
	contentType.Store(".si7", "image/vnd.stiwap.sis")
	contentType.Store(".si9", "image/vnd.lgtwap.sis")
	contentType.Store(".sis", "application/vnd.symbian.install")
	contentType.Store(".sit", "application/x-stuffit")
	contentType.Store(".skd", "application/x-Koan")
	contentType.Store(".skm", "application/x-Koan")
	contentType.Store(".skp", "application/x-Koan")
	contentType.Store(".skt", "application/x-Koan")
	contentType.Store(".slc", "application/x-salsa")
	contentType.Store(".smd", "audio/x-smd")
	contentType.Store(".smi", "application/smil")
	contentType.Store(".smil", "application/smil")
	contentType.Store(".smp", "application/studiom")
	contentType.Store(".smz", "audio/x-smd")
	contentType.Store(".snd", "audio/basic")
	contentType.Store(".spc", "application/x-pkcs7-certificates")
	contentType.Store(".spl", "application/futuresplash")
	contentType.Store(".spr", "application/x-sprite")
	contentType.Store(".sprite", "application/x-sprite")
	contentType.Store(".sdp", "application/sdp")
	contentType.Store(".spt", "application/x-spt")
	contentType.Store(".src", "application/x-wais-source")
	contentType.Store(".sst", "application/vnd.ms-pkicertstore")
	contentType.Store(".stk", "application/hyperstudio")
	contentType.Store(".stl", "application/vnd.ms-pkistl")
	contentType.Store(".stm", "text/html")
	contentType.Store(".svg", "image/svg+xml")
	contentType.Store(".sv4cpio", "application/x-sv4cpio")
	contentType.Store(".sv4crc", "application/x-sv4crc")
	contentType.Store(".svf", "image/vnd")
	contentType.Store(".svg", "image/svg+xml")
	contentType.Store(".svh", "image/svh")
	contentType.Store(".svr", "x-world/x-svr")
	contentType.Store(".swf", "application/x-shockwave-flash")
	contentType.Store(".swfl", "application/x-shockwave-flash")
	contentType.Store(".t", "application/x-troff")
	contentType.Store(".tad", "application/octet-stream")
	contentType.Store(".talk", "text/x-speech")
	contentType.Store(".tar", "application/x-tar")
	contentType.Store(".taz", "application/x-tar")
	contentType.Store(".tbp", "application/x-timbuktu")
	contentType.Store(".tbt", "application/x-timbuktu")
	contentType.Store(".tcl", "application/x-tcl")
	contentType.Store(".tex", "application/x-tex")
	contentType.Store(".texi", "application/x-texinfo")
	contentType.Store(".texinfo", "application/x-texinfo")
	contentType.Store(".tgz", "application/x-compressed")
	contentType.Store(".thm", "application/vnd.eri.thm")
	contentType.Store(".tif", "image/tiff")
	contentType.Store(".tiff", "image/tiff")
	contentType.Store(".tki", "application/x-tkined")
	contentType.Store(".tkined", "application/x-tkined")
	contentType.Store(".toc", "application/toc")
	contentType.Store(".toy", "image/toy")
	contentType.Store(".tr", "application/x-troff")
	contentType.Store(".trk", "x-lml/x-gps")
	contentType.Store(".trm", "application/x-msterminal")
	contentType.Store(".tsi", "audio/tsplayer")
	contentType.Store(".tsp", "application/dsptype")
	contentType.Store(".tsv", "text/tab-separated-values")
	contentType.Store(".ttf", "application/octet-stream")
	contentType.Store(".ttz", "application/t-time")
	contentType.Store(".txt", "text/plain")
	contentType.Store(".uls", "text/iuls")
	contentType.Store(".ult", "audio/x-mod")
	contentType.Store(".ustar", "application/x-ustar")
	contentType.Store(".uu", "application/x-uuencode")
	contentType.Store(".uue", "application/x-uuencode")
	contentType.Store(".vcd", "application/x-cdlink")
	contentType.Store(".vcf", "text/x-vcard")
	contentType.Store(".vdo", "video/vdo")
	contentType.Store(".vib", "audio/vib")
	contentType.Store(".viv", "video/vivo")
	contentType.Store(".vivo", "video/vivo")
	contentType.Store(".vmd", "application/vocaltec-media-desc")
	contentType.Store(".vmf", "application/vocaltec-media-file")
	contentType.Store(".vmi", "application/x-dreamcast-vms-info")
	contentType.Store(".vms", "application/x-dreamcast-vms")
	contentType.Store(".vox", "audio/voxware")
	contentType.Store(".vqe", "audio/x-twinvq-plugin")
	contentType.Store(".vqf", "audio/x-twinvq")
	contentType.Store(".vql", "audio/x-twinvq")
	contentType.Store(".vre", "x-world/x-vream")
	contentType.Store(".vrml", "x-world/x-vrml")
	contentType.Store(".vrt", "x-world/x-vrt")
	contentType.Store(".vrw", "x-world/x-vream")
	contentType.Store(".vts", "workbook/formulaone")
	contentType.Store(".wav", "audio/x-wav")
	contentType.Store(".wax", "audio/x-ms-wax")
	contentType.Store(".wbmp", "image/vnd.wap.wbmp")
	contentType.Store(".wcm", "application/vnd.ms-works")
	contentType.Store(".wdb", "application/vnd.ms-works")
	contentType.Store(".web", "application/vnd.xara")
	contentType.Store(".wi", "image/wavelet")
	contentType.Store(".wis", "application/x-InstallShield")
	contentType.Store(".wks", "application/vnd.ms-works")
	contentType.Store(".wm", "video/x-ms-wm")
	contentType.Store(".wma", "audio/x-ms-wma")
	contentType.Store(".wmd", "application/x-ms-wmd")
	contentType.Store(".wmf", "application/x-msmetafile")
	contentType.Store(".wml", "text/vnd.wap.wml")
	contentType.Store(".wmlc", "application/vnd.wap.wmlc")
	contentType.Store(".wmls", "text/vnd.wap.wmlscript")
	contentType.Store(".wmlsc", "application/vnd.wap.wmlscriptc")
	contentType.Store(".wmlscript", "text/vnd.wap.wmlscript")
	contentType.Store(".wmv", "audio/x-ms-wmv")
	contentType.Store(".wmx", "video/x-ms-wmx")
	contentType.Store(".wmz", "application/x-ms-wmz")
	contentType.Store(".wpng", "image/x-up-wpng")
	contentType.Store(".wps", "application/vnd.ms-works")
	contentType.Store(".wpt", "x-lml/x-gps")
	contentType.Store(".wri", "application/x-mswrite")
	contentType.Store(".wrl", "x-world/x-vrml")
	contentType.Store(".wrz", "x-world/x-vrml")
	contentType.Store(".ws", "text/vnd.wap.wmlscript")
	contentType.Store(".wsc", "application/vnd.wap.wmlscriptc")
	contentType.Store(".wv", "video/wavelet")
	contentType.Store(".wvx", "video/x-ms-wvx")
	contentType.Store(".wxl", "application/x-wxl")
	contentType.Store(".x-gzip", "application/x-gzip")
	contentType.Store(".xaf", "x-world/x-vrml")
	contentType.Store(".xar", "application/vnd.xara")
	contentType.Store(".xbm", "image/x-xbitmap")
	contentType.Store(".xdm", "application/x-xdma")
	contentType.Store(".xdma", "application/x-xdma")
	contentType.Store(".xdw", "application/vnd.fujixerox.docuworks")
	contentType.Store(".xht", "application/xhtml+xml")
	contentType.Store(".xhtm", "application/xhtml+xml")
	contentType.Store(".xhtml", "application/xhtml+xml")
	contentType.Store(".xla", "application/vnd.ms-excel")
	contentType.Store(".xlc", "application/vnd.ms-excel")
	contentType.Store(".xll", "application/x-excel")
	contentType.Store(".xlm", "application/vnd.ms-excel")
	contentType.Store(".xls", "application/vnd.ms-excel")
	contentType.Store(".xlsx", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	contentType.Store(".xlt", "application/vnd.ms-excel")
	contentType.Store(".xlw", "application/vnd.ms-excel")
	contentType.Store(".xm", "audio/x-mod")
	contentType.Store(".xml", "text/plain")
	contentType.Store(".xml", "application/xml")
	contentType.Store(".xmz", "audio/x-mod")
	contentType.Store(".xof", "x-world/x-vrml")
	contentType.Store(".xpi", "application/x-xpinstall")
	contentType.Store(".xpm", "image/x-xpixmap")
	contentType.Store(".xsit", "text/xml")
	contentType.Store(".xsl", "text/xml")
	contentType.Store(".xul", "text/xul")
	contentType.Store(".xwd", "image/x-xwindowdump")
	contentType.Store(".xyz", "chemical/x-pdb")
	contentType.Store(".yz1", "application/x-yz1")
	contentType.Store(".z", "application/x-compress")
	contentType.Store(".zac", "application/x-zaurus-zac")
	contentType.Store(".zip", "application/zip")
	contentType.Store(".json", "application/json")
}
