package bzip2

const (
	// beginStreamMagic is the magic number, BZ.
	beginStreamMagic = 0x425a
	// version is the version of bzip file used. 'h' indicates a bzip2 file.
	version = 'h'
	// endStreamMagic is the BCD(binary coded decimal) of the square root of pi.
	endStreamMagic = 0x177245385090
)
