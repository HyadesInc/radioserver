package protocol

import (
	"fmt"
	"unsafe"
)

// Most of the constants here are to be compatible with SpyServer
// Some extensions were made to support different modes / SDRs

// Defined by ((Major) << 24) | ((Minor) << 16) | (Revision)
// Spyserver Standard

type Version struct {
	Major    int
	Minor    int
	Revision int
}

func (v *Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Revision)
}

func (v *Version) ToUint32() uint32 {
	return GenProtocolVersion(*v)
}

func GenProtocolVersion(version Version) uint32 {
	return uint32(((version.Major) << 24) | ((version.Minor) << 16) | (version.Revision))
}

func SplitProtocolVersion(protocol uint32) Version {
	major := int(((protocol & (0xFF << 24)) >> 24) & 0xFF)
	minor := int(((protocol & (0xFF << 16)) >> 16) & 0xFF)
	revision := int(protocol & 0xFFFF)

	return Version{
		Major:    major,
		Minor:    minor,
		Revision: revision,
	}
}

const DefaultPort = 5555
const DefaultFFTDisplayPixels = 2000
const DefaultFFTRange = 127

// region Limit Values
const FFTMaxDisplayPixels = 1 << 15
const FFTMinDisplayPixels = 100
const FFTMaxDBRange = 150
const FFTMinDBRange = 10
const FFTMaxDBOffset = 100

// endregion

// region Internal States
const (
	ParserAcquiringHeader = iota
	ParserReadingData     = iota
)

// endregion

// DeviceIds IDs
const (
	// Spyserver Standard
	DeviceInvalid   = 0
	DeviceAirspyOne = 1
	DeviceAirspyHf  = 2
	DeviceRtlsdr    = 3

	// Radio Server Standard
	DeviceLimeSDRUSB  = 100000
	DeviceLimeSDRMini = 100001
	DeviceSpyServer   = 100002
	DeviceHackRF      = 100003
)

// DeviceNames names of the devices
const (
	DeviceInvalidName     = "Invalid Device"
	DeviceAirspyOneName   = "Airspy Mini / R2"
	DeviceAirspyHFName    = "Airspy HF / HF+"
	DeviceRtlsdrName      = "RTLSDR"
	DeviceLimeSDRUSBName  = "LimeSDR USB"
	DeviceLimeSDRMiniName = "LimeSDR Mini"
	DeviceHackRFName      = "HackRF"
	DeviceSpyserverName   = "SpyServer"
)

// DeviceName list of device names by their ids
var DeviceName = map[uint32]string{
	DeviceInvalid:     DeviceInvalidName,
	DeviceAirspyOne:   DeviceAirspyOneName,
	DeviceAirspyHf:    DeviceAirspyHFName,
	DeviceRtlsdr:      DeviceRtlsdrName,
	DeviceLimeSDRUSB:  DeviceLimeSDRUSBName,
	DeviceLimeSDRMini: DeviceLimeSDRMiniName,
	DeviceHackRF:      DeviceHackRFName,
	DeviceSpyServer:   DeviceSpyserverName,
}

const (
	CmdHello      = 0
	CmdGetSetting = 1
	CmdSetSetting = 2
	CmdPing       = 3
)

const (
	SettingStreamingMode    = 0
	SettingStreamingEnabled = 1
	SettingGain             = 2
	SettingIqFormat         = 100
	SettingIqFrequency      = 101
	SettingIqDecimation     = 102
	SettingDigitalGain      = 103
	SettingFFTFormat        = 200
	SettingFFTFrequency     = 201
	SettingFFTDecimation    = 202
	SettingFFTDbOffset      = 203
	SettingFFTDbRange       = 204
	SettingFFTDisplayPixels = 205
)

// SettingNames list of device names by their ids
var SettingNames = map[uint32]string{
	SettingStreamingMode:    "Streaming Mode",
	SettingStreamingEnabled: "Streaming Enabled",
	SettingGain:             "Gain",
	SettingIqFormat:         "IQ Format",
	SettingIqFrequency:      "IQ Frequency",
	SettingIqDecimation:     "IQ Decimation",
	SettingDigitalGain:      "Digital Gain",
	SettingFFTFormat:        "FFT Format",
	SettingFFTFrequency:     "FFT Frequency",
	SettingFFTDecimation:    "FFT Decimation",
	SettingFFTDbOffset:      "FFT dB Offset",
	SettingFFTDbRange:       "FFT dB Range",
	SettingFFTDisplayPixels: "FFT Display Pixels",
}

var PossibleSettings = []uint32{
	SettingStreamingMode,
	SettingStreamingEnabled,
	SettingGain,

	SettingIqFormat,
	SettingIqFrequency,
	SettingIqDecimation,

	SettingFFTFormat,
	SettingFFTFrequency,
	SettingFFTDecimation,
	SettingFFTDbOffset,
	SettingFFTDbRange,
	SettingFFTDisplayPixels,
}

var GlobalAffectedSettings = []uint32{
	SettingGain,
}

func IsSettingPossible(setting uint32) bool {
	for _, v := range PossibleSettings {
		if setting == v {
			return true
		}
	}

	return false
}

func SettingAffectsGlobal(setting uint32) bool {
	for _, v := range GlobalAffectedSettings {
		if setting == v {
			return true
		}
	}

	return false
}

// StreamTypes is a enum that defines which stream types the spyserver supports.
const (
	StreamTypeStatus = 0
	StreamTypeIQ     = 1
	StreamTypeAF     = 2
	StreamTypeFFT    = 4
)

const (
	// StreamModeIQOnly only enables IQ Channel
	StreamModeIQOnly = StreamTypeIQ

	StreamModeAFOnly = StreamTypeAF

	// StreamModeFFTOnly only enables FFT Channel
	StreamModeFFTOnly = StreamTypeFFT

	// StreamModeFFTOnly only enables both IQ and FFT Channels
	StreamModeFFTIQ = StreamTypeFFT | StreamTypeIQ

	StreamModeFFTAF = StreamTypeFFT | StreamTypeAF
)

const (
	StreamFormatInvalid = 0
	StreamFormatUint8   = 1
	StreamFormatInt16   = 2
	StreamFormatInt24   = 3
	StreamFormatFloat   = 4
	StreamFormatDint4   = 5
)

const (
	MsgTypeDeviceInfo  = 0
	MsgTypeClientSync  = 1
	MsgTypePong        = 2
	MsgTypeReadSetting = 3
	MsgTypeUint8IQ     = 100
	MsgTypeInt16IQ     = 101
	MsgTypeInt24IQ     = 102
	MsgTypeFloatIQ     = 103
	MsgTypeUint8AF     = 200
	MsgTypeInt16AF     = 201
	MsgTypeInt24AF     = 202
	MsgTypeFloatAF     = 203
	MsgTypeDint4FFT    = 300
	MsgTypeUint8FFT    = 301
)

type MessageHeader struct {
	ProtocolID     uint32
	MessageType    uint32
	StreamType     uint32
	SequenceNumber uint32
	BodySize       uint32
}

type CommandHeader struct {
	CommandType uint32
	BodySize    uint32
}

type DeviceInfo struct {
	DeviceType           uint32
	DeviceSerial         uint32
	MaximumSampleRate    uint32
	MaximumBandwidth     uint32
	DecimationStageCount uint32
	GainStageCount       uint32
	MaximumGainIndex     uint32
	MinimumFrequency     uint32
	MaximumFrequency     uint32
	Resolution           uint32
	MinimumIQDecimation  uint32
	ForcedIQFormat       uint32
}

type ClientSync struct {
	CanControl                uint32
	Gain                      uint32
	DeviceCenterFrequency     uint32
	IQCenterFrequency         uint32
	FFTCenterFrequency        uint32
	MinimumIQCenterFrequency  uint32
	MaximumIQCenterFrequency  uint32
	MinimumFFTCenterFrequency uint32
	MaximumFFTCenterFrequency uint32
}

type PingPacket struct {
	Timestamp int64
}

const MessageHeaderSize = uint32(unsafe.Sizeof(MessageHeader{}))
const CommandHeaderSize = uint32(unsafe.Sizeof(CommandHeader{}))
const MaxMessageBodySize = 1 << 20
