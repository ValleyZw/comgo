package comgo

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
	"time"
)

const TimeFormat = "02/01/2006T15:04:05.000000"

// New returns configuration parameters of COMTRADE files.
func New() CFG {
	return CFG{}
}

/*
 * CFG - Configuration parameters
 * @StationName: Name of the station
 * @RecordDeviceId: Identification of the recording device
 * @RevisionYear: COMTRADE standard revision year
 * @ChannelNumber: Number of channels
 * @ChannelType: Type of channels
 * @AnalogDetail: Analog channel details
 * @DigitDetail: Digit channel details
 * @LineFrequency: line frequency
 * @SampleRateNum: Sampling rate(s)
 * @SampleDetail: Number of samples at each rate
 * @StartTime: Date and time of first data point
 * @TriggerTime: Date and time of trigger point
 * @DataFileType: Data file type
 * @TimeFactor: Time Stamp multiplication factor
 * @DataFileContent: Store data file content
 */
type CFG struct {
	StationName     string
	RecordDeviceId  string
	RevisionYear    uint16
	ChannelNumber   uint16
	AnalogDetail    *ChannelA
	DigitDetail     *ChannelD
	LineFrequency   uint16
	SampleRateNum   uint16
	SampleDetail    []SampleRate
	StartTime       time.Time
	TriggerTime     time.Time
	DataFileType    string
	TimeFactor      float64
	DataFileContent []byte
}

func (cfg *CFG) GetStationName() string {
	if cfg != nil {
		return cfg.StationName
	}
	return ""
}

func (cfg *CFG) GetRecordDeviceId() string {
	if cfg != nil {
		return cfg.RecordDeviceId
	}
	return ""
}

func (cfg *CFG) GetRevisionYear() uint16 {
	if cfg != nil {
		return cfg.RevisionYear
	}
	return 0
}

func (cfg *CFG) GetChannelNumber() uint16 {
	if cfg != nil {
		return cfg.ChannelNumber
	}
	return 0
}

func (cfg *CFG) GetAnalogDetail() *ChannelA {
	if cfg != nil {
		return cfg.AnalogDetail
	}
	return nil
}

func (cfg *CFG) GetDigitDetail() *ChannelD {
	if cfg != nil {
		return cfg.DigitDetail
	}
	return nil
}

func (cfg *CFG) GetLineFrequency() uint16 {
	if cfg != nil {
		return cfg.LineFrequency
	}
	return 0
}

func (cfg *CFG) GetSampleRateNum() uint16 {
	if cfg != nil {
		return cfg.SampleRateNum
	}
	return 0
}

func (cfg *CFG) GetSampleDetail() []SampleRate {
	if cfg != nil {
		return cfg.SampleDetail
	}
	return nil
}

func (cfg *CFG) GetStartTime() time.Time {
	if cfg != nil {
		return cfg.StartTime
	}
	return time.Time{}
}

func (cfg *CFG) GetTriggerTime() time.Time {
	if cfg != nil {
		return cfg.TriggerTime
	}
	return time.Time{}
}

func (cfg *CFG) GetDataFileType() string {
	if cfg != nil {
		return cfg.DataFileType
	}
	return ""
}

func (cfg *CFG) GetTimeFactor() float64 {
	if cfg != nil {
		return cfg.TimeFactor
	}
	return 0
}

func (cfg *CFG) GetDataFileContent() []byte {
	if cfg != nil {
		return cfg.DataFileContent
	}
	return nil
}

// Return the sampling rate
// only one sampling rate is taking into account
func (cfg *CFG) GetSamplingRate() float64 {
	sampleDetail := cfg.GetSampleDetail()
	if sampleDetail == nil || len(sampleDetail) == 0 {
		return 0
	}
	return sampleDetail[0].GetRate()
}

// Return the number of samples
// only one sampling rate is taking into account
func (cfg *CFG) GetSamplingNumber() int {
	sampleDetail := cfg.GetSampleDetail()
	if sampleDetail == nil || len(sampleDetail) == 0 {
		return 0
	}
	return sampleDetail[0].GetNumber()
}

// Return the names of all analog channel
func (cfg *CFG) GetAnalogChannelNames() []string {
	analogDetail := cfg.GetAnalogDetail()
	if analogDetail != nil {
		return analogDetail.ChannelNames
	}
	return nil
}

/*
 * ChannelA - Analog channel parameters
 * @ChannelTotal: Total number of channels
 * @ChannelNumber: Channel number series
 * @ChannelNames: Names of each channel
 * @ChannelPhases: Phases of each channel
 * @ChannelElements: Channel element (usually null)
 * @ChannelUnits: Units of each channel
 * @ConversionFactors: Conversion factor A and B
 * @TimeFactors: Time factors of each channels
 * @ValueMin: Min Value of each channels
 * @ValueMax: Max Value of each channels
 * @Primary: Primary ratios
 * @Secondary: Secondary ratios
 */
type ChannelA struct {
	ChannelTotal      uint16
	ChannelNumber     []uint16
	ChannelNames      []string
	ChannelPhases     []string
	ChannelElements   []string
	ChannelUnits      []string
	ConversionFactors map[string][]float64
	TimeFactors       []float64
	ValueMin          []int
	ValueMax          []int
	Primary           []float64
	Secondary         []float64
}

func (m *ChannelA) GetChannelTotal() uint16 {
	if m != nil {
		return m.ChannelTotal
	}
	return 0
}

func (m *ChannelA) GetChannelNumber() []uint16 {
	if m != nil {
		return m.ChannelNumber
	}
	return nil
}

func (m *ChannelA) GetChannelNames() []string {
	if m != nil {
		return m.ChannelNames
	}
	return nil
}

func (m *ChannelA) GetChannelPhases() []string {
	if m != nil {
		return m.ChannelPhases
	}
	return nil
}

func (m *ChannelA) GetChannelElements() []string {
	if m != nil {
		return m.ChannelElements
	}
	return nil
}

func (m *ChannelA) GetChannelUnits() []string {
	if m != nil {
		return m.ChannelUnits
	}
	return nil
}

func (m *ChannelA) GetConversionFactors() map[string][]float64 {
	if m != nil {
		return m.ConversionFactors
	}
	return nil
}

func (m *ChannelA) GetTimeFactors() []float64 {
	if m != nil {
		return m.TimeFactors
	}
	return nil
}

func (m *ChannelA) GetValueMin() []int {
	if m != nil {
		return m.ValueMin
	}
	return nil
}

func (m *ChannelA) GetValueMax() []int {
	if m != nil {
		return m.ValueMax
	}
	return nil
}

func (m *ChannelA) GetPrimary() []float64 {
	if m != nil {
		return m.Primary
	}
	return nil
}

func (m *ChannelA) GetSecondary() []float64 {
	if m != nil {
		return m.Secondary
	}
	return nil
}

/*
 * ChannelD - Digit channel parameters
 * @ChannelTotal: Total number of channels
 * @ChannelNumber: Channel number series
 * @ChannelNames: Names of each channel
 * @ChannelPhases: Phases of each channel
 * @ChannelElements: Channel element (usually null)
 */
type ChannelD struct {
	ChannelTotal    uint16
	ChannelNumber   []uint16
	ChannelNames    []string
	ChannelPhases   []string
	ChannelElements []string
	InitialState    []uint8
}

func (m *ChannelD) GetChannelTotal() uint16 {
	if m != nil {
		return m.ChannelTotal
	}
	return 0
}

func (m *ChannelD) GetChannelNumber() []uint16 {
	if m != nil {
		return m.ChannelNumber
	}
	return nil
}

func (m *ChannelD) GetChannelNames() []string {
	if m != nil {
		return m.ChannelNames
	}
	return nil
}

func (m *ChannelD) GetChannelPhases() []string {
	if m != nil {
		return m.ChannelPhases
	}
	return nil
}

func (m *ChannelD) GetChannelElements() []string {
	if m != nil {
		return m.ChannelElements
	}
	return nil
}

func (m *ChannelD) GetInitialState() []uint8 {
	if m != nil {
		return m.InitialState
	}
	return nil
}

/*
 * SampleRate - Sampling rate and sampling number
 * @Rate: Sampling rate
 * @Number: Total number under current sampling rate
 */
type SampleRate struct {
	Rate   float64
	Number int
}

func (m *SampleRate) GetRate() float64 {
	if m != nil {
		return m.Rate
	}
	return 0
}

func (m *SampleRate) GetNumber() int {
	if m != nil {
		return m.Number
	}
	return 0
}

/*
 * BinData - Dat date structure
 * @Sample: Sample series
 * @Stamp: Time Stamp
 * @Value: Analog values: y = factorA * x + factorB
 */
type BinData struct {
	Sample int32
	Stamp  int32
	Value  []int16
}

func (m *BinData) GetSample() int32 {
	if m != nil {
		return m.Sample
	}
	return 0
}

func (m *BinData) GetStamp() int32 {
	if m != nil {
		return m.Stamp
	}
	return 0
}

func (m *BinData) GetValue() []int16 {
	if m != nil {
		return m.Value
	}
	return nil
}

// Reads the Comtrade header file (.cfg).
// return empty CFG and error if err != nil
func (cfg *CFG) ReadCFG(rd io.Reader) (err error) {
	var tempList [][]byte
	content, err := ioutil.ReadAll(rd)
	if err != nil {
		return err
	}
	lines := bytes.Split(content, []byte("\n"))

	// Processing first line
	tempList = bytes.Split(lines[0], []byte(","))
	if len(tempList) < 2 {
		return errors.New("cfg format error")
	}
	cfg.StationName = ByteToString(tempList[0])
	cfg.RecordDeviceId = ByteToString(tempList[1])
	// checking vector length to avoid IndexError
	if len(tempList) > 2 {
		if value, err := strconv.ParseUint(ByteToString(tempList[2]), 10, 16); err != nil {
			return err
		} else {
			cfg.RevisionYear = uint16(value)
		}
	}

	// Processing second line
	tempList = bytes.Split(lines[1], []byte(","))
	if len(tempList) < 3 {
		return errors.New("cfg format error")
	}
	// Total channel number
	if value, err := strconv.ParseUint(ByteToString(tempList[0]), 10, 16); err != nil {
		return err
	} else {
		cfg.ChannelNumber = uint16(value)
	}

	if !bytes.Contains(tempList[1], []byte("A")) || !bytes.Contains(tempList[2], []byte("D")) {
		return errors.New("cfg format error")
	}

	// Initialize analog and digit channels
	chA, chD := ChannelA{}, ChannelD{}
	cfg.AnalogDetail, cfg.DigitDetail = &chA, &chD
	chA.ConversionFactors = make(map[string][]float64)

	// Analog channel total number
	if value, err := strconv.ParseUint(string(bytes.TrimSuffix(bytes.TrimSpace(tempList[1]), []byte("A"))), 10, 16); err != nil {
		return err
	} else {
		chA.ChannelTotal = uint16(value)
	}

	// Digit channel total number
	if value, err := strconv.ParseUint(string(bytes.TrimSuffix(bytes.TrimSpace(tempList[2]), []byte("D"))), 10, 16); err != nil {
		return err
	} else {
		chD.ChannelTotal = uint16(value)
	}

	// Processing analog channels
	for i := 0; i < int(chA.GetChannelTotal()); i++ {
		tempList = bytes.Split(lines[2+i], []byte(","))
		if len(tempList) < 10 {
			return errors.New("cfg format error")
		}
		if num, err := strconv.Atoi(ByteToString(tempList[0])); err != nil {
			return err
		} else {
			chA.ChannelNumber = append(chA.GetChannelNumber(), uint16(num))
		}
		// Format ids to xxx_xxx_xxx
		chA.ChannelNames = append(chA.GetChannelNames(), ByteToString(bytes.Join(bytes.Split(tempList[1], []byte(" ")), []byte("_"))))
		chA.ChannelPhases = append(chA.GetChannelPhases(), ByteToString(tempList[2]))
		// Channel element (usually null)
		chA.ChannelElements = append(chA.GetChannelElements(), ByteToString(tempList[3]))
		chA.ChannelUnits = append(chA.GetChannelUnits(), ByteToString(tempList[4]))
		// Conversion factor A
		if num, err := strconv.ParseFloat(ByteToString(tempList[5]), 64); err != nil {
			return err
		} else {
			chA.ConversionFactors["a"] = append(chA.GetConversionFactors()["a"], num)
		}
		// Conversion factor B
		if num, err := strconv.ParseFloat(ByteToString(tempList[6]), 64); err != nil {
			return err
		} else {
			chA.ConversionFactors["b"] = append(chA.GetConversionFactors()["b"], num)
		}
		// Time factor
		if num, err := strconv.ParseFloat(ByteToString(tempList[7]), 64); err != nil {
			return err
		} else {
			chA.TimeFactors = append(chA.GetTimeFactors(), num)
		}
		// Min Value at current channel
		if num, err := strconv.Atoi(ByteToString(tempList[8])); err != nil {
			return err
		} else {
			chA.ValueMin = append(chA.GetValueMin(), num)
		}
		// Max Value at current channel
		if num, err := strconv.Atoi(ByteToString(tempList[9])); err != nil {
			return err
		} else {
			chA.ValueMax = append(chA.GetValueMax(), num)
		}

		if len(tempList) > 10 {
			if num, err := strconv.ParseFloat(ByteToString(tempList[10]), 64); err == nil {
				chA.Primary = append(chA.GetPrimary(), num)
			}
		}
		if len(tempList) > 11 {
			if num, err := strconv.ParseFloat(ByteToString(tempList[11]), 64); err == nil {
				chA.Secondary = append(chA.GetSecondary(), num)
			}
		}
	}

	// Processing digit channels
	for i := 0; i < int(chD.GetChannelTotal()); i++ {
		tempList = bytes.Split(lines[2+int(chA.GetChannelTotal())+i], []byte(","))
		if len(tempList) < 3 {
			return errors.New("cfg format error")
		}
		if num, err := strconv.Atoi(ByteToString(tempList[0])); err != nil {
			return err
		} else {
			chD.ChannelNumber = append(chD.GetChannelNumber(), uint16(num))
		}
		chD.ChannelNames = append(chD.GetChannelNames(), ByteToString(bytes.Join(bytes.Split(tempList[1], []byte(" ")), []byte("_"))))
		chD.ChannelPhases = append(chD.GetChannelPhases(), ByteToString(tempList[2]))

		// checking vector length to avoid IndexError
		if len(tempList) > 3 {
			// Channel element (usually null)
			chD.ChannelElements = append(chD.GetChannelElements(), ByteToString(tempList[3]))
		}
		if len(tempList) > 4 {
			if num, err := strconv.ParseUint(ByteToString(tempList[4]), 10, 8); err != nil {
				return err
			} else {
				chD.InitialState = append(chD.GetInitialState(), uint8(num))
			}
		}
	}

	// Read line frequency
	tempList = bytes.Split(lines[2+chA.GetChannelTotal()+chD.GetChannelTotal()], []byte(","))
	if num, err := strconv.ParseFloat(ByteToString(tempList[0]), 64); err != nil {
		return err
	} else {
		cfg.LineFrequency = uint16(num)
	}

	// Read sampling rate num
	tempList = bytes.Split(lines[3+chA.GetChannelTotal()+chD.GetChannelTotal()], []byte(","))
	if num, err := strconv.ParseUint(ByteToString(tempList[0]), 10, 16); err != nil {
		return err
	} else {
		cfg.SampleRateNum = uint16(num)
	}

	// Read Sample number (@TODO only one sampling rate is taking into account)
	for i := 0; i < int(cfg.GetSampleRateNum()); i++ {
		sampleRate := SampleRate{}
		tempList = bytes.Split(lines[4+i+int(chA.GetChannelTotal())+int(chD.GetChannelTotal())], []byte(","))
		if num, err := strconv.ParseFloat(ByteToString(tempList[0]), 64); err != nil {
			return err
		} else {
			sampleRate.Rate = num
		}
		if num, err := strconv.ParseFloat(ByteToString(tempList[1]), 64); err != nil {
			return err
		} else {
			sampleRate.Number = int(num)
		}
		cfg.SampleDetail = append(cfg.GetSampleDetail(), sampleRate)
	}

	// Read start date and time ([dd,mm,yyyy,hh,mm,ss.ssssss])
	tempList = bytes.Split(lines[4+cfg.GetSampleRateNum()+chA.GetChannelTotal()+chD.GetChannelTotal()], []byte(","))
	if start, err := time.Parse(TimeFormat, ByteToString(bytes.Join(tempList, []byte("T")))); err != nil {
		return err
	} else {
		cfg.TriggerTime = start
	}

	// Read trigger date and time ([dd,mm,yyyy,hh,mm,ss.ssssss])
	tempList = bytes.Split(lines[5+cfg.GetSampleRateNum()+chA.GetChannelTotal()+chD.GetChannelTotal()], []byte(","))
	if trigger, err := time.Parse(TimeFormat, ByteToString(bytes.Join(tempList, []byte("T")))); err != nil {
		return err
	} else {
		cfg.StartTime = trigger
	}

	// Read dat content type
	tempList = bytes.Split(lines[6+cfg.GetSampleRateNum()+chA.GetChannelTotal()+chD.GetChannelTotal()], []byte(","))
	cfg.DataFileType = ByteToString(tempList[0])

	// Read time multiplication factor
	tempList = bytes.Split(lines[7+cfg.GetSampleRateNum()+chA.GetChannelTotal()+chD.GetChannelTotal()], []byte(","))
	if !bytes.Equal(tempList[0], []byte("")) {
		if num, err := strconv.ParseFloat(ByteToString(tempList[0]), 64); err != nil {
			return err
		} else {
			cfg.TimeFactor = num
		}
	} else {
		cfg.TimeFactor = 1
	}

	return nil
}

// Reads the contents of the Comtrade .dat file
// Store the contents in a private variable
func (cfg *CFG) ReadDAT(rd io.Reader) (err error) {
	content, err := ioutil.ReadAll(rd)
	if err != nil {
		return err
	}
	cfg.DataFileContent = content
	return nil
}

// Returns an array of numbers containing the data values of the channel number
// num is the number of the channel as in .cfg file
func (cfg *CFG) GetAnalogChannelData(num uint16) (result []float64, err error) {
	if cfg == nil {
		return nil, errors.New("invalid cfg file, read .cfg first")
	}

	if cfg.GetDataFileContent() == nil || bytes.Equal(cfg.GetDataFileContent(), []byte("")) {
		return nil, errors.New("not data content, read .dat first")
	}

	analogDetail := cfg.GetAnalogDetail()
	if analogDetail == nil {
		return nil, errors.New("invalid analog channel")
	}

	if num > analogDetail.GetChannelTotal() {
		return nil, errors.New("analog channel number greater than the total number of channels")
	}

	if num < 1 {
		return nil, errors.New("analog channel number cannot be less than 1")
	}

	digitDetail := cfg.GetDigitDetail()
	if digitDetail == nil {
		return nil, errors.New("invalid digital channel")
	}

	// Number of bytes per Sample:
	NB := 8 + int(analogDetail.GetChannelTotal())<<1 + int(math.Ceil(float64(int(digitDetail.GetChannelTotal()))/float64(16)))<<1

	sampleDetail := cfg.GetSampleDetail()
	if sampleDetail == nil || len(sampleDetail) == 0 {
		return nil, errors.New("invalid or not enough sample detail")
	}

	dataFileContent := cfg.GetDataFileContent()
	if dataFileContent == nil {
		return nil, errors.New("invalid dat file content")
	}

	factor := analogDetail.GetConversionFactors()

	// Number of samples: @TODO - only take 1 rate into account
	// Reading the values from datFileContent string
	for i := 0; i < sampleDetail[0].GetNumber(); i++ {
		s := dataFileContent[i*NB : i*NB+NB]

		var data struct {
			Sample int32
			Stamp  int32
		}

		value := make([]int16, (NB-8)/2) // dynamic slice

		err = binary.Read(bytes.NewReader(s[:8]), binary.LittleEndian, &data)
		if err != nil {
			return nil, err
		}

		err = binary.Read(bytes.NewReader(s[8:]), binary.LittleEndian, &value)
		if err != nil {
			return nil, err
		}

		result = append(result, float64(value[num-1])*factor["a"][num-1]+factor["b"][num-1])
	}

	return result, nil
}

// Convert []byte type file content to string
// Delete extra space
func ByteToString(b []byte) string {
	return strings.TrimSpace(string(b))
}
