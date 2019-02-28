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

/*
 * SampleRate - Sampling rate and sampling number
 * @Rate: Sampling rate
 * @Number: Total number under current sampling rate
 */
type SampleRate struct {
	Rate   float64
	Number int
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
	for i := 0; i < int(chA.ChannelTotal); i++ {
		tempList = bytes.Split(lines[2+i], []byte(","))
		if len(tempList) < 10 {
			return errors.New("cfg format error")
		}
		if num, err := strconv.Atoi(ByteToString(tempList[0])); err != nil {
			return err
		} else {
			chA.ChannelNumber = append(chA.ChannelNumber, uint16(num))
		}
		// Format ids to xxx_xxx_xxx
		chA.ChannelNames = append(chA.ChannelNames, ByteToString(bytes.Join(bytes.Split(tempList[1], []byte(" ")), []byte("_"))))
		chA.ChannelPhases = append(chA.ChannelPhases, ByteToString(tempList[2]))
		// Channel element (usually null)
		chA.ChannelElements = append(chA.ChannelElements, ByteToString(tempList[3]))
		chA.ChannelUnits = append(chA.ChannelUnits, ByteToString(tempList[4]))
		// Conversion factor A
		if num, err := strconv.ParseFloat(ByteToString(tempList[5]), 64); err != nil {
			return err
		} else {
			chA.ConversionFactors["a"] = append(chA.ConversionFactors["a"], num)
		}
		// Conversion factor B
		if num, err := strconv.ParseFloat(ByteToString(tempList[6]), 64); err != nil {
			return err
		} else {
			chA.ConversionFactors["b"] = append(chA.ConversionFactors["b"], num)
		}
		// Time factor
		if num, err := strconv.ParseFloat(ByteToString(tempList[7]), 64); err != nil {
			return err
		} else {
			chA.TimeFactors = append(chA.TimeFactors, num)
		}
		// Min Value at current channel
		if num, err := strconv.Atoi(ByteToString(tempList[8])); err != nil {
			return err
		} else {
			chA.ValueMin = append(chA.ValueMin, num)
		}
		// Max Value at current channel
		if num, err := strconv.Atoi(ByteToString(tempList[9])); err != nil {
			return err
		} else {
			chA.ValueMax = append(chA.ValueMax, num)
		}

		if len(tempList) > 10 {
			if num, err := strconv.ParseFloat(ByteToString(tempList[10]), 64); err == nil {
				chA.Primary = append(chA.Primary, num)
			}
		}
		if len(tempList) > 11 {
			if num, err := strconv.ParseFloat(ByteToString(tempList[11]), 64); err == nil {
				chA.Secondary = append(chA.Secondary, num)
			}
		}
	}

	// Processing digit channels
	for i := 0; i < int(chD.ChannelTotal); i++ {
		tempList = bytes.Split(lines[2+int(chA.ChannelTotal)+i], []byte(","))
		if len(tempList) < 3 {
			return errors.New("cfg format error")
		}
		if num, err := strconv.Atoi(ByteToString(tempList[0])); err != nil {
			return err
		} else {
			chD.ChannelNumber = append(chD.ChannelNumber, uint16(num))
		}
		chD.ChannelNames = append(chD.ChannelNames, ByteToString(bytes.Join(bytes.Split(tempList[1], []byte(" ")), []byte("_"))))
		chD.ChannelPhases = append(chD.ChannelPhases, ByteToString(tempList[2]))

		// checking vector length to avoid IndexError
		if len(tempList) > 3 {
			// Channel element (usually null)
			chD.ChannelElements = append(chD.ChannelElements, ByteToString(tempList[3]))
		}
		if len(tempList) > 4 {
			if num, err := strconv.ParseUint(ByteToString(tempList[4]), 10, 8); err != nil {
				return err
			} else {
				chD.InitialState = append(chD.InitialState, uint8(num))
			}
		}
	}

	// Read line frequency
	tempList = bytes.Split(lines[2+chA.ChannelTotal+chD.ChannelTotal], []byte(","))
	if num, err := strconv.ParseFloat(ByteToString(tempList[0]), 64); err != nil {
		return err
	} else {
		cfg.LineFrequency = uint16(num)
	}

	// Read sampling rate num
	tempList = bytes.Split(lines[3+chA.ChannelTotal+chD.ChannelTotal], []byte(","))
	if num, err := strconv.ParseUint(ByteToString(tempList[0]), 10, 16); err != nil {
		return err
	} else {
		cfg.SampleRateNum = uint16(num)
	}

	// Read Sample number (@TODO only one sampling rate is taking into account)
	for i := 0; i < int(cfg.SampleRateNum); i++ {
		sampleRate := SampleRate{}
		tempList = bytes.Split(lines[4+i+int(chA.ChannelTotal)+int(chD.ChannelTotal)], []byte(","))
		if num, err := strconv.ParseFloat(ByteToString(tempList[0]), 64); err != nil {
			return err
		} else {
			sampleRate.Rate = num
		}
		if num, err := strconv.Atoi(ByteToString(tempList[1])); err != nil {
			return err
		} else {
			sampleRate.Number = num
		}
		cfg.SampleDetail = append(cfg.SampleDetail, sampleRate)
	}

	// Read start date and time ([dd,mm,yyyy,hh,mm,ss.ssssss])
	tempList = bytes.Split(lines[4+cfg.SampleRateNum+chA.ChannelTotal+chD.ChannelTotal], []byte(","))
	if start, err := time.Parse(TimeFormat, ByteToString(bytes.Join(tempList, []byte("T")))); err != nil {
		return err
	} else {
		cfg.TriggerTime = start
	}

	// Read trigger date and time ([dd,mm,yyyy,hh,mm,ss.ssssss])
	tempList = bytes.Split(lines[5+cfg.SampleRateNum+chA.ChannelTotal+chD.ChannelTotal], []byte(","))
	if trigger, err := time.Parse(TimeFormat, ByteToString(bytes.Join(tempList, []byte("T")))); err != nil {
		return err
	} else {
		cfg.StartTime = trigger
	}

	// Read dat content type
	tempList = bytes.Split(lines[6+cfg.SampleRateNum+chA.ChannelTotal+chD.ChannelTotal], []byte(","))
	cfg.DataFileType = ByteToString(tempList[0])

	// Read time multiplication factor
	tempList = bytes.Split(lines[7+cfg.SampleRateNum+chA.ChannelTotal+chD.ChannelTotal], []byte(","))
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
	if bytes.Equal(cfg.DataFileContent, []byte("")) {
		return nil, errors.New("not data content, read .dat first")
	}

	if num > cfg.AnalogDetail.ChannelTotal {
		return nil, errors.New("channel number greater than the total number of channels")
	}

	if num < 1 {
		return nil, errors.New("channel number cannot be less than 1")
	}

	// Number of bytes per Sample:
	NB := 8 + int(cfg.AnalogDetail.ChannelTotal)<<1 + int(math.Ceil(float64(int(cfg.DigitDetail.ChannelTotal))/float64(16)))<<1
	// Number of samples: @TODO - only take 1 rate into account
	N := cfg.SampleDetail[0].Number

	// Reading the values from datFileContent string
	for i := 0; i < N; i++ {
		s := cfg.DataFileContent[i*NB : i*NB+NB]

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

		result = append(result, float64(value[num-1])*cfg.AnalogDetail.ConversionFactors["a"][num-1]+cfg.AnalogDetail.ConversionFactors["b"][num-1])
	}

	return result, nil
}

// Return the sampling rate
// only one sampling rate is taking into account
func (cfg *CFG) GetSamplingRate() float64 {
	return cfg.SampleDetail[0].Rate
}

// Return the number of samples
// only one sampling rate is taking into account
func (cfg *CFG) GetSamplingNumber() int {
	return cfg.SampleDetail[0].Number
}

// Return the names of all analog channel
func (cfg *CFG) GetAnalogChannelNames() []string {
	return cfg.AnalogDetail.ChannelNames
}

// Return sampling start time
func (cfg *CFG) GetStartTime() time.Time {
	return cfg.StartTime
}

// Convert []byte type file content to string
// Delete extra space
func ByteToString(b []byte) string {
	return strings.TrimSpace(string(b))
}
